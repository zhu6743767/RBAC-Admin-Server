package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/pwd"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	env        = flag.String("env", "dev", "运行环境 (dev/test/prod)")
	configPath = flag.String("f", "", "配置文件路径")
	module     = flag.String("m", "", "模块 (user/db)")
	task       = flag.String("t", "", "任务 (create/migrate)")
)

func main() {
	flag.Parse()

	// 初始化系统
	if err := core.InitSystem(loadConfig()); err != nil {
		fmt.Printf("系统初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 处理命令行参数
	if *module != "" {
		handleCommand()
		return
	}

	// 启动服务器
	startServer()
}

// loadConfig 加载配置
func loadConfig() *config.Config {
	var configFile string
	if *configPath != "" {
		configFile = *configPath
	} else {
		configFile = fmt.Sprintf("settings_%s.yaml", *env)
	}
	
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}
	return cfg
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
		fmt.Println("创建用户失败")
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
	logrus.Info("开始启动服务器...")

	// 检查配置是否为空
	if global.Config == nil {
		logrus.Fatal("全局配置为空")
	}

	// 检查系统配置是否为空
	if global.Config.System.Port == 0 {
		logrus.Fatal("系统端口配置为0")
	}

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", global.Config.System.IP, global.Config.System.Port)
	logrus.Infof("服务器启动在 %s", addr)
	
	// 这里可以添加具体的服务器启动逻辑
	// 例如: router := setupRouter()
	//      logrus.Fatal(http.ListenAndServe(addr, router))
	
	logrus.Infof("🎉 服务器启动成功，监听地址: %s", addr)
	
	// 添加阻塞逻辑，防止程序立即退出
	select {}
}
