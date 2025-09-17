package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"rbac.admin/config"
	"rbac.admin/models"
)

// InitSystem 初始化系统
// 按顺序初始化配置、日志、验证器、数据库、Redis
func InitSystem(cfg *config.Config) error {
	logrus.Info("开始初始化系统...")

	// 1. 初始化日志系统
	if err := InitLogger(&cfg.Log); err != nil {
		return fmt.Errorf("初始化日志系统失败: %v", err)
	}
	logrus.Info("✅ 日志系统初始化成功")

	// 2. 初始化验证器
	if err := InitValidator(); err != nil {
		return fmt.Errorf("初始化验证器失败: %v", err)
	}
	logrus.Info("✅ 验证器初始化成功")

	// 3. 初始化数据库
	if err := InitGorm(&cfg.DB); err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}
	logrus.Info("✅ 数据库初始化成功")

	// 4. 初始化Redis
	if err := InitRedis(&cfg.Redis); err != nil {
		return fmt.Errorf("初始化Redis失败: %v", err)
	}
	logrus.Info("✅ Redis初始化成功")

	// 5. 自动迁移数据库表结构
	if err := AutoMigrateModels(); err != nil {
		return fmt.Errorf("自动迁移数据库表结构失败: %v", err)
	}
	logrus.Info("✅ 数据库表结构自动迁移成功")

	logrus.Info("🎉 系统初始化完成")
	return nil
}

// AutoMigrateModels 自动迁移所有模型
func AutoMigrateModels() error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 需要迁移的模型列表
	models := []interface{}{
		// 基础模型
		&models.User{},
		&models.Department{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},

		// 菜单模型
		&models.Menu{},
		&models.RoleMenu{},

		// API和配置模型
		&models.API{},
		&models.Dict{},
		&models.DictItem{},
		&models.Config{},

		// 文件和日志模型
		&models.File{},
		&models.Log{},
	}

	// 执行自动迁移
	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型失败: %v", err)
		}
	}

	return nil
}

// CleanupSystem 清理系统资源
func CleanupSystem() {
	logrus.Info("开始清理系统资源...")

	// 关闭数据库连接
	if err := CloseDB(); err != nil {
		logrus.Errorf("关闭数据库连接失败: %v", err)
	} else {
		logrus.Info("✅ 数据库连接已关闭")
	}

	// 关闭Redis连接
	if err := CloseRedis(); err != nil {
		logrus.Errorf("关闭Redis连接失败: %v", err)
	} else {
		logrus.Info("✅ Redis连接已关闭")
	}

	logrus.Info("✅ 系统资源清理完成")
}

// WaitForSignal 等待系统信号
func WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logrus.Info("等待系统信号...")
	sig := <-sigChan
	logrus.Infof("接收到信号: %v，开始优雅关闭...", sig)

	// 清理系统资源
	CleanupSystem()

	logrus.Info("系统已优雅关闭")
}

// GetSystemStatus 获取系统状态
func GetSystemStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// 数据库状态
	if DB != nil && SQLDB != nil {
		if err := SQLDB.Ping(); err != nil {
			status["db"] = map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
		} else {
			status["db"] = map[string]interface{}{
				"status": "connected",
			}
		}
	} else {
		status["db"] = map[string]interface{}{
			"status": "not_initialized",
		}
	}

	// Redis状态
	if RedisClient != nil {
		if err := RedisClient.Ping(RedisCtx).Err(); err != nil {
			status["redis"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			status["redis"] = map[string]interface{}{
				"status": "connected",
			}
		}
	} else {
		status["redis"] = map[string]interface{}{
			"status": "not_initialized",
		}
	}

	return status
}
