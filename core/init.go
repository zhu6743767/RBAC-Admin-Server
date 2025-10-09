package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"rbac_admin_server/core/init_casbin"
	"rbac_admin_server/core/init_gorm"
	"rbac_admin_server/core/init_redis"
	"rbac_admin_server/global"
)

// InitSystem 初始化系统核心组件
// 按顺序初始化：验证器 -> 数据库 -> Redis -> 数据库表迁移 -> Casbin权限管理
func InitSystem() error {
	global.Logger.Info("开始初始化系统核心组件...")

	// 1. 初始化验证器
	if err := InitValidator(); err != nil {
		return fmt.Errorf("初始化验证器失败: %v", err)
	}
	global.Logger.Info("✅ 验证器初始化成功")

	// 2. 初始化数据库连接
	db, err := init_gorm.InitGorm()
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}
	global.DB = db
	global.Logger.Info("✅ 数据库初始化成功")

	// 3. 初始化Redis缓存
	redisClient, err := init_redis.InitRedis()
	if err != nil {
		global.Logger.Warnf("⚠️ Redis初始化失败: %v, 将继续运行", err)
	} else {
		global.Redis = redisClient
		global.Logger.Info("✅ Redis初始化成功")
	}

	// 4. 自动迁移数据库表结构
	if err := AutoMigrateModels(); err != nil {
		return fmt.Errorf("自动迁移数据库表结构失败: %v", err)
	}
	global.Logger.Info("✅ 数据库表结构自动迁移成功")

	// 5. 初始化Casbin权限管理
	casbinEnforcer, err := init_casbin.InitCasbin()
	if err != nil {
		global.Logger.Warnf("⚠️ Casbin权限管理初始化失败: %v, 将继续运行", err)
	} else {
		global.Casbin = casbinEnforcer
		global.Logger.Info("✅ Casbin权限管理初始化成功")
	}

	global.Logger.Info("🎉 系统核心组件初始化完成")
	return nil
}

// InitCasbin 初始化Casbin权限管理（废弃，请使用init_casbin包中的InitCasbin函数）
func InitCasbin() error {
	global.Logger.Warnf("⚠️ core/InitCasbin已废弃，请使用init_casbin包中的InitCasbin函数")
	return fmt.Errorf("该函数已废弃，请使用init_casbin包中的InitCasbin函数")
}

// AutoMigrateModels 自动迁移所有模型
func AutoMigrateModels() error {
	return init_gorm.MigrateTables(global.DB)
}

// CleanupSystem 清理系统资源
func CleanupSystem() {
	global.Logger.Info("开始清理系统资源...")

	// 关闭数据库连接
	if err := CloseDB(); err != nil {
		global.Logger.Errorf("关闭数据库连接失败: %v", err)
	} else {
		global.Logger.Info("✅ 数据库连接已关闭")
	}

	// 关闭Redis连接
	if err := CloseRedis(); err != nil {
		global.Logger.Errorf("关闭Redis连接失败: %v", err)
	} else {
		global.Logger.Info("✅ Redis连接已关闭")
	}

	global.Logger.Info("✅ 系统资源清理完成")
}

// WaitForSignal 等待系统信号
func WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	global.Logger.Info("等待系统信号...")
	sig := <-sigChan
	global.Logger.Infof("接收到信号: %v，开始优雅关闭...", sig)

	// 清理系统资源
	CleanupSystem()

	global.Logger.Info("系统已优雅关闭")
}

// GetSystemStatus 获取系统状态
func GetSystemStatus() map[string]interface{} {
	status := make(map[string]interface{})

	// 数据库状态
	if global.DB != nil {
		sqlDB, err := global.DB.DB()
		if err == nil {
			if err := sqlDB.Ping(); err != nil {
				status["db"] = map[string]interface{}{
					"status": "error",
					"error":	 err.Error(),
				}
			} else {
				status["db"] = map[string]interface{}{
					"status": "connected",
				}
			}
		} else {
			status["db"] = map[string]interface{}{
				"status": "error",
				"error":	 err.Error(),
			}
		}
	} else {
		status["db"] = map[string]interface{}{
			"status": "not_initialized",
		}
	}

	// Redis状态
	if global.Redis != nil {
		if err := global.Redis.Ping(RedisCtx).Err(); err != nil {
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
