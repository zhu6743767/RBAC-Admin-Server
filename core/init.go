package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"rbac.admin/global"
	"rbac.admin/models"
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
	if err := InitGorm(&global.Config.DB); err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}
	global.Logger.Info("✅ 数据库初始化成功")
	
	// 设置全局数据库连接实例
	global.DB = DB

	// 3. 初始化Redis缓存
	if err := InitRedis(); err != nil {
		return fmt.Errorf("初始化Redis失败: %v", err)
	}
	global.Logger.Info("✅ Redis初始化成功")

	// 4. 自动迁移数据库表结构
	if err := AutoMigrateModels(); err != nil {
		return fmt.Errorf("自动迁移数据库表结构失败: %v", err)
	}
	global.Logger.Info("✅ 数据库表结构自动迁移成功")

	// 5. 初始化Casbin权限管理
	if err := InitCasbin(); err != nil {
		return fmt.Errorf("初始化Casbin权限管理失败: %v", err)
	}
	global.Logger.Info("✅ Casbin权限管理初始化成功")

	global.Logger.Info("🎉 系统核心组件初始化完成")
	return nil
}

// InitCasbin 初始化Casbin权限管理
func InitCasbin() error {
	// 创建GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		return fmt.Errorf("创建Casbin适配器失败: %v", err)
	}

	// 加载模型配置文件
	modelPath := filepath.Join("config", "casbin", "model.conf")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		// 如果默认路径不存在，尝试使用替代路径
		modelPath = filepath.Join("../config", "casbin", "model.conf")
	}

	// 初始化Casbin执行器
	enforcer, err := casbin.NewCachedEnforcer(modelPath, adapter)
	if err != nil {
		return fmt.Errorf("初始化Casbin执行器失败: %v", err)
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("加载Casbin策略失败: %v", err)
	}

	// 设置全局Casbin实例
	global.Casbin = enforcer
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
