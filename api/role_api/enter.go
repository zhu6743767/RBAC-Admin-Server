package role_api

import "github.com/gin-gonic/gin"

// RoleApi 角色API结构体
type RoleApi struct{}

// NewRoleApi 创建角色API实例
func NewRoleApi() *RoleApi {
	return &RoleApi{}
}

// RegisterRoutes 注册角色API路由
func (r *RoleApi) RegisterRoutes(router *gin.RouterGroup) {
	roleRouter := router.Group("/role")
	{
		roleRouter.GET("/list", r.GetRoleList)
		roleRouter.POST("/create", r.CreateRole)
		roleRouter.PUT("/update", r.UpdateRole)
		roleRouter.DELETE("/delete", r.DeleteRole)
		roleRouter.GET("/permissions", r.GetRolePermissions)
		roleRouter.POST("/set-permissions", r.SetRolePermissions)
		roleRouter.GET("/users", r.GetRoleUsers)
	}
}