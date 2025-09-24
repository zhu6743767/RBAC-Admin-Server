package permission_api

import "github.com/gin-gonic/gin"

// PermissionApi 权限API结构体
type PermissionApi struct{}

// NewPermissionApi 创建权限API实例
func NewPermissionApi() *PermissionApi {
	return &PermissionApi{}
}

// RegisterRoutes 注册权限API路由
func (p *PermissionApi) RegisterRoutes(router *gin.RouterGroup) {
	permissionRouter := router.Group("/permission")
	{
		permissionRouter.GET("/list", p.GetPermissionList)
		permissionRouter.POST("/create", p.CreatePermission)
		permissionRouter.PUT("/update", p.UpdatePermission)
		permissionRouter.DELETE("/delete", p.DeletePermission)
		permissionRouter.GET("/tree", p.GetPermissionTree)
		permissionRouter.GET("/role-permissions", p.GetRolePermissions)
	}
}