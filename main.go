package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/routes"
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
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
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
	rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
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
	hashPwd := pwd.HashedPassword(string(password))

	// 创建管理员用户
	err = global.DB.Create(&models.User{
		Username: username,
		Password: hashPwd,
		Nickname: "管理员",
		Status:   1,
		IsAdmin:  true,
	}).Error
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("创建用户成功")
	os.Exit(0)
}

// migrateDatabase 数据库迁移
func migrateDatabase() {
	fmt.Println("开始数据库迁移...")
	err := global.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Department{},
		&models.Menu{},
		&models.File{},
		&models.Log{},
		&gormadapter.CasbinRule{},
	)
	if err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("数据库迁移成功")
	os.Exit(0)
}

// startServer 启动服务器
func startServer() {
	global.Logger.Info("开始启动服务器...")

	// 检查配置是否为空
	if global.Config == nil {
		global.Logger.Fatal("全局配置为空")
	}

	// 检查系统配置是否为空
	if global.Config.System.Port == 0 {
		global.Logger.Fatal("系统端口配置为0")
	}

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", global.Config.System.IP, global.Config.System.Port)
	
	// 设置路由
	router := routes.SetupRouter()
	
	// 创建HTTP服务器
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	
	// 启动服务器的goroutine
	go func() {
		global.Logger.Infof("🎉 服务器启动成功，监听地址: %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatalf("服务器启动失败: %v", err)
		}
	}()
	
	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Logger.Info("正在优雅关闭服务器...")
	
	// 关闭Redis连接
	if global.Redis != nil {
		if err := global.Redis.Close(); err != nil {
			global.Logger.Errorf("关闭Redis连接失败: %v", err)
		} else {
			global.Logger.Info("Redis连接已关闭")
		}
	}
	
	// 关闭HTTP服务器
	if err := server.Shutdown(nil); err != nil {
		global.Logger.Fatalf("服务器强制关闭: %v", err)
	}
	
	global.Logger.Info("服务器已优雅关闭")
}
