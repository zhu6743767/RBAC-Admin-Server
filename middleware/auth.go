package middleware

import (
	"net/http"
	"rbac_admin_server/global"
	"time"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// 允许跨域请求，设置必要的HTTP头
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		h.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		h.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		h.Set("Access-Control-Allow-Credentials", "true")

		// 处理 OPTIONS 请求
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Auth 认证中间件
// 验证JWT token的有效性，并将用户信息存入上下文
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			c.Abort()
			return
		}

		// 移除 Bearer 前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := global.ParseToken(token)
		if err != nil {
			global.Logger.Warnf("Token解析失败: %s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token 无效"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleList", claims.RoleList)
		global.Logger.Debugf("用户认证成功: %s, 角色: %v", claims.Username, claims.RoleList)
		c.Next()
	}
}

// Admin 管理员中间件
// 验证用户是否具有管理员权限
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请先登录"})
			c.Abort()
			return
		}

		// 查询用户信息
		var user struct {
			IsAdmin bool
		}
		if err := global.DB.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
			global.Logger.Error("查询用户信息失败: " + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户不存在"})
			c.Abort()
			return
		}

		if !user.IsAdmin {
			global.Logger.Warnf("用户ID: %v 尝试访问管理员资源但无权限", userID)
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无管理员权限"})
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

		global.Logger.Infof("[%s] %s %s %d %v",
			clientIP,
			method,
			path,
			statusCode,
			latency,
		)
	}
}