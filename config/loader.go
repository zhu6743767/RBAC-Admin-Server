// Package config 提供配置加载和环境变量解析功能
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// parseDuration 解析时间字符串为time.Duration
func parseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	return 0
}

// Load 从指定文件加载配置
func Load(filename string) (*Config, error) {
	// 首先尝试加载 .env 文件
	_ = godotenv.Load()

	// 读取配置文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 替换环境变量占位符
	content := string(data)
	content = replaceEnvVars(content)

	// 解析YAML配置
	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 应用环境变量
	applyEnvironmentVariables(&cfg)

	// 应用默认值
	applyDefaults(&cfg)

	// 验证配置
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// replaceEnvVars 替换配置文件中的环境变量占位符
func replaceEnvVars(content string) string {
	// 替换 ${VAR_NAME} 格式的环境变量
	for {
		start := strings.Index(content, "${")
		if start == -1 {
			break
		}
		end := strings.Index(content[start:], "}")
		if end == -1 {
			break
		}
		fullMatch := content[start : start+end+1]
		varName := content[start+2 : start+end]
		// 获取环境变量值，如果不存在则使用默认值
		value := os.Getenv(varName)
		if value == "" {
			// 尝试从 .env 文件加载
			value = getEnvFromFile(varName)
		}
		content = strings.ReplaceAll(content, fullMatch, value)
	}
	return content
}

// getEnvFromFile 从 .env 文件获取环境变量
func getEnvFromFile(key string) string {
	// 简单的 .env 文件解析
	if _, err := os.Stat(".env"); err == nil {
		data, _ := os.ReadFile(".env")
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, "=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}
	return ""
}

// applyDefaults 应用默认值
func applyDefaults(cfg *Config) {
	if cfg.System.Port == 0 {
		cfg.System.Port = 8080
	}
	if cfg.System.IP == "" {
		cfg.System.IP = "127.0.0.1"
	}

	if cfg.DB.Mode == "" {
		cfg.DB.Mode = "sqlite"
	}
	if cfg.DB.Port == 0 {
		cfg.DB.Port = 3306
	}
	if cfg.DB.MaxOpenConns == 0 {
		cfg.DB.MaxOpenConns = 100
	}
	if cfg.DB.MaxIdleConns == 0 {
		cfg.DB.MaxIdleConns = 10
	}
	if cfg.DB.ConnMaxLifetime == 0 {
		cfg.DB.ConnMaxLifetime = 1 * time.Hour
	}
	if cfg.DB.ConnMaxIdleTime == 0 {
		cfg.DB.ConnMaxIdleTime = 30 * time.Minute
	}

	if cfg.DB.Timeout == "" {
		cfg.DB.Timeout = "30s"
	}

	if cfg.JWT.ExpireHours == 0 {
		cfg.JWT.ExpireHours = 24
	}
	if cfg.JWT.RefreshExpireHours == 0 {
		cfg.JWT.RefreshExpireHours = 168
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "rbac-admin-server"
	}
	if cfg.JWT.Audience == "" {
		cfg.JWT.Audience = "rbac-client"
	}

	// Redis配置可选，不设置默认值

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "text"
	}
	if cfg.Log.Output == "" {
		cfg.Log.Output = "both"
	}
	if cfg.Log.LogDir == "" {
		cfg.Log.LogDir = "./logs"
	}
	if cfg.Log.MaxSize == 0 {
		cfg.Log.MaxSize = 100
	}
	if cfg.Log.MaxAge == 0 {
		cfg.Log.MaxAge = 7
	}
	if cfg.Log.MaxBackups == 0 {
		cfg.Log.MaxBackups = 3
	}

	if cfg.Performance.MaxUploadSize == "" {
		cfg.Performance.MaxUploadSize = "10MB"
	}
	if cfg.Performance.RequestRateLimit == 0 {
		cfg.Performance.RequestRateLimit = 100
	}
	if cfg.Performance.BurstRateLimit == 0 {
		cfg.Performance.BurstRateLimit = 200
	}
	if cfg.Performance.CompressionLevel == 0 {
		cfg.Performance.CompressionLevel = 6
	}

	if cfg.Upload.MaxFileSize == "" {
		cfg.Upload.MaxFileSize = "10MB"
	}
	if cfg.Upload.StorageType == "" {
		cfg.Upload.StorageType = "local"
	}
	if cfg.Upload.StoragePath == "" {
		cfg.Upload.StoragePath = "./uploads"
	}
	if cfg.Upload.MaxFilesPerRequest == 0 {
		cfg.Upload.MaxFilesPerRequest = 5
	}

	if cfg.Monitoring.HealthCheckPath == "" {
		cfg.Monitoring.HealthCheckPath = "/health"
	}
	if cfg.Monitoring.MetricsPath == "" {
		cfg.Monitoring.MetricsPath = "/metrics"
	}

	if cfg.CORS.MaxAge == 0 {
		cfg.CORS.MaxAge = 12 * time.Hour
	}

	if cfg.Security.BcryptCost == 0 {
		cfg.Security.BcryptCost = 10
	}
	if cfg.Security.MaxLoginAttempts == 0 {
		cfg.Security.MaxLoginAttempts = 5
	}
	if cfg.Security.LockDurationMinutes == 0 {
		cfg.Security.LockDurationMinutes = 30
	}
	if cfg.Security.SessionTimeout == 0 {
		cfg.Security.SessionTimeout = 2 * time.Hour
	}
	if cfg.Security.APIKeyHeader == "" {
		cfg.Security.APIKeyHeader = "X-API-Key"
	}

	if cfg.Swagger.Title == "" {
		cfg.Swagger.Title = "RBAC管理员API"
	}
	if cfg.Swagger.Version == "" {
		cfg.Swagger.Version = "1.0.0"
	}
	if cfg.Swagger.Host == "" {
		cfg.Swagger.Host = "localhost:8080"
	}
	if cfg.Swagger.BasePath == "" {
		cfg.Swagger.BasePath = "/api/v1"
	}

	if cfg.App.Name == "" {
		cfg.App.Name = "RBAC管理员"
	}
	if cfg.App.Version == "" {
		cfg.App.Version = "1.0.0"
	}
	if cfg.App.Environment == "" {
		cfg.App.Environment = "development"
	}
}

// validateConfig 验证配置有效性
func validateConfig(cfg *Config) error {
	if cfg.System.Port <= 0 || cfg.System.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	if cfg.DB.Mode == "" {
		return fmt.Errorf("数据库类型不能为空")
	}

	if cfg.DB.Mode == "sqlite" && cfg.DB.Path == "" {
		cfg.DB.Path = ":memory:"
	}

	if cfg.DB.Mode != "sqlite" && cfg.DB.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}

	return nil
}

// applyEnvironmentVariables 应用环境变量到配置
func applyEnvironmentVariables(cfg *Config) {
	// 系统配置
	if ip := os.Getenv("SYSTEM_IP"); ip != "" {
		cfg.System.IP = ip
	}
	if port := os.Getenv("SYSTEM_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.System.Port = p
		}
	}

	// 数据库配置
	if dbMode := os.Getenv("DB_MODE"); dbMode != "" {
		cfg.DB.Mode = dbMode
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.DB.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.DB.Port = p
		}
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		cfg.DB.User = username
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.DB.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.DB.Password = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DB.DbNAME = dbName
	}
	if dbname := os.Getenv("DB_DBNAME"); dbname != "" {
		cfg.DB.DbNAME = dbname
	}
	if path := os.Getenv("DB_PATH"); path != "" {
		cfg.DB.Path = path
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if expire := os.Getenv("JWT_EXPIRE_HOURS"); expire != "" {
		if h, err := strconv.Atoi(expire); err == nil {
			cfg.JWT.ExpireHours = h
		}
	}
	if refreshExpire := os.Getenv("JWT_REFRESH_EXPIRE_HOURS"); refreshExpire != "" {
		if h, err := strconv.Atoi(refreshExpire); err == nil {
			cfg.JWT.RefreshExpireHours = h
		}
	}
	if issuer := os.Getenv("JWT_ISSUER"); issuer != "" {
		cfg.JWT.Issuer = issuer
	}
	if audience := os.Getenv("JWT_AUDIENCE"); audience != "" {
		cfg.JWT.Audience = audience
	}

	// Redis配置
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}

	// 日志配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if logDir := os.Getenv("LOG_DIR"); logDir != "" {
		cfg.Log.LogDir = logDir
	}

	// 应用配置
	if name := os.Getenv("APP_NAME"); name != "" {
		cfg.App.Name = name
	}
	if version := os.Getenv("APP_VERSION"); version != "" {
		cfg.App.Version = version
	}
	if env := os.Getenv("APP_ENVIRONMENT"); env != "" {
		cfg.App.Environment = env
	}
	if debug := os.Getenv("APP_DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			cfg.App.Debug = d
		}
	}
}