package main

import (
	"os"
	"os/signal"
	"syscall"

	"rbac_admin_server/core"
	"rbac_admin_server/flags"
	"rbac_admin_server/global"
	"rbac_admin_server/routes"

	"github.com/sirupsen/logrus"
)

func main() {
	// 解析命令行参数
	args := flags.ParseCommandLineArgs()

	// 处理命令行参数
	if err := flags.HandleCommandLineArgs(args); err != nil {
		logrus.Fatalf("❌ 命令处理失败: %v", err)
	}

	// 如果是服务器模式，则启动服务器
	if args.Mode == flags.ModeServer {
		// 初始化系统核心组件
		global.Logger.Info("开始初始化系统核心组件...")
		if err := core.InitSystem(); err != nil {
			global.Logger.Fatalf("❌ 系统初始化失败: %v", err)
		}
		global.Logger.Info("✅ 系统核心组件初始化完成")

		// 启动应用
		global.Logger.Info("开始启动HTTP服务器...")
		routes.Run()

		// 等待系统信号，优雅退出
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		// 清理系统资源
		core.CleanupSystem()

		global.Logger.Info("✅ 服务已优雅停止")
	} else {
		// 非服务器模式，处理完成后退出
		global.Logger.Info("✅ 命令执行完成")
	}
}
