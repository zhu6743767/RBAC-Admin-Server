package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load 从配置文件加载配置，支持环境变量替换
// 如果配置文件不存在，则创建默认配置
func Load(configPath string) (*Config, error) {
	// 确保配置文件路径是绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件绝对路径失败: %w", err)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		// 创建默认配置文件
		fmt.Printf("配置文件 %s 不存在，创建默认配置...\n", absPath)
		defaultConfig := DefaultConfig()
		if err := Save(absPath, defaultConfig); err != nil {
			return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
		}
		return defaultConfig, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 替换环境变量
	configStr := replaceEnvVars(string(data))

	// 解析配置
	var cfg Config
	if err := yaml.Unmarshal([]byte(configStr), &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := Validate(&cfg); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &cfg, nil
}

// Save 保存配置到文件
func Save(configPath string, config *Config) error {
	// 确保目录存在
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("获取配置文件绝对路径失败: %w", err)
	}

	// 创建目录（如果不存在）
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 序列化配置
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Validate 验证配置的有效性
func Validate(config *Config) error {
	return validateAll(config)
}

// replaceEnvVars 替换字符串中的环境变量
func replaceEnvVars(str string) string {
	return os.Expand(str, func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		// 支持默认值，格式: ${VAR:-default}
		if strings.Contains(key, ":-") {
			parts := strings.SplitN(key, ":-", 2)
			if len(parts) == 2 {
				if value, exists := os.LookupEnv(parts[0]); exists {
					return value
				}
				return parts[1]
			}
		}
		return ""
	})
}

// validateAll 执行所有验证
func validateAll(config *Config) error {
	if err := validateServer(config.Server); err != nil {
		return fmt.Errorf("服务器配置错误: %w", err)
	}
	if err := validateDatabase(config.Database); err != nil {
		return fmt.Errorf("数据库配置错误: %w", err)
	}
	if err := validateJWT(config.JWT); err != nil {
		return fmt.Errorf("JWT配置错误: %w", err)
	}
	if err := validateRedis(config.Redis); err != nil {
		return fmt.Errorf("Redis配置错误: %w", err)
	}
	if err := validateLog(config.Log); err != nil {
		return fmt.Errorf("日志配置错误: %w", err)
	}
	if err := validateSecurity(config.Security); err != nil {
		return fmt.Errorf("安全配置错误: %w", err)
	}
	if err := validateCORS(config.CORS); err != nil {
		return fmt.Errorf("CORS配置错误: %w", err)
	}
	if err := validatePerformance(config.Performance); err != nil {
		return fmt.Errorf("性能配置错误: %w", err)
	}
	if err := validateUpload(config.Upload); err != nil {
		return fmt.Errorf("上传配置错误: %w", err)
	}
	if err := validateMonitoring(config.Monitoring); err != nil {
		return fmt.Errorf("监控配置错误: %w", err)
	}
	if err := validateSwagger(config.Swagger); err != nil {
		return fmt.Errorf("Swagger配置错误: %w", err)
	}

	return nil
}

func validateServer(server ServerConfig) error {
	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("端口必须在1-65535之间")
	}
	if server.ReadTimeout <= 0 {
		return fmt.Errorf("读取超时必须大于0")
	}
	if server.WriteTimeout <= 0 {
		return fmt.Errorf("写入超时必须大于0")
	}
	if server.ShutdownTimeout <= 0 {
		return fmt.Errorf("关闭超时必须大于0")
	}
	return nil
}

func validateDatabase(db DatabaseConfig) error {
	if db.Host == "" {
		return fmt.Errorf("主机不能为空")
	}
	if db.Port <= 0 || db.Port > 65535 {
		return fmt.Errorf("端口必须在1-65535之间")
	}
	if db.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if db.MaxIdleConns < 0 {
		return fmt.Errorf("最大空闲连接数不能为负数")
	}
	if db.MaxOpenConns < 0 {
		return fmt.Errorf("最大打开连接数不能为负数")
	}
	if db.ConnMaxLifetime <= 0 {
		return fmt.Errorf("连接最大生命周期必须大于0")
	}
	return nil
}

func validateJWT(jwt JWTConfig) error {
	if jwt.Secret == "" {
		return fmt.Errorf("密钥不能为空")
	}
	if jwt.ExpireHours <= 0 {
		return fmt.Errorf("过期时间必须大于0")
	}
	if jwt.RefreshExpireHours <= 0 {
		return fmt.Errorf("刷新过期时间必须大于0")
	}
	return nil
}

func validateRedis(redis RedisConfig) error {
	if redis.Host == "" {
		return fmt.Errorf("主机不能为空")
	}
	if redis.Port <= 0 || redis.Port > 65535 {
		return fmt.Errorf("端口必须在1-65535之间")
	}
	if redis.DB < 0 {
		return fmt.Errorf("数据库索引不能为负数")
	}
	if redis.MaxRetries < 0 {
		return fmt.Errorf("最大重试次数不能为负数")
	}
	if redis.PoolSize <= 0 {
		return fmt.Errorf("连接池大小必须大于0")
	}
	return nil
}

func validateLog(log LogConfig) error {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[log.Level] {
		return fmt.Errorf("日志级别必须是debug、info、warn或error")
	}
	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[log.Format] {
		return fmt.Errorf("日志格式必须是json或text")
	}
	if log.MaxSize <= 0 {
		return fmt.Errorf("最大文件大小必须大于0")
	}
	if log.MaxBackups <= 0 {
		return fmt.Errorf("最大备份文件数必须大于0")
	}
	if log.MaxAge <= 0 {
		return fmt.Errorf("最大保存天数必须大于0")
	}
	return nil
}

func validateSecurity(security SecurityConfig) error {
	if security.BcryptCost < 4 || security.BcryptCost > 31 {
		return fmt.Errorf("Bcrypt成本必须在4-31之间")
	}
	if security.MaxLoginAttempts <= 0 {
		return fmt.Errorf("最大登录尝试次数必须大于0")
	}
	if security.LockDurationMinutes <= 0 {
		return fmt.Errorf("账户锁定时间必须大于0")
	}
	return nil
}

func validateCORS(cors CORSConfig) error {
	if cors.Enable && len(cors.AllowOrigins) == 0 {
		return fmt.Errorf("启用CORS时必须指定允许的来源")
	}
	return nil
}

func validatePerformance(performance PerformanceConfig) error {
	// 移除旧验证逻辑，新结构中没有这些字段
	return nil
}

func validateUpload(upload UploadConfig) error {
	if upload.MaxFileSize == "" {
		return fmt.Errorf("最大文件大小不能为空")
	}
	if len(upload.AllowedExtensions) == 0 {
		return fmt.Errorf("必须指定允许的文件扩展名")
	}
	return nil
}

func validateMonitoring(monitoring MonitoringConfig) error {
	// 移除旧验证逻辑，新结构中没有这些字段
	return nil
}

func validateSwagger(swagger SwaggerConfig) error {
	if swagger.Enable && swagger.Host == "" {
		return fmt.Errorf("启用Swagger时必须指定主机地址")
	}
	return nil
}