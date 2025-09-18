package routes

import (
	"rbac.admin/api"
	"rbac.admin/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建默认路由
	r := gin.Default()

	// 设置跨域中间件
	r.Use(middleware.Cors())

	// 设置静态文件
	r.Static("/static", "./static")

	// 设置API路由
	apiGroup := r.Group("/api")
	{
		// 公共路由
		publicGroup := apiGroup.Group("/public")
		{
			publicGroup.POST("/login", api.Login)
			publicGroup.POST("/register", api.Register)
		}

		// 需要认证的路由
		authGroup := apiGroup.Group("/admin")
		authGroup.Use(middleware.Auth())
		{
			// 用户管理
			authGroup.GET("/user/list", api.GetUserList)
			authGroup.POST("/user/create", api.CreateUser)
			authGroup.PUT("/user/update", api.UpdateUser)
			authGroup.DELETE("/user/delete", api.DeleteUser)

			// 角色管理
			roleGroup := authGroup.Group("/role")
			{
				roleGroup.GET("/list", api.GetRoleList)
				roleGroup.POST("/create", api.CreateRole)
				roleGroup.PUT("/update", api.UpdateRole)
				roleGroup.DELETE("/delete", api.DeleteRole)
			}

			// 权限管理
			permissionGroup := authGroup.Group("/permission")
			{
				permissionGroup.GET("/list", api.GetPermissionList)
				permissionGroup.POST("/create", api.CreatePermission)
				permissionGroup.PUT("/update", api.UpdatePermission)
				permissionGroup.DELETE("/delete", api.DeletePermission)
			}

			// 部门管理
			deptGroup := authGroup.Group("/dept")
			{
				deptGroup.GET("/list", api.GetDepartmentList)
				deptGroup.POST("/create", api.CreateDepartment)
				deptGroup.PUT("/update", api.UpdateDepartment)
				deptGroup.DELETE("/delete", api.DeleteDepartment)
			}

			// 菜单管理
			menuGroup := authGroup.Group("/menu")
			{
				menuGroup.GET("/list", api.GetMenuList)
				menuGroup.POST("/create", api.CreateMenu)
				menuGroup.PUT("/update", api.UpdateMenu)
				menuGroup.DELETE("/delete", api.DeleteMenu)
			}

			// 文件管理
			fileGroup := authGroup.Group("/file")
			{
				fileGroup.POST("/upload", api.UploadFile)
				fileGroup.GET("/list", api.GetFileList)
				fileGroup.DELETE("/delete", api.DeleteFile)
			}

			// 日志管理
			logGroup := authGroup.Group("/log")
			{
				logGroup.GET("/list", api.GetLogList)
				logGroup.DELETE("/delete", api.DeleteLog)
			}

			// 个人中心
			profileGroup := authGroup.Group("/profile")
			{
				profileGroup.GET("/info", api.GetProfile)
				profileGroup.PUT("/update", api.UpdateProfile)
				profileGroup.PUT("/password", api.UpdatePassword)
			}
		}
	}

	return r
}