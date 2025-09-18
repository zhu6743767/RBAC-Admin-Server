package middleware

import (
	"net/http"
	"rbac.admin/global"
	"rbac.admin/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}

		// 验证token
		token = strings.Replace(token, "Bearer ", "", 1)
		claims, err := global.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}

		// 获取用户信息
		var user models.User
		err = global.DB.Where("id = ?", claims.UserID).First(&user).Error
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "msg": "用户不存在"})
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != 1 {
			c.JSON(401, gin.H{"code": 401, "msg": "用户已被禁用"})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// Admin 管理员权限中间件
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(403, gin.H{"code": 403, "msg": "未登录"})
			c.Abort()
			return
		}

		u := user.(models.User)
		if !u.IsAdmin {
			c.JSON(403, gin.H{"code": 403, "msg": "无权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logrus.Infof("[%s] %s %s %d %v",
			clientIP,
			method,
			path,
			statusCode,
			latency,
		)
	}
}