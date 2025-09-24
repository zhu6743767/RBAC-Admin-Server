package log_api

import "github.com/gin-gonic/gin"

// LogApi 日志API结构体
type LogApi struct{}

// NewLogApi 创建日志API实例
func NewLogApi() *LogApi {
	return &LogApi{}
}

// RegisterRoutes 注册日志API路由
func (l *LogApi) RegisterRoutes(router *gin.RouterGroup) {
	logRouter := router.Group("/log")
	{
		logRouter.GET("/list", l.GetLogList)
		logRouter.GET("/user-logs", l.GetUserLogs)
		logRouter.DELETE("/delete", l.DeleteLog)
		logRouter.DELETE("/delete-multiple", l.DeleteMultipleLogs)
		logRouter.GET("/dashboard", l.GetLogDashboard)
	}
}