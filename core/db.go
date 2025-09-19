package core

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"rbac.admin/config"
	"rbac.admin/global"

	// 纯Go SQLite驱动，无需CGO
	glebarezsqlite "github.com/glebarez/sqlite"
)

var (
	// DB 全局数据库连接
	DB *gorm.DB
	// SQLDB 全局SQL数据库连接
	SQLDB *sql.DB
)

// InitGorm 初始化GORM数据库连接
// 支持MySQL、PostgreSQL、SQLite三种数据库
// 自动配置连接池参数，支持日志级别设置
func InitGorm(cfg *config.DBConfig) error {
	var dialector gorm.Dialector

	// 验证配置
	if cfg == nil {
		return fmt.Errorf("数据库配置为空")
	}

	// 根据数据库类型选择驱动
	switch cfg.Mode {
	case "mysql":
		dialector = mysql.Open(buildMysqlDSN(cfg))
	case "pgsql", "postgres", "postgresql":
		dialector = postgres.Open(buildPostgresDSN(cfg))
	case "sqlite":
		// 使用纯Go SQLite驱动，无需CGO
		dialector = glebarezsqlite.Open(buildSqliteDSN(cfg))
	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Mode)
	}

	// GORM配置
	gormConfig := &gorm.Config{
		// 命名策略（使用单数表名）
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 日志级别（默认Info级别）
	gormConfig.Logger = logger.Default.LogMode(logger.Info)

	// 打开数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 获取底层SQL连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取SQL连接失败: %v", err)
	}

	// 设置连接池参数（使用配置文件中的值）
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	DB = db
	SQLDB = sqlDB

	// 日志输出连接信息（使用全局日志系统）
	if cfg.Mode == "sqlite" {
		global.Logger.Infof("✅ 数据库连接成功: SQLite @ %s", cfg.Path)
	} else {
		global.Logger.Infof("✅ 数据库连接成功: %s@%s:%d/%s", cfg.User, cfg.Host, cfg.Port, cfg.DbNAME)
	}

	// 记录连接池配置
	global.Logger.Infof("🔧 数据库连接池配置: MaxIdleConns=%d, MaxOpenConns=%d, MaxLifetime=%v, MaxIdleTime=%v",
		cfg.MaxIdleConns, cfg.MaxOpenConns, cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime)

	return nil
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

	// 添加SSL模式
	if cfg.SSLMode != "" {
		dsn += fmt.Sprintf("&tls=%s", cfg.SSLMode)
	}

	// 添加连接超时
	if cfg.Timeout != "" {
		dsn += fmt.Sprintf("&timeout=%s", cfg.Timeout)
	}

	return dsn
}

// buildPostgresDSN 构建PostgreSQL连接字符串
func buildPostgresDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password, // 使用正确的字段名
		cfg.DbNAME,   // 使用正确的字段名
		cfg.Port,
		cfg.SSLMode)

	// 添加连接超时
	if cfg.Timeout != "" {
		dsn += fmt.Sprintf(" connect_timeout=%s", cfg.Timeout)
	}

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
	if SQLDB != nil {
		return SQLDB.Close()
	}
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
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
