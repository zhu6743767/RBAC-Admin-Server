package examples

import (
	"fmt"
	"rbac.admin/config"
	"rbac.admin/core/logger"
)

// FactoryUsageExample 工厂模式使用示例
func FactoryUsageExample() {
	fmt.Println("=== 日志工厂模式使用示例 ===")

	// 示例1: 使用默认的logrus日志
	fmt.Println("\n1. 使用logrus日志:")
	logrusConfig := config.LogConfig{
		Type:   "logrus",
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	
	logrusFactory, _ := logger.GetFactory("logrus")
	// 转换配置类型
	logConfig := &logger.Config{
		Type:   logrusConfig.Type,
		Level:  logrusConfig.Level,
		Format: logrusConfig.Format,
		Output: logrusConfig.Output,
	}
	logrusLogger, _ := logrusFactory.Create(logConfig)
	logrusLogger.Info("这是logrus日志消息")
	logrusLogger.Debug("这是debug级别消息")
	logrusLogger.Close()

	// 示例2: 动态切换日志类型
	fmt.Println("\n2. 动态切换日志类型:")
	configurations := []config.LogConfig{
		{Type: "logrus", Level: "info", Format: "text", Output: "stdout"},
		{Type: "logrus", Level: "debug", Format: "json", Output: "file", LogDir: "./logs"},
	}

	for _, cfg := range configurations {
		factory, err := logger.GetFactory(cfg.Type)
		if err != nil {
			fmt.Printf("获取工厂失败: %v\n", err)
			continue
		}

		// 转换配置类型
		logConfig := &logger.Config{
			Type:   cfg.Type,
			Level:  cfg.Level,
			Format: cfg.Format,
			Output: cfg.Output,
		}
		logger, err := factory.Create(logConfig)
		if err != nil {
			fmt.Printf("创建日志器失败: %v\n", err)
			continue
		}

		logger.Info(fmt.Sprintf("使用 %s 日志器，级别: %s", cfg.Type, cfg.Level))
		logger.Close()
	}

	// 示例3: 注册自定义日志工厂
	fmt.Println("\n3. 注册自定义日志工厂:")
	// 可以在这里注册自定义的日志工厂
	// logger.RegisterFactory("custom", &CustomLoggerFactory{})

	// 示例4: 获取所有支持的日志类型
	fmt.Println("\n4. 支持的日志类型:")
	// 目前只支持logrus
	fmt.Printf("支持的日志类型: [logrus]\n")

	// 示例5: 使用日志工厂创建不同配置的日志实例
	fmt.Println("\n5. 创建不同配置的日志实例:")
	
	// 创建多个日志实例
	cfg1 := &logger.Config{Type: "logrus", Level: "info", Format: "text", Output: "stdout"}
	cfg2 := &logger.Config{Type: "logrus", Level: "debug", Format: "json", Output: "stdout"}
	
	logger1, _ := logrusFactory.Create(cfg1)
	logger2, _ := logrusFactory.Create(cfg2)
	
	logger1.Info("应用日志")
	logger2.Debug("调试日志")
	
	// 关闭所有日志
	logger1.Close()
	logger2.Close()
}

// 运行示例
// go run examples/factory_usage_example.go