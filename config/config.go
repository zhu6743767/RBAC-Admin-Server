package config

import (
	"time"
)

// Config 全局配置结构体
type Config struct {
	System      SystemConfig      `yaml:"system"`
	DB          DBConfig          `yaml:"db"`
	Redis       RedisConfig       `yaml:"redis"`
	Log         LogConfig         `yaml:"log"`
	JWT         JWTConfig         `yaml:"jwt"`
	Security    SecurityConfig    `yaml:"security"`
	CORS        CORSConfig        `yaml:"cors"`
	Performance PerformanceConfig `yaml:"performance"`
	Upload      UploadConfig      `yaml:"upload"`
	Monitoring  MonitoringConfig  `yaml:"monitoring"`
	Swagger     SwaggerConfig     `yaml:"swagger"`
	App         AppConfig         `yaml:"app"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	IP        string `yaml:"ip"`
	Port      int    `yaml:"port"`
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	Timezone  string `yaml:"timezone"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Mode              string        `yaml:"mode"`
	Host              string        `yaml:"host"`
	Port              int           `yaml:"port"`
	User              string        `yaml:"user"`
	Password          string        `yaml:"password"`
	DbNAME            string        `yaml:"dbname"`
	SSLMode           string        `yaml:"sslmode"`
	Timeout           string        `yaml:"timeout"`
	Charset           string        `yaml:"charset"`
	Collation         string        `yaml:"collation"`
	MaxIdleConns      int           `yaml:"max_idle_conns"`
	MaxOpenConns      int           `yaml:"max_open_conns"`
	ConnMaxLifetime   time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime   time.Duration `yaml:"conn_max_idle_time"`
	Path              string        `yaml:"path"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr               string        `yaml:"addr"`
	Password           string        `yaml:"password"`
	DB                 int           `yaml:"db"`
	PoolSize           int           `yaml:"pool_size"`
	MinIdleConns       int           `yaml:"min_idle_conns"`
	MaxConnAge         time.Duration `yaml:"max_conn_age"`
	PoolTimeout        time.Duration `yaml:"pool_timeout"`
	IdleTimeout        time.Duration `yaml:"idle_timeout"`
	IdleCheckFrequency time.Duration `yaml:"idle_check_frequency"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	DialTimeout        time.Duration `yaml:"dial_timeout"`
	MaxRetries         int           `yaml:"max_retries"`
	MinRetryBackoff    time.Duration `yaml:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `yaml:"max_retry_backoff"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret              string `yaml:"secret"`
	ExpireHours         int    `yaml:"expire_hours"`
	RefreshExpireHours  int    `yaml:"refresh_expire_hours"`
	Issuer              string `yaml:"issuer"`
	Audience            string `yaml:"audience"`
	SigningMethod       string `yaml:"signing_method"`
	TokenName           string `yaml:"token_name"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"dir"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	Stdout     bool   `yaml:"stdout"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	LogDir     string `yaml:"log_dir"`
	EnableCaller bool  `yaml:"enable_caller"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORSOrigins             []string        `yaml:"cors_origins"`
	CSRFSecret              string          `yaml:"csrf_secret"`
	XSSProtection           bool            `yaml:"xss_protection"`
	FrameOptions            string          `yaml:"frame_options"`
	ContentSecurityPolicy   string          `yaml:"content_security_policy"`
	RateLimit               int             `yaml:"rate_limit"`
	BruteForceProtection    bool            `yaml:"brute_force_protection"`
	PasswordComplexity      int             `yaml:"password_complexity"`
	LoginAttemptsLimit      int             `yaml:"login_attempts_limit"`
	LoginLockoutTime        int             `yaml:"login_lockout_time"`
	BcryptCost              int             `yaml:"bcrypt_cost"`
	MaxLoginAttempts        int             `yaml:"max_login_attempts"`
	LockDurationMinutes     int             `yaml:"lock_duration_minutes"`
	SessionTimeout          time.Duration   `yaml:"session_timeout"`
	APIKeyHeader            string          `yaml:"api_key_header"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string      `yaml:"allow_origins"`
	AllowMethods     []string      `yaml:"allow_methods"`
	AllowHeaders     []string      `yaml:"allow_headers"`
	AllowCredentials bool          `yaml:"allow_credentials"`
	ExposeHeaders    []string      `yaml:"expose_headers"`
	MaxAge           time.Duration `yaml:"max_age"`
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxRequestSize      int           `yaml:"max_request_size"`
	MaxUploadSize       string        `yaml:"max_upload_size"`
	RequestTimeout      time.Duration `yaml:"request_timeout"`
	ResponseCompression bool          `yaml:"response_compression"`
	GzipLevel           int           `yaml:"gzip_level"`
	CacheControl        string        `yaml:"cache_control"`
	ETag                bool          `yaml:"etag"`
	RequestRateLimit    int           `yaml:"request_rate_limit"`
	BurstRateLimit      int           `yaml:"burst_rate_limit"`
	CompressionLevel    int           `yaml:"compression_level"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Dir                string   `yaml:"dir"`
	MaxSize            int      `yaml:"max_size"`
	AllowedTypes       []string `yaml:"allowed_types"`
	FilePermissions    int      `yaml:"file_permissions"`
	DirPermissions     int      `yaml:"dir_permissions"`
	MaxFileSize        string   `yaml:"max_file_size"`
	StorageType        string   `yaml:"storage_type"`
	StoragePath        string   `yaml:"storage_path"`
	MaxFilesPerRequest int      `yaml:"max_files_per_request"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	HealthCheckPath string `yaml:"health_check_path"`
	MetricsPath     string `yaml:"metrics_path"`
}

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Title    string `yaml:"title"`
	Version  string `yaml:"version"`
	Host     string `yaml:"host"`
	BasePath string `yaml:"base_path"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		System: SystemConfig{
			IP:       "127.0.0.1",
			Port:     8080,
			Name:     "RBAC管理员系统",
			Version:  "1.0.0",
			Timezone: "Asia/Shanghai",
		},
		DB: DBConfig{
			Mode:            "sqlite",
			Host:            "localhost",
			Port:            3306,
			User:            "root",
			Password:        "",
			DbNAME:          "rbac_admin.db",
			SSLMode:         "disable",
			Timeout:         "30s",
			Charset:         "utf8mb4",
			Collation:       "utf8mb4_general_ci",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
			Path:            ":memory:",
		},
		Redis: RedisConfig{
			Addr:               "localhost:6379",
			Password:           "",
			DB:                 0,
			PoolSize:           10,
			MinIdleConns:       5,
			MaxConnAge:         time.Hour,
			PoolTimeout:        30 * time.Second,
			IdleTimeout:        5 * time.Minute,
			IdleCheckFrequency: time.Minute,
			ReadTimeout:        3 * time.Second,
			WriteTimeout:       3 * time.Second,
			DialTimeout:        5 * time.Second,
			MaxRetries:         3,
			MinRetryBackoff:    time.Millisecond,
			MaxRetryBackoff:    500 * time.Millisecond,
		},
		JWT: JWTConfig{
			Secret:              "your_jwt_secret_key_minimum_32_characters",
			ExpireHours:         24,
			RefreshExpireHours:  168,
			Issuer:              "rbac-admin-server",
			Audience:            "rbac-client",
			SigningMethod:       "HS256",
			TokenName:           "Authorization",
		},
		Log: LogConfig{
			Level:      "info",
			Dir:        "./logs",
			Filename:   "app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			Stdout:     true,
			Format:     "text",
			Output:     "both",
			LogDir:     "./logs",
			EnableCaller: true,
		},
		Security: SecurityConfig{
			CORSOrigins:           []string{"*"},
			CSRFSecret:            "your_csrf_secret_key_here",
			XSSProtection:         true,
			FrameOptions:          "DENY",
			ContentSecurityPolicy: "default-src 'self'",
			RateLimit:             100,
			BruteForceProtection:  true,
			PasswordComplexity:    8,
			LoginAttemptsLimit:    5,
			LoginLockoutTime:      30,
			BcryptCost:            10,
			MaxLoginAttempts:      5,
			LockDurationMinutes:   30,
			SessionTimeout:        2 * time.Hour,
			APIKeyHeader:          "X-API-Key",
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
			ExposeHeaders:    []string{},
			MaxAge:           12 * time.Hour,
		},
		Performance: PerformanceConfig{
			MaxRequestSize:      10,
			MaxUploadSize:       "10MB",
			RequestTimeout:      30 * time.Second,
			ResponseCompression: true,
			GzipLevel:           6,
			CacheControl:        "no-cache",
			ETag:                true,
			RequestRateLimit:    100,
			BurstRateLimit:      200,
			CompressionLevel:    6,
		},
		Upload: UploadConfig{
			Dir:                "./uploads",
			MaxSize:            50,
			AllowedTypes:       []string{"image/jpeg", "image/png", "application/pdf", "application/zip"},
			FilePermissions:    0644,
			DirPermissions:     0755,
			MaxFileSize:        "10MB",
			StorageType:        "local",
			StoragePath:        "./uploads",
			MaxFilesPerRequest: 5,
		},
		Monitoring: MonitoringConfig{
			HealthCheckPath: "/health",
			MetricsPath:     "/metrics",
		},
		Swagger: SwaggerConfig{
			Title:    "RBAC管理员API",
			Version:  "1.0.0",
			Host:     "localhost:8080",
			BasePath: "/api/v1",
		},
		App: AppConfig{
			Name:        "RBAC管理员",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       false,
		},
	}
}
