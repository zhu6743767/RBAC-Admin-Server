package main

import (
	"flag"
	"fmt"
	"os"

	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
)

var (
	env = flag.String("env", "dev", "运行环境 (dev/test/prod)")
)

func main() {
	flag.Parse()

	// 加载配置文件
	configFile := fmt.Sprintf("settings_%s.yaml", *env)
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}
	global.Config = cfg

	// 初始化全局日志
	if err := core.InitLogger(&global.Config.Log); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化系统
	if err := core.InitSystem(); err != nil {
		fmt.Printf("系统初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 打印系统配置，验证初始化是否成功
	fmt.Println("\n系统初始化成功！")
	fmt.Printf("服务器配置: %s:%d\n", global.Config.System.IP, global.Config.System.Port)
	fmt.Printf("数据库类型: %s\n", global.Config.DB.Mode)
	fmt.Printf("JWT配置: 有效期=%d小时, 刷新令牌有效期=%d小时\n", global.Config.JWT.ExpireHours, global.Config.JWT.RefreshExpireHours)
	
	// 检查Redis连接
	if global.Redis != nil {
		fmt.Println("Redis连接已初始化")
	} else {
		fmt.Println("Redis未初始化（可能是因为配置中未设置）")
	}
	
	// 检查Casbin实例
	if global.Casbin != nil {
		fmt.Println("Casbin权限管理已初始化")
	} else {
		fmt.Println("Casbin未初始化")
	}
}