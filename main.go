package main

import (
	"fmt"
	"log"
	"strings"

	"rbac.admin/core"
	"rbac.admin/global"
)

// main 程序入口点 - 使用全新的初始化器架构
func main() {
	fmt.Println("🚀 RBAC管理员服务器启动中...")

	// 创建项目初始化器
	initializer := core.NewInitializer("./settings.yaml")
	
	// 执行完整项目初始化
	if err := initializer.Initialize(); err != nil {
		log.Fatalf("❌ 项目初始化失败: %v", err)
	}

	// 显示启动信息
	displayStartupInfo()

	// 设置优雅关闭
	initializer.WaitForSignal()
}

// displayStartupInfo 显示启动信息
func displayStartupInfo() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("RBAC管理员服务器启动成功!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("端口: %d\n", global.Config.Server.Port)
	fmt.Printf("数据库: %s@%s:%d/%s\n", 
		global.Config.Database.Username,
		global.Config.Database.Host,
		global.Config.Database.Port,
		global.Config.Database.Database)
	fmt.Printf("日志级别: %s\n", global.Config.Log.Level)
	fmt.Println(strings.Repeat("=", 50))
}

// maskJWTSecret 掩码显示JWT密钥
func maskJWTSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
