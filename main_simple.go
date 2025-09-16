package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"rbac.admin/config"
)

func main() {
	var (
		env       = flag.String("env", "dev", "运行环境: dev/test/prod")
		configPath = flag.String("config", "", "配置文件路径")
	)
	
	flag.Parse()

	// 显示启动横幅
	fmt.Println(`
	🚀 RBAC管理员服务器配置测试
	╔═══════════════════════════════════════╗
	║          Configuration Test          ║
	╚═══════════════════════════════════════╝
	`)

	// 确定配置文件
	var cfgFile string
	if *configPath != "" {
		cfgFile = *configPath
		fmt.Printf("📁 使用指定配置文件: %s\n", cfgFile)
	} else {
		switch strings.ToLower(*env) {
		case "dev", "development":
			cfgFile = "settings_dev.yaml"
		case "test", "testing":
			cfgFile = "settings_test.yaml"
		case "prod", "production":
			cfgFile = "settings_prod.yaml"
		default:
			log.Fatalf("不支持的环境: %s", *env)
		}
		fmt.Printf("🌍 运行环境: %s，使用配置文件: %s\n", *env, cfgFile)
	}

	// 加载配置
	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 显示配置摘要
	fmt.Printf("\n✅ 配置加载成功!\n")
	fmt.Printf("📋 应用信息: %s v%s (%s)\n", cfg.App.Name, cfg.App.Version, cfg.App.Environment)
	fmt.Printf("🖥️  服务器端口: %d\n", cfg.Server.Port)
	fmt.Printf("🗄️  数据库类型: %s\n", cfg.Database.Type)
	fmt.Printf("🔐 JWT颁发者: %s\n", cfg.JWT.Issuer)
	fmt.Printf("📝 日志级别: %s\n", cfg.Log.Level)

	// 检查配置完整性
	if cfg.App.Name == "" {
		log.Fatal("❌ 应用名称不能为空")
	}
	if cfg.JWT.Secret == "" {
		log.Fatal("❌ JWT密钥不能为空")
	}
	if cfg.Database.Type == "" {
		log.Fatal("❌ 数据库类型不能为空")
	}

	fmt.Println("\n🎉 所有配置验证通过!")
}