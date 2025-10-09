package profile_api

import (
	"github.com/gin-gonic/gin"
	"rbac_admin_server/middleware"
)

// ProfileApi 个人信息管理API
type ProfileApi struct{}

// NewProfileApi 创建ProfileApi实例
func NewProfileApi() *ProfileApi {
	return &ProfileApi{}
}

// RegisterRoutes 注册个人信息管理相关路由
func (p *ProfileApi) RegisterRoutes(router *gin.RouterGroup) {
	profileRouter := router.Group("/profile")
	{
		// 需要认证的路由
		profileRouter.Use(middleware.Auth())
		{
			profileRouter.GET("/info", p.GetUserInfo)           // 获取用户个人信息
			profileRouter.PUT("/info", p.UpdateUserInfo)         // 更新用户个人信息
			profileRouter.PUT("/password", p.UpdatePassword)      // 修改密码
			profileRouter.GET("/dashboard", p.GetDashboardData)   // 获取仪表盘数据
			profileRouter.GET("/settings", p.GetUserSettings)     // 获取用户设置
			profileRouter.PUT("/settings", p.UpdateUserSettings)  // 更新用户设置
		}
	}
}