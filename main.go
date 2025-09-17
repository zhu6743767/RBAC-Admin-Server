package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
)

var (
	env        = flag.String("env", "dev", "运行环境: dev(开发环境), test(测试环境), prod(生产环境)")
	configPath = flag.String("config", "", "配置文件路径，优先级高于环境选择")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 显示启动横幅
	displayBanner()

	// 根据环境或指定路径加载配置
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 设置全局配置
	global.Config = cfg

	// 初始化系统（包含日志、验证器、数据库、Redis等）
	if err := core.InitSystem(cfg); err != nil {
		log.Fatalf("❌ 系统初始化失败: %v", err)
	}

	// 显示环境信息
	displayEnvironmentInfo()

	// 显示启动信息
	displayStartupInfo()

	// 等待退出信号
	core.WaitForSignal()
}

// loadConfig 根据环境加载对应的配置文件
func loadConfig() (*config.Config, error) {
	var cfgFile string

	if *configPath != "" {
		// 如果指定了配置文件路径，直接使用
		cfgFile = *configPath
		fmt.Printf("📁 使用指定配置文件: %s\n", cfgFile)
	} else {
		// 根据环境选择配置文件
		switch strings.ToLower(*env) {
		case "dev", "development":
			cfgFile = "settings_dev.yaml"
		case "test", "testing":
			cfgFile = "settings_test.yaml"
		case "prod", "production":
			cfgFile = "settings_prod.yaml"
		default:
			return nil, fmt.Errorf("不支持的环境: %s，请使用 dev/test/prod", *env)
		}
		fmt.Printf("🌍 运行环境: %s，使用配置文件: %s\n", *env, cfgFile)
	}

	// 加载配置
	return config.Load(cfgFile)
}

// displayBanner 显示启动横幅
func displayBanner() {
	fmt.Println(`
	🚀 RBAC管理员服务器启动中...
	╔═══════════════════════════════════════╗
	║          RBAC Admin Server            ║
	║    Role-Based Access Control System   ║
	╚═══════════════════════════════════════╝
	`)
}

// displayEnvironmentInfo 显示环境信息
func displayEnvironmentInfo() {
	fmt.Printf("🌍 运行环境: %s\n", strings.ToUpper(*env))
	fmt.Printf("📁 配置文件: %s\n", getCurrentConfigFile())
	fmt.Printf("🗄️  数据库: %s\n", getDatabaseInfo())
	fmt.Println(strings.Repeat("─", 50))
}

// displayStartupInfo 显示启动信息
func displayStartupInfo() {
	fmt.Println(strings.Repeat("═", 50))
	fmt.Println("✅ RBAC管理员服务器启动成功!")
	fmt.Println(strings.Repeat("═", 50))
	fmt.Printf("🌐 访问地址: http://localhost:%d\n", global.Config.System.Port)
	fmt.Printf("📊 健康检查: http://localhost:%d/health\n", global.Config.System.Port)
	fmt.Printf("📈 监控指标: http://localhost:%d/metrics\n", global.Config.System.Port)

	if global.Config.Swagger.Enable && global.Config.Swagger.EnableUI {
		fmt.Printf("📚 API文档: http://localhost:%d/swagger/index.html\n", global.Config.System.Port)
	}

	fmt.Printf("🗄️  数据库: %s@%s:%d/%s\n",
		global.Config.DB.User,
		global.Config.DB.Host,
		global.Config.DB.Port,
		global.Config.DB.DbNAME)
	fmt.Printf("📊 日志级别: %s\n", global.Config.Log.Level)
	fmt.Println(strings.Repeat("═", 50))
}

// getCurrentConfigFile 获取当前使用的配置文件
func getCurrentConfigFile() string {
	if *configPath != "" {
		return *configPath
	}

	switch strings.ToLower(*env) {
	case "dev", "development":
		return "settings_dev.yaml"
	case "test", "testing":
		return "settings_test.yaml"
	case "prod", "production":
		return "settings_prod.yaml"
	default:
		return "settings_dev.yaml"
	}
}

// getDatabaseInfo 获取数据库信息
func getDatabaseInfo() string {
	if global.Config.DB.Mode == "sqlite" {
		return fmt.Sprintf("SQLite(%s)", global.Config.DB.DbNAME)
	}
	return fmt.Sprintf("MySQL(%s@%s:%d/%s)",
		global.Config.DB.User,
		global.Config.DB.Host,
		global.Config.DB.Port,
		global.Config.DB.DbNAME)
}

// 简化版的初始化器（适配新的配置加载方式）
type Initializer struct {
	config *config.Config
}

// NewInitializer 创建新的初始化器
func NewInitializer(cfg *config.Config) *Initializer {
	return &Initializer{
		config: cfg,
	}
}

// Initialize 执行完整项目初始化
func (i *Initializer) Initialize() error {
	// 这里可以添加数据库连接、Redis连接等初始化逻辑
	// 目前简化处理，实际项目中可以扩展
	return nil
}

// WaitForSignal 等待退出信号
func (i *Initializer) WaitForSignal() {
	// 简化的优雅关闭逻辑
	fmt.Println("🔄 服务器运行中... 按 Ctrl+C 退出")
	// 等待信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\n🛑 服务器正在关闭...")
	fmt.Println("✅ 服务器已安全关闭")
}
