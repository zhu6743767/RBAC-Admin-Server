package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"rbac_admin_server/core"
	"rbac_admin_server/global"
	"rbac_admin_server/routes"

	"github.com/sirupsen/logrus"
)

// 命令行参数定义
var (
	configFile = flag.String("settings", "settings.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 初始化日志
	core.InitLogger("logs")

	// 加载配置
	global.Config = core.ReadConfig(*configFile)

	// 初始化系统核心组件
	if err := core.InitSystem(); err != nil {
		global.Logger.Fatalf("❌ 系统初始化失败: %v", err)
	}

	// 启动应用
	routes.Run()

	// 等待系统信号，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 清理系统资源
	core.CleanupSystem()

	logrus.Info("✅ 服务已优雅停止")
}
