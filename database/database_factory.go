package database

import (
	"fmt"
	"rbac.admin/config"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DatabaseType 数据库类型枚举
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite"
	SQLServer  DatabaseType = "sqlserver"
)

// DatabaseConfig 数据库配置接口
type DatabaseConfig interface {
	GetDriver() string
	GetDSN() string
	GetConfig() *gorm.Config
}

// BaseConfig 基础数据库配置
type BaseConfig struct {
	Type     DatabaseType `yaml:"type"`
	Host     string       `yaml:"host"`
	Port     int          `yaml:"port"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password"`
	Database string       `yaml:"database"`
	Charset  string       `yaml:"charset"`
	SSLMode  string       `yaml:"ssl_mode"`
	Timeout  string       `yaml:"timeout"`
	Path     string       `yaml:"path"` // SQLite专用
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	BaseConfig
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// PostgreSQLConfig PostgreSQL配置
type PostgreSQLConfig struct {
	BaseConfig
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// SQLiteConfig SQLite配置
type SQLiteConfig struct {
	BaseConfig
	Cache    string `yaml:"cache"`
	Mode     string `yaml:"mode"`
	Journal  string `yaml:"journal"`
}

// SQLServerConfig SQLServer配置
type SQLServerConfig struct {
	BaseConfig
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// DatabaseFactory 数据库工厂
type DatabaseFactory struct {
	configs map[DatabaseType]DatabaseConfig
}

// NewDatabaseFactory 创建数据库工厂
func NewDatabaseFactory() *DatabaseFactory {
	return &DatabaseFactory{
		configs: make(map[DatabaseType]DatabaseConfig),
	}
}

// Register 注册数据库配置
func (f *DatabaseFactory) Register(dbType DatabaseType, config DatabaseConfig) {
	f.configs[dbType] = config
}

// CreateConnection 创建数据库连接
func (f *DatabaseFactory) CreateConnection(dbType DatabaseType) (*gorm.DB, error) {
	config, exists := f.configs[dbType]
	if !exists {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	dsn := config.GetDSN()
	gormConfig := config.GetConfig()

	db, err := gorm.Open(getDialector(dbType, dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", dbType, err)
	}

	// 设置连接池参数
	if sqlDB, err := db.DB(); err == nil {
		switch dbType {
		case MySQL:
			if mysqlConfig, ok := config.(*MySQLConfig); ok {
				sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
				sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
			}
		case PostgreSQL:
			if pgConfig, ok := config.(*PostgreSQLConfig); ok {
				sqlDB.SetMaxIdleConns(pgConfig.MaxIdleConns)
				sqlDB.SetMaxOpenConns(pgConfig.MaxOpenConns)
			}
		case SQLServer:
			if sqlConfig, ok := config.(*SQLServerConfig); ok {
				sqlDB.SetMaxIdleConns(sqlConfig.MaxIdleConns)
				sqlDB.SetMaxOpenConns(sqlConfig.MaxOpenConns)
			}
		}
	}

	return db, nil
}

// getDialector 获取数据库驱动
func getDialector(dbType DatabaseType, dsn string) gorm.Dialector {
	switch dbType {
	case MySQL:
		return mysql.Open(dsn)
	case PostgreSQL:
		return postgres.Open(dsn)
	case SQLite:
		return sqlite.Open(dsn)
	case SQLServer:
		return sqlserver.Open(dsn)
	default:
		return mysql.Open(dsn) // 默认使用MySQL
	}
}

// GetDSN 实现接口方法 (MySQLConfig)
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetConfig 实现接口方法 (MySQLConfig)
func (c *MySQLConfig) GetConfig() *gorm.Config {
	return &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// GetDSN 实现接口方法 (PostgreSQLConfig)
func (c *PostgreSQLConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
}

// GetConfig 实现接口方法 (PostgreSQLConfig)
func (c *PostgreSQLConfig) GetConfig() *gorm.Config {
	return &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// GetDSN 实现接口方法 (SQLiteConfig)
func (c *SQLiteConfig) GetDSN() string {
	return fmt.Sprintf("%s?cache=%s&mode=%s&journal=%s",
		c.Path, c.Cache, c.Mode, c.Journal)
}

// GetConfig 实现接口方法 (SQLiteConfig)
func (c *SQLiteConfig) GetConfig() *gorm.Config {
	return &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// GetDSN 实现接口方法 (SQLServerConfig)
func (c *SQLServerConfig) GetDSN() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

// GetConfig 实现接口方法 (SQLServerConfig)
func (c *SQLServerConfig) GetConfig() *gorm.Config {
	return &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// GetDriver 实现接口方法
func (c *BaseConfig) GetDriver() string {
	return string(c.Type)
}

// AutoMigrate 自动迁移数据库结构
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	factory    *DatabaseFactory
	db         *gorm.DB
	DBType     DatabaseType
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(config *config.Config) (*DatabaseManager, error) {
	factory := NewDatabaseFactory()

	// 根据配置注册对应的数据库类型
	var dbType DatabaseType
	switch strings.ToLower(config.Database.Type) {
	case "mysql":
		dbType = MySQL
		mysqlConfig := &MySQLConfig{
			BaseConfig: BaseConfig{
				Type:     MySQL,
				Host:     config.Database.Host,
				Port:     config.Database.Port,
				Username: config.Database.Username,
				Password: config.Database.Password,
				Database: config.Database.Database,
				Charset:  config.Database.Charset,
				SSLMode:  config.Database.SSLMode,
			},
			MaxIdleConns:    config.Database.MaxIdleConns,
			MaxOpenConns:    config.Database.MaxOpenConns,
			ConnMaxLifetime: config.Database.ConnMaxLifetime,
		}
		factory.Register(dbType, mysqlConfig)
	case "postgres":
		dbType = PostgreSQL
		pgConfig := &PostgreSQLConfig{
			BaseConfig: BaseConfig{
				Type:     PostgreSQL,
				Host:     config.Database.Host,
				Port:     config.Database.Port,
				Username: config.Database.Username,
				Password: config.Database.Password,
				Database: config.Database.Database,
				SSLMode:  config.Database.SSLMode,
			},
			MaxIdleConns:    config.Database.MaxIdleConns,
			MaxOpenConns:    config.Database.MaxOpenConns,
			ConnMaxLifetime: config.Database.ConnMaxLifetime,
		}
		factory.Register(dbType, pgConfig)
	case "sqlite":
		dbType = SQLite
		sqliteConfig := &SQLiteConfig{
			BaseConfig: BaseConfig{
				Type: SQLite,
				Path: config.Database.Path,
			},
			Cache:   "shared",
			Mode:    "memory",
			Journal: "wal",
		}
		factory.Register(dbType, sqliteConfig)
	case "sqlserver":
		dbType = SQLServer
		sqlConfig := &SQLServerConfig{
			BaseConfig: BaseConfig{
				Type:     SQLServer,
				Host:     config.Database.Host,
				Port:     config.Database.Port,
				Username: config.Database.Username,
				Password: config.Database.Password,
				Database: config.Database.Database,
				SSLMode:  config.Database.SSLMode,
			},
			MaxIdleConns:    config.Database.MaxIdleConns,
			MaxOpenConns:    config.Database.MaxOpenConns,
			ConnMaxLifetime: config.Database.ConnMaxLifetime,
		}
		factory.Register(dbType, sqlConfig)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Database.Type)
	}

	db, err := factory.CreateConnection(dbType)
	if err != nil {
		return nil, err
	}

	// TODO: 使用日志系统记录数据库连接成功
	fmt.Printf("数据库连接成功: type=%s, host=%s, port=%d, database=%s\n",
		dbType, config.Database.Host, config.Database.Port, config.Database.Database)

	return &DatabaseManager{
		factory: factory,
		db:      db,
		DBType:  dbType,
	}, nil
}

// GetDB 获取数据库连接
func (m *DatabaseManager) GetDB() *gorm.DB {
	return m.db
}

// GetType 获取数据库类型
func (m *DatabaseManager) GetType() DatabaseType {
	return m.DBType
}

// Close 关闭数据库连接
func (m *DatabaseManager) Close() error {
	if sqlDB, err := m.db.DB(); err == nil {
		return sqlDB.Close()
	}
	return nil
}

// Ping 检查数据库连接
func (m *DatabaseManager) Ping() error {
	if sqlDB, err := m.db.DB(); err == nil {
		return sqlDB.Ping()
	}
	return fmt.Errorf("failed to get sql.DB")
}