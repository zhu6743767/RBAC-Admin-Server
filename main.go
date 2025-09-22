package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/routes"
	"rbac.admin/utils"
)

var (
	env        = flag.String("env", "dev", "运行环境 (dev/test/prod)")
	configPath = flag.String("f", "", "配置文件路径")
	module     = flag.String("m", "", "模块 (user/db)")
	task       = flag.String("t", "", "任务 (create/migrate)")
)

func main() {
	flag.Parse()

	// 加载配置文件
	cfg, err := loadConfig()
	if err != nil {
		logrus.Fatalf("加载配置文件失败: %v", err)
	}
	
	// 初始化全局日志系统
	if err := core.InitLogger(&cfg.Log); err != nil {
		logrus.Fatalf("初始化日志系统失败: %v", err)
	}
	
	// 设置全局配置
	global.Config = cfg

	// 初始化核心系统组件
	if err := core.InitSystem(); err != nil {
		global.Logger.Fatalf("系统初始化失败: %v", err)
	}

	// 处理命令行参数
	if *module != "" {
		handleCommand()
		return
	}

	// 启动HTTP服务器
	startServer()
}

// loadConfig 加载配置
func loadConfig() (*config.Config, error) {
	var configFile string
	if *configPath != "" {
		configFile = *configPath
	} else {
		configFile = fmt.Sprintf("settings_%s.yaml", *env)
	}
	
	return config.Load(configFile)
}

// handleCommand 处理命令行参数
func handleCommand() {
	switch *module {
	case "user":
		handleUserCommand()
	case "db":
		handleDBCommand()
	default:
		fmt.Printf("不支持的模块: %s\n", *module)
		os.Exit(1)
	}
}

// handleUserCommand 处理用户相关命令
func handleUserCommand() {
	switch *task {
	case "create":
		createAdminUser()
	default:
		fmt.Printf("不支持的用户任务: %s\n", *task)
		os.Exit(1)
	}
}

// handleDBCommand 处理数据库相关命令
func handleDBCommand() {
	switch *task {
	case "migrate":
		migrateDatabase()
	default:
		fmt.Printf("不支持的数据库任务: %s\n", *task)
		os.Exit(1)
	}
}

// createAdminUser 创建管理员用户
func createAdminUser() {
	fmt.Println("请输入用户名")
	var username string
	fmt.Scanln(&username)

	if username == "" {
		fmt.Println("用户名不能为空")
		os.Exit(1)
	}

	var user models.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err == nil {
		fmt.Println("用户已存在")
		os.Exit(1)
	}

	fmt.Println("请输入密码")
	password, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Println("读取密码失败")
		os.Exit(1)
	}
	fmt.Println()

	if len(password) < 6 {
		fmt.Println("密码长度不能少于6位")
		os.Exit(1)
	}

	fmt.Println("请再次输入密码")
	rePassword, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Println("读取密码失败")
		os.Exit(1)
	}
	fmt.Println()

	if string(password) != string(rePassword) {
		fmt.Println("两次密码不一致")
		os.Exit(1)
	}

	// 密码加密
	hashPwd := utils.HashedPassword(string(password))

	// 创建用户
	newUser := models.User{
		Username: username,
		Password: hashPwd,
		Nickname: "管理员",
		Avatar:   "/default-avatar.png",
		Email:    "",
		Phone:    "",
		Status:   1,
		IsAdmin:  true,
	}

	if err := global.DB.Create(&newUser).Error; err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("管理员用户创建成功")
}

// migrateDatabase 数据库迁移
func migrateDatabase() {
	fmt.Println("开始执行数据库迁移...")

	// 执行自动迁移
	if err := core.AutoMigrateModels(); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("数据库迁移成功")
}

// startServer 启动HTTP服务器
func startServer() {
	// 设置HTTP服务器
	router := routes.SetupRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", global.Config.System.IP, global.Config.System.Port),
		Handler: router,
	}

	// 启动服务器（异步）
	go func() {
		global.Logger.Infof("🚀 服务器启动成功，访问地址: http://%s:%d", global.Config.System.IP, global.Config.System.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待系统信号，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	global.Logger.Info("开始优雅关闭服务器...")

	// 清理系统资源
	core.CleanupSystem()

	global.Logger.Info("服务器已优雅关闭")
}
