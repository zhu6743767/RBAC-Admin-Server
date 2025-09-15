// Package config 提供统一的配置结构定义
package config

import (
	"time"
)

// Config 全局配置结构体
// 包含RBAC管理员服务器的所有配置信息
// 通过YAML文件进行配置，支持热重载
// 该结构体与settings.yaml文件完全对应
// 修改后请确保settings.yaml格式同步更新

type Config struct {
	Server      ServerConfig      `yaml:"server"`      // 服务器配置
	CORS        CORSConfig        `yaml:"cors"`        // CORS跨域配置
	Database    DatabaseConfig    `yaml:"database"`    // 数据库配置
	JWT         JWTConfig         `yaml:"jwt"`         // JWT配置
	Redis       RedisConfig       `yaml:"redis"`       // Redis配置
	Log         LogConfig         `yaml:"log"`         // 日志配置
	Security    SecurityConfig    `yaml:"security"`    // 安全配置
	Performance PerformanceConfig `yaml:"performance"` // 性能配置
	Upload      UploadConfig      `yaml:"upload"`      // 文件上传配置
	Monitoring  MonitoringConfig  `yaml:"monitoring"`  // 监控配置
	Swagger     SwaggerConfig     `yaml:"swagger"`     // Swagger文档配置
}

// ServerConfig 服务器配置结构体
type ServerConfig struct {
	Port            int           `yaml:"port"`             // HTTP服务监听端口，默认8080
	Mode            string        `yaml:"mode"`             // 运行模式: debug(开发)、release(生产)、test(测试)
	ReadTimeout     time.Duration `yaml:"read_timeout"`     // 读取超时
	WriteTimeout    time.Duration `yaml:"write_timeout"`    // 写入超时
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"` // 优雅关闭超时
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Type            string        `yaml:"type"`              // 数据库类型: mysql, postgres, sqlite, sqlserver
	Driver          string        `yaml:"driver"`            // 数据库驱动（可选，按Type自动推导）
	Host            string        `yaml:"host"`              // 数据库服务器地址
	Port            int           `yaml:"port"`              // 数据库端口，MySQL默认3306
	Username        string        `yaml:"username"`          // 数据库用户名
	Password        string        `yaml:"password"`          // 数据库密码（建议使用环境变量）
	Database        string        `yaml:"database"`          // 数据库名称
	Path            string        `yaml:"path"`              // SQLite专用路径
	Charset         string        `yaml:"charset"`           // 字符集，推荐utf8mb4支持emoji
	Collation       string        `yaml:"collation"`         // 排序规则
	MaxOpenConns    int           `yaml:"max_open_conns"`    // 最大连接数，防止连接耗尽
	MaxIdleConns    int           `yaml:"max_idle_conns"`    // 空闲连接数，提高响应速度
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"` // 连接生命周期
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`// 空闲连接超时
	SSLMode         string        `yaml:"ssl_mode"`          // SSL模式
	Timeout         string        `yaml:"timeout"`           // 连接超时
}

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret             string `yaml:"secret"`               // JWT签名密钥，必须保密且足够复杂
	ExpireHours        int    `yaml:"expire_hours"`         // 令牌有效期，单位小时
	RefreshExpireHours int    `yaml:"refresh_expire_hours"` // 刷新令牌有效期，单位小时
	Issuer             string `yaml:"issuer"`               // 令牌颁发者
	Audience           string `yaml:"audience"`             // 令牌受众
}

// RedisConfig Redis配置结构体
type RedisConfig struct {
	Host         string        `yaml:"host"`                // Redis服务器地址
	Port         int           `yaml:"port"`                // Redis端口，默认6379
	Password     string        `yaml:"password"`            // Redis密码（如果没有密码留空）
	DB           int           `yaml:"db"`                  // 数据库编号，0-15
	MaxRetries   int           `yaml:"max_retries"`         // 最大重试次数
	PoolSize     int           `yaml:"pool_size"`           // 连接池大小
	MinIdleConns int           `yaml:"min_idle_conns"`      // 最小空闲连接数
	MaxConnAge   time.Duration `yaml:"max_conn_age"`        // 连接最大生命周期
	DialTimeout  time.Duration `yaml:"dial_timeout"`        // 连接超时
	ReadTimeout  time.Duration `yaml:"read_timeout"`        // 读取超时
	WriteTimeout time.Duration `yaml:"write_timeout"`       // 写入超时
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Type         string `yaml:"type"`           // 日志类型: logrus, zap, zerolog (默认logrus)
	Level        string `yaml:"level"`          // 日志级别: debug, info, warn, error
	Format       string `yaml:"format"`         // 日志格式: json(结构化), text(可读性)
	Output       string `yaml:"output"`         // 输出位置: stdout(控制台), file(文件), both(两者)
	LogDir       string `yaml:"log_dir"`        // 日志根目录
	MaxSize      int    `yaml:"max_size"`       // 日志文件最大大小(MB)
	MaxAge       int    `yaml:"max_age"`        // 日志文件最大保留天数
	MaxBackups   int    `yaml:"max_backups"`    // 日志文件最大备份数
	Compress     bool   `yaml:"compress"`       // 是否压缩旧日志
	LocalTime    bool   `yaml:"local_time"`     // 是否使用本地时间
	EnableCaller bool   `yaml:"enable_caller"`  // 记录调用者信息
	EnableTrace  bool   `yaml:"enable_trace"`   // 记录堆栈跟踪
}

// SecurityConfig 安全配置结构体
type SecurityConfig struct {
	BcryptCost          int           `yaml:"bcrypt_cost"`           // 密码加密强度，值越大越安全但越慢
	MaxLoginAttempts    int           `yaml:"max_login_attempts"`    // 最大登录尝试次数，防暴力破解
	LockDurationMinutes int           `yaml:"lock_duration_minutes"` // 账户锁定时间，单位分钟
	SessionTimeout      time.Duration `yaml:"session_timeout"`       // 会话超时时间
	APIKeyHeader        string        `yaml:"api_key_header"`        // API密钥请求头
	EnableCSRF          bool          `yaml:"enable_csrf"`           // 是否启用CSRF保护
	CSRFSecret          string        `yaml:"csrf_secret"`           // CSRF密钥
}

// CORSConfig CORS跨域配置结构体
type CORSConfig struct {
	Enable           bool          `yaml:"enable"`             // 是否启用CORS
	AllowOrigins     []string      `yaml:"allow_origins"`      // 允许的源
	AllowMethods     []string      `yaml:"allow_methods"`      // 允许的方法
	AllowHeaders     []string      `yaml:"allow_headers"`      // 允许的请求头
	ExposeHeaders    []string      `yaml:"expose_headers"`     // 暴露的响应头
	AllowCredentials bool          `yaml:"allow_credentials"`  // 允许携带凭证
	MaxAge           time.Duration `yaml:"max_age"`            // 预检缓存时间
}

// PerformanceConfig 性能配置结构体
type PerformanceConfig struct {
	EnablePprof        bool   `yaml:"enable_pprof"`         // 性能分析工具
	MaxUploadSize      string `yaml:"max_upload_size"`       // 最大上传文件大小
	RequestRateLimit   int    `yaml:"request_rate_limit"`    // 请求频率限制
	BurstRateLimit     int    `yaml:"burst_rate_limit"`      // 突发请求限制
	EnableCompression  bool   `yaml:"enable_compression"`    // 启用响应压缩
	CompressionLevel   int    `yaml:"compression_level"`     // 压缩级别1-9
}

// UploadConfig 文件上传配置结构体
type UploadConfig struct {
	MaxFileSize        string   `yaml:"max_file_size"`         // 最大文件大小
	AllowedExtensions  []string `yaml:"allowed_extensions"`    // 允许的文件扩展名
	AllowedMimeTypes   []string `yaml:"allowed_mime_types"`    // 允许的MIME类型
	StorageType        string   `yaml:"storage_type"`          // 存储类型
	StoragePath        string   `yaml:"storage_path"`          // 本地存储路径
	MaxFilesPerRequest int      `yaml:"max_files_per_request"` // 单次请求最大文件数
}

// MonitoringConfig 监控配置结构体
type MonitoringConfig struct {
	EnableHealthCheck bool   `yaml:"enable_health_check"` // 启用健康检查
	HealthCheckPath   string `yaml:"health_check_path"`   // 健康检查路径
	EnableMetrics     bool   `yaml:"enable_metrics"`      // 启用指标收集
	MetricsPath       string `yaml:"metrics_path"`        // 指标路径
	EnableTracing     bool   `yaml:"enable_tracing"`      // 启用链路追踪
	TracingEndpoint   string `yaml:"tracing_endpoint"`    // 追踪服务端点
}

// SwaggerConfig Swagger文档配置结构体
type SwaggerConfig struct {
	Enable      bool   `yaml:"enable"`      // 是否启用Swagger
	Title       string `yaml:"title"`       // API标题
	Description string `yaml:"description"` // API描述
	Version     string `yaml:"version"`     // API版本
	Host        string `yaml:"host"`        // API主机
	BasePath    string `yaml:"base_path"`   // API基础路径
	EnableUI    bool   `yaml:"enable_ui"`   // 是否启用Swagger UI
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            8080,
			Mode:            "debug",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Type:            "sqlite",
			Host:            "localhost",
			Port:            3306,
			Database:        "rbac_admin",
			Username:        "root",
			Password:        "",
			Charset:         "utf8mb4",
			SSLMode:         "disable",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 3600 * time.Second,
			ConnMaxIdleTime: 300 * time.Second,
		},
		Redis: RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			PoolSize:     100,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		Log: LogConfig{
			Type:         "logrus",
			Level:        "info",
			Format:       "json",
			Output:       "both",
			LogDir:       "logs",
			MaxSize:      100,
			MaxAge:       7,
			MaxBackups:   10,
			Compress:     true,
			LocalTime:    true,
			EnableCaller: true,
			EnableTrace:  false,
		},
		JWT: JWTConfig{
			Secret:             "your-secret-key-change-this",
			ExpireHours:        24,
			RefreshExpireHours: 168,
			Issuer:             "rbac-admin",
			Audience:           "rbac-admin",
		},
		Security: SecurityConfig{
			BcryptCost:          10,
			MaxLoginAttempts:    5,
			LockDurationMinutes: 30,
			SessionTimeout:      24 * time.Hour,
			APIKeyHeader:        "X-API-Key",
			EnableCSRF:          true,
			CSRFSecret:          "csrf-secret-key",
		},
		CORS: CORSConfig{
			Enable:           true,
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
		Performance: PerformanceConfig{
			EnablePprof:        false,
			MaxUploadSize:      "32MB",
			RequestRateLimit:   100,
			BurstRateLimit:     50,
			EnableCompression:  true,
			CompressionLevel:   6,
		},
		Upload: UploadConfig{
			MaxFileSize:        "10MB",
			AllowedExtensions:  []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"},
			AllowedMimeTypes:   []string{"image/jpeg", "image/png", "image/gif", "application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
			StorageType:        "local",
			StoragePath:        "./uploads",
			MaxFilesPerRequest: 5,
		},
		Monitoring: MonitoringConfig{
			EnableHealthCheck: true,
			HealthCheckPath:   "/health",
			EnableMetrics:     true,
			MetricsPath:       "/metrics",
			EnableTracing:     false,
			TracingEndpoint:   "http://localhost:14268/api/traces",
		},
		Swagger: SwaggerConfig{
			Enable:      true,
			Title:       "RBAC Admin API",
			Description: "基于角色的访问控制管理员系统API文档",
			Version:     "1.0.0",
			Host:        "localhost:8080",
			BasePath:    "/api/v1",
			EnableUI:    true,
		},
	}
}
