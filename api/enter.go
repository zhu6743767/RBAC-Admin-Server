package api

import (
	"rbac_admin_server/api/dept_api"
	"rbac_admin_server/api/file_api"
	"rbac_admin_server/api/log_api"
	"rbac_admin_server/api/menu_api"
	"rbac_admin_server/api/permission_api"
	"rbac_admin_server/api/profile_api"
	"rbac_admin_server/api/role_api"
	"rbac_admin_server/api/user_api"
)

// Api 全局API实例，提供所有API接口的访问入口
type Api struct {
	UserApi       *user_api.UserApi
	RoleApi       *role_api.RoleApi
	PermissionApi *permission_api.PermissionApi
	DeptApi       *dept_api.DeptApi
	MenuApi       *menu_api.MenuApi
	FileApi       *file_api.FileApi
	LogApi        *log_api.LogApi
	ProfileApi    *profile_api.ProfileApi
}

// App 全局API实例，供外部调用
var App = new(Api)

// InitApi 初始化所有API模块
func InitApi() {
	App.UserApi = user_api.NewUserApi()
	App.RoleApi = role_api.NewRoleApi()
	App.PermissionApi = permission_api.NewPermissionApi()
	App.DeptApi = dept_api.NewDeptApi()
	App.MenuApi = menu_api.NewMenuApi()
	App.FileApi = file_api.NewFileApi()
	App.LogApi = log_api.NewLogApi()
	App.ProfileApi = profile_api.NewProfileApi()
}