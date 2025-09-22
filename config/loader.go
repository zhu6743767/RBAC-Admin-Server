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

	// 获取默认配置
	cfg := DefaultConfig()

	// 如果指定了配置文件，读取并解析
	if filename != "" {
		// 读取配置文件
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		// 替换环境变量占位符
		content := string(data)
		content = replaceEnvVars(content)

		// 解析YAML配置
		if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 应用环境变量
	applyEnvironmentVariables(cfg)

	// 验证配置
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return cfg, nil
}

// replaceEnvVars 替换配置文件中的环境变量占位符
func replaceEnvVars(content string) string {
	// 替换 ${VAR_NAME} 格式的环境变量
	for {
		start := strings.Index(content, "$")
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

// applyDefaults 应用默认值（保留函数结构，但使用DefaultConfig）
func applyDefaults(cfg *Config) {
	// 默认值已在DefaultConfig()中设置，这里可以保留一些特定的覆盖逻辑
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
	if mode := os.Getenv("SYSTEM_MODE"); mode != "" {
		cfg.System.Mode = mode
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
	if issuer := os.Getenv("JWT_ISSUER"); issuer != "" {
		cfg.JWT.Issuer = issuer
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
		cfg.Log.Dir = logDir
	}
	if stdout := os.Getenv("LOG_STDOUT"); stdout != "" {
		if b, err := strconv.ParseBool(stdout); err == nil {
			cfg.Log.Stdout = b
		}
	}

	// 应用配置
	if name := os.Getenv("APP_NAME"); name != "" {
		cfg.App.Name = name
	}
	if version := os.Getenv("APP_VERSION"); version != "" {
		cfg.App.Version = version
	}
	if debug := os.Getenv("APP_DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			cfg.App.Debug = d
		}
	}
}