// Package config 提供配置管理功能
package config

import (
	"time"
)

// Config 主配置结构体
type Config struct {
	System      SystemConfig      `mapstructure:"system" json:"system" yaml:"system"`
	DB          DBConfig          `mapstructure:"db" json:"db" yaml:"db"`
	Redis       RedisConfig       `mapstructure:"redis" json:"redis" yaml:"redis"`
	Log         LogConfig         `mapstructure:"log" json:"log" yaml:"log"`
	JWT         JWTConfig         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	CORS        CORSConfig        `mapstructure:"cors" json:"cors" yaml:"cors"`
	Security    SecurityConfig    `mapstructure:"security" json:"security" yaml:"security"`
	Performance PerformanceConfig `mapstructure:"performance" json:"performance" yaml:"performance"`
	Upload      UploadConfig      `mapstructure:"upload" json:"upload" yaml:"upload"`
	Monitoring  MonitoringConfig  `mapstructure:"monitoring" json:"monitoring" yaml:"monitoring"`
	Swagger     SwaggerConfig     `mapstructure:"swagger" json:"swagger" yaml:"swagger"`
	App         AppConfig         `mapstructure:"app" json:"app" yaml:"app"`
}

// SystemConfig 系统配置结构体
type SystemConfig struct {
	IP   string `mapstructure:"ip" json:"ip" yaml:"ip"`     // 监听IP
	Port int    `mapstructure:"port" json:"port" yaml:"port"` // HTTP服务监听端口，默认8080
}

// DBConfig 数据库配置结构体
type DBConfig struct {
	Mode            string        `mapstructure:"mode" json:"mode" yaml:"mode"`                        // 数据库类型: mysql, postgres, sqlite, sqlserver
	Driver          string        `mapstructure:"driver" json:"driver" yaml:"driver"`                  // 数据库驱动（可选，按Mode自动推导）
	Host            string        `mapstructure:"host" json:"host" yaml:"host"`                        // 数据库服务器地址
	Port            int           `mapstructure:"port" json:"port" yaml:"port"`                          // 数据库端口，MySQL默认3306
	User            string        `mapstructure:"user" json:"user" yaml:"user"`                          // 数据库用户名
	PASSWORD        string        `mapstructure:"password" json:"password" yaml:"password"`              // 数据库密码（建议使用环境变量）
	DbNAME          string        `mapstructure:"dbname" json:"dbname" yaml:"dbname"`                  // 数据库名称
	Path            string        `mapstructure:"path" json:"path" yaml:"path"`                          // SQLite专用路径
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`    // 最大连接数，防止连接耗尽
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`    // 空闲连接数，提高响应速度
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 连接生命周期
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 空闲连接超时
	SSLMode         string        `mapstructure:"ssl_mode" json:"ssl_mode" yaml:"ssl_mode"`              // SSL模式
	Timeout         string        `mapstructure:"timeout" json:"timeout" yaml:"timeout"`                  // 连接超时
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
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`         // Redis地址，格式 host:port
	Password string `mapstructure:"password" json:"password" yaml:"password"` // Redis密码（如果没有密码留空）
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                 // 数据库编号，0-15
}

// LogConfig 日志配置结构体
type LogConfig struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`
	Format       string `mapstructure:"format" json:"format" yaml:"format"`
	Output       string `mapstructure:"output" json:"output" yaml:"output"`
	LogDir       string `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
	MaxSize      int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge       int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Compress     bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
	EnableCaller bool   `mapstructure:"enable_caller" json:"enable_caller" yaml:"enable_caller"`
	EnableTrace  bool   `mapstructure:"enable_trace" json:"enable_trace" yaml:"enable_trace"`
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

// AppConfig 应用配置结构体
type AppConfig struct {
	Name        string `yaml:"name"`        // 应用名称
	Version     string `yaml:"version"`     // 应用版本
	Environment string `yaml:"environment"` // 应用环境
	Debug       bool   `yaml:"debug"`       // 是否调试模式
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		System: SystemConfig{
			IP:   "127.0.0.1",
			Port: 8080,
		},
		DB: DBConfig{
			Mode:            "sqlite",
			Host:            "localhost",
			Port:            3306,
			User:            "root",
			PASSWORD:        "",
			DbNAME:          "rbac_admin",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 1 * time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
			SSLMode:         "disable",
			Timeout:         "30s",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:             "your-jwt-secret-key-here",
			ExpireHours:        24,
			RefreshExpireHours: 168,
			Issuer:             "rbac-admin-server",
			Audience:           "rbac-client",
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
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
			ExposeHeaders:    []string{"X-Total-Count"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
		Log: LogConfig{
			Level:        "info",
			Format:       "text",
			Output:       "both",
			LogDir:       "./logs",
			MaxSize:      100,
			MaxBackups:   3,
			MaxAge:       7,
			Compress:     true,
			EnableCaller: true,
			EnableTrace:  false,
		},
		Performance: PerformanceConfig{
			EnablePprof:        false,
			MaxUploadSize:      "10MB",
			RequestRateLimit:   100,
			BurstRateLimit:     200,
			EnableCompression:  true,
			CompressionLevel:   6,
		},
		Upload: UploadConfig{
			MaxFileSize:         "10MB",
			AllowedExtensions:   []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"},
			AllowedMimeTypes:    []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
			StorageType:         "local",
			StoragePath:         "./uploads",
			MaxFilesPerRequest:  5,
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
			Title:       "RBAC管理员API",
			Description: "企业级RBAC权限管理API",
			Version:     "1.0.0",
			Host:        "localhost:8080",
			BasePath:    "/api/v1",
			EnableUI:    true,
		},
		App: AppConfig{
			Name:        "RBAC管理员服务器",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
	}
}
