package core

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"rbac_admin_server/config"
	"rbac_admin_server/global"
)

var (
	// DB 全局数据库连接
	DB *gorm.DB
	// SQLDB 全局SQL数据库连接
	SQLDB *sql.DB
)

// InitGorm 初始化GORM数据库连接（废弃，请使用init_gorm包中的InitGorm函数）
// 支持MySQL、PostgreSQL、SQLite三种数据库
// 自动配置连接池参数，支持日志级别设置
func InitGorm(cfg *config.DBConfig) error {
	global.Logger.Warnf("⚠️ core/InitGorm已废弃，请使用init_gorm包中的InitGorm函数")
	return fmt.Errorf("该函数已废弃，请使用init_gorm包中的InitGorm函数")
}

// buildMysqlDSN 构建MySQL连接字符串
func buildMysqlDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password, // 使用正确的字段名
		cfg.Host,
		cfg.Port,
		cfg.DbNAME,   // 使用正确的字段名
		"utf8mb4")
	// 默认使用utf8mb4_general_ci排序规则
	dsn += "&collation=utf8mb4_general_ci"

	// 添加SSL模式（使用默认设置）
	dsn += "&tls=false"

	return dsn
}

// buildPostgresDSN 构建PostgreSQL连接字符串
func buildPostgresDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password, // 使用正确的字段名
		cfg.DbNAME,   // 使用正确的字段名
		cfg.Port)

	return dsn
}

// buildSqliteDSN 构建SQLite连接字符串
func buildSqliteDSN(cfg *config.DBConfig) string {
	path := cfg.Path
	if path == "" {
		path = cfg.DbNAME + ".db"
	}

	// SQLite连接参数
	// _pragma=foreign_keys=on 启用外键约束
	// cache=shared 启用共享缓存
	// mode=rwc 创建和读写模式
	return fmt.Sprintf("%s?_pragma=foreign_keys=on&cache=shared&mode=rwc", path)
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if global.DB != nil {
		sqlDB, err := global.DB.DB()
		if err == nil {
			global.Logger.Info("🔄 正在关闭数据库连接...")
			return sqlDB.Close()
		}
		return err
	}
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return global.DB
}

// GetSQLDB 获取SQL数据库连接
func GetSQLDB() *sql.DB {
	return SQLDB
}



// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("数据库连接未初始化")
	}
	return DB.AutoMigrate(models...)
}

// IsRecordNotFound 判断是否为记录未找到错误
func IsRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
