package core

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"rbac.admin/config"

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

	// 设置连接池参数（使用默认值）
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	DB = db
	SQLDB = sqlDB

	fmt.Printf("✅ 数据库连接成功: %s@%s:%d/%s\n", cfg.User, cfg.Host, cfg.Port, cfg.DbNAME)
	return nil
}

// buildMysqlDSN 构建MySQL连接字符串
func buildMysqlDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.PASSWORD,
		cfg.Host,
		cfg.Port,
		cfg.DbNAME,
		"utf8mb4")
	// 默认使用utf8mb4_general_ci排序规则
	dsn += "&collation=utf8mb4_general_ci"

	return dsn
}

// buildPostgresDSN 构建PostgreSQL连接字符串
func buildPostgresDSN(cfg *config.DBConfig) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.PASSWORD,
		cfg.DbNAME,
		cfg.Port)

	// 添加额外参数
	dsn += " client_encoding=utf8"

	return dsn
}

// buildSqliteDSN 构建SQLite连接字符串
func buildSqliteDSN(cfg *config.DBConfig) string {
	// 如果配置了路径则使用路径，否则使用内存数据库
	if cfg.Path != "" {
		return cfg.Path
	}
	// 如果配置了数据库名称则使用数据库名称作为文件名
	if cfg.DbNAME != "" && cfg.DbNAME != ":memory:" {
		return cfg.DbNAME + ".db"
	}
	// 默认使用内存数据库
	return ":memory:"
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
