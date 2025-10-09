package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"rbac_admin_server/global"
)

// HealthApi 健康检查API处理器
type HealthApi struct {}

// NewHealthApi 创建健康检查API实例
func NewHealthApi() *HealthApi {
	return &HealthApi{}
}

// HealthCheck 处理健康检查请求
func (h *HealthApi) HealthCheck(c *gin.Context) {
	// 检查数据库连接
dbStatus := "OK"
// 使用Exec执行简单查询来测试数据库连接，避免空结构体导致的model accessible fields required错误
dbErr := global.DB.Exec("SELECT 1").Error
if dbErr != nil {
	dbStatus = "ERROR: " + dbErr.Error()
}

	// 检查Redis连接
	redisStatus := "OK"
	redisErr := global.Redis.Ping(global.RedisCtx).Err()
	if redisErr != nil {
		redisStatus = "ERROR: " + redisErr.Error()
	}

	// 返回健康状态
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "success",
		"status":   "ok",
		"database": dbStatus,
		"redis":    redisStatus,
		"timestamp": time.Now().Unix(),
	})
}

// RegisterRoutes 注册健康检查路由
func (h *HealthApi) RegisterRoutes(router *gin.Engine) {
	// 使用配置中的健康检查路径，如果没有配置则使用默认路径
	healthCheckPath := global.Config.Monitoring.HealthCheckPath
	if healthCheckPath == "" {
		healthCheckPath = "/health"
	}

	router.GET(healthCheckPath, h.HealthCheck)
}