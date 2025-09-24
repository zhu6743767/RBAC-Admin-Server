package init_gorm

import (
	"fmt"
	"rbac_admin_server/config"
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitGorm 初始化GORM数据库连接
func InitGorm() (*gorm.DB, error) {
	config := global.Config.DB

	// 配置日志级别
	logLevel := logger.Info
	if global.Config.System.Mode == "release" {
		logLevel = logger.Warn
	}

	// 创建新的日志器
	newLogger := logger.New(
		global.Logger, // 使用全局日志实例
		logger.Config{
			SlowThreshold: 200, // 慢SQL阈值
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	// 通用配置
	dbConfig := &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	var db *gorm.DB
	var err error

	// 根据数据库类型初始化连接
	switch config.Mode {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)
		db, err = gorm.Open(mysql.Open(dsn), dbConfig)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			config.Host, config.User, config.Password, config.DBName, config.Port)
		db, err = gorm.Open(postgres.Open(dsn), dbConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DBName), dbConfig)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Mode)
	}

	if err != nil {
		return nil, err
	}

	// 获取数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	logrus.Info("✅ 数据库连接初始化成功")
	return db, nil
}

// MigrateTables 迁移数据表
func MigrateTables(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 定义要迁移的数据表模型
	tables := []interface{}{
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

	// 执行迁移
	if err := db.AutoMigrate(tables...); err != nil {
		return fmt.Errorf("迁移模型失败: %v", err)
	}

	logrus.Info("✅ 数据库表迁移成功")
	return nil
}