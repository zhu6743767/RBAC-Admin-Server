package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rbac_admin_server/api"
	"rbac_admin_server/api/user_api"
	"rbac_admin_server/global"
	"rbac_admin_server/middleware"
)

// SetDebugMode 设置Gin为调试模式
func SetDebugMode() {
	gin.SetMode(gin.DebugMode)
}

// SetReleaseMode 设置Gin为发布模式
func SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建默认路由
	r := gin.Default()

	// 设置跨域中间件
	r.Use(middleware.Cors())

	// 设置静态文件目录
	r.Static("/uploads", "./uploads")

	// 初始化API
	api.InitApi()

	// 创建用户API实例用于处理登录和注册
	userApi := user_api.NewUserApi()

	// 公共路由组
	public := r.Group("/public")
	{
		// 登录接口
		public.POST("/login", userApi.Login)
		// 注册接口
		public.POST("/register", userApi.Register)
	}

	// 需要认证的路由组
	admin := r.Group("/admin")
	admin.Use(middleware.Auth())
	{
		// 用户管理模块
		api.App.UserApi.RegisterRoutes(admin)

		// 角色管理模块
		api.App.RoleApi.RegisterRoutes(admin)

		// 权限管理模块
		api.App.PermissionApi.RegisterRoutes(admin)

		// 部门管理模块
		api.App.DeptApi.RegisterRoutes(admin)

		// 菜单管理模块
		api.App.MenuApi.RegisterRoutes(admin)

		// 文件管理模块
		api.App.FileApi.RegisterRoutes(admin)

		// 日志管理模块
		api.App.LogApi.RegisterRoutes(admin)

		// 个人中心模块
		api.App.ProfileApi.RegisterRoutes(admin)
	}

	return r
}

// Run 运行路由和HTTP服务器
func Run() {
	// 创建路由
	router := SetupRouter()

	// 获取系统配置
	s := global.Config.System

	// 启动HTTP服务器
	logrus.Infof("后端服务运行在 http://%s:%d", s.IP, s.Port)

	// 启动服务器
	if err := router.Run(fmt.Sprintf("%s:%d", s.IP, s.Port)); err != nil {
		logrus.Fatalf("服务器启动失败: %v", err)
	}
}