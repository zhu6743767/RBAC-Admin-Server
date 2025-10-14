package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"rbac_admin_server/api"
	"rbac_admin_server/api/captcha_api"
	"rbac_admin_server/api/email_api"
	"rbac_admin_server/api/user_api"
	"rbac_admin_server/global"
	"rbac_admin_server/middleware"
	"rbac_admin_server/utils/captcha"
)

// Run 运行路由和HTTP服务器
func Run() {
	// 获取系统配置
	s := global.Config.System

	// 设置Gin模式
	if s.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建默认路由
	r := gin.Default()

	// 设置跨域中间件
	r.Use(middleware.Cors())

	// 设置静态文件目录
	r.Static("/uploads", "./uploads")

	// 初始化API
	api.InitApi()

	// 注册健康检查路由
	api.App.HealthApi.RegisterRoutes(r)

	// 启动邮件验证码清理定时器
	captcha.EmailStore.StartCleanupTimer()

	// 创建API实例
	userApi := user_api.NewUserApi()
	captchaApi := &captcha_api.CaptchaApi{}
	emailApi := &email_api.EmailApi{}

	// 公共路由组
	public := r.Group("/public")
	{
		// 登录接口
		public.POST("/login", userApi.Login)
		// 注册接口
		public.POST("/register", userApi.Register)
		// 验证码路由
		captchaApi.RegisterRoutes(public)
		// 邮箱路由
		emailApi.RegisterRoutes(public)
	}

	// 需要认证的路由组
	admin := r.Group("/admin")
	// 使用Auth中间件进行身份验证
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

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	global.Logger.Infof("后端服务运行在 http://%s", addr)

	// 启动服务器
	if err := r.Run(addr); err != nil {
		global.Logger.Fatalf("服务器启动失败: %v", err)
	}
}