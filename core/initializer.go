package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"rbac.admin/config"
	"rbac.admin/core/audit"
	"rbac.admin/core/logger"
	"rbac.admin/database"
	"rbac.admin/global"
)

// Initializer 项目初始化器
type Initializer struct {
	configPath      string
	config          *config.Config
	logger          logger.Logger
	operationLogger *audit.OperationLogger
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewInitializer 创建初始化器
func NewInitializer(configPath string) *Initializer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Initializer{
		configPath: configPath,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Initialize 执行完整的项目初始化
func (i *Initializer) Initialize() error {
	log.Println("🚀 开始项目初始化...")

	// 1. 加载配置
	if err := i.loadConfig(); err != nil {
		return fmt.Errorf("配置加载失败: %w", err)
	}

	// 2. 初始化日志系统
	if err := i.initLogger(); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 3. 设置全局配置
	i.setupGlobals()

	// 4. 初始化数据库
	if err := i.initDatabase(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	// 5. 初始化Redis
	if err := i.initRedis(); err != nil {
		return fmt.Errorf("Redis初始化失败: %w", err)
	}

	// 6. 初始化其他服务
	if err := i.initServices(); err != nil {
		return fmt.Errorf("服务初始化失败: %w", err)
	}

	// 7. 初始化操作日志系统
	i.operationLogger = audit.NewOperationLogger(i.logger)

	log.Println("✅ 项目初始化完成")
	return nil
}

// loadConfig 加载项目配置
func (i *Initializer) loadConfig() error {
	cfg, err := config.Load(i.configPath)
	if err != nil {
		return err
	}

	i.config = cfg
	log.Printf("✅ 配置加载成功: %s", i.configPath)
	return nil
}

// initLogger 初始化日志系统
func (i *Initializer) initLogger() error {
	// 使用工厂模式创建日志实例
	factory, err := logger.GetFactory(i.config.Log.Type)
	if err != nil {
		return fmt.Errorf("获取日志工厂失败: %w", err)
	}

	// 转换配置类型
	logConfig := &logger.Config{
		Type:         i.config.Log.Type,
		Level:        i.config.Log.Level,
		Format:       i.config.Log.Format,
		Output:       i.config.Log.Output,
		LogDir:       i.config.Log.LogDir,
		MaxSize:      i.config.Log.MaxSize,
		MaxAge:       i.config.Log.MaxAge,
		MaxBackups:   i.config.Log.MaxBackups,
		Compress:     i.config.Log.Compress,
		LocalTime:    i.config.Log.LocalTime,
		EnableCaller: i.config.Log.EnableCaller,
		EnableTrace:  i.config.Log.EnableTrace,
	}

	l, err := factory.Create(logConfig)
	if err != nil {
		return fmt.Errorf("创建日志实例失败: %w", err)
	}

	i.logger = l
	log.Printf("✅ 日志系统初始化完成 (type: %s)", i.config.Log.Type)
	return nil
}

// setupGlobals 设置全局变量
func (i *Initializer) setupGlobals() {
	global.Config = i.config
	global.Logger = i.logger
	log.Println("✅ 全局配置设置完成")
}

// initDatabase 初始化数据库连接
func (i *Initializer) initDatabase() error {
	dbManager, err := database.NewDatabaseManager(i.config)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	global.DBManager = dbManager

	// 执行数据库迁移
	migrator := database.NewMigrator(dbManager.GetDB())
	if err := migrator.AutoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化基础数据
	if err := migrator.SeedData(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	log.Printf("📊 数据库初始化完成 (type: %s)", dbManager.GetType())
	return nil
}

// initRedis 初始化Redis连接
func (i *Initializer) initRedis() error {
	// TODO: 实现Redis初始化逻辑
	log.Println("🔄 Redis初始化完成")
	return nil
}

// initServices 初始化其他服务
func (i *Initializer) initServices() error {
	// TODO: 实现其他服务初始化
	log.Println("🔧 服务初始化完成")
	return nil
}

// Shutdown 优雅关闭
func (i *Initializer) Shutdown() {
	log.Println("🛑 开始优雅关闭...")

	// 取消上下文
	i.cancel()

	// 关闭日志
	if i.logger != nil {
		i.logger.Close()
	}

	log.Println("✅ 项目已优雅关闭")
}

// WaitForSignal 等待系统信号
func (i *Initializer) WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("📡 接收到关闭信号")
	i.Shutdown()
}
