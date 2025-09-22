package config

// Config 全局配置
// 所有配置项将从配置文件加载
// 如果配置文件中未指定，则使用默认值
// 通过环境变量可以覆盖配置文件中的值


// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		System: SystemConfig{
			Mode:        "dev",
			IP:          "0.0.0.0",
			Port:        8080,
			ReadTimeout: 30,
			WriteTimeout: 30,
			MaxHeaderBytes: 1 << 20,
		},
		DB: DBConfig{
			Mode:       "sqlite",
			Host:       "localhost",
			Port:       3306,
			User:       "root",
			Password:   "123456",
			DbNAME:     "rbac_admin.db",
			Path:       "./data",
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			ConnMaxLifetime: 60,
		},
		Redis: RedisConfig{
			Addr:        "localhost:6379",
			Password:    "",
			DB:          0,
			PoolSize:    100,
			MinIdleConns: 10,
			DialTimeout: 5,
			ReadTimeout: 3,
			WriteTimeout: 3,
			IdleTimeout: 60,
		},
		JWT: JWTConfig{
			Secret:     "your-secret-key",
			ExpireHours: 24,
			Issuer:     "rbac-admin",
			RenewWindow: 6,
		},
		Log: LogConfig{
			Level:      "info",
			Dir:        "./logs",
			Filename:   "app.log",
			Format:     "text",
			MaxSize:    50,
			MaxAge:     7,
			MaxBackups: 10,
			Compress:   false,
			Stdout:     true,
			EnableCaller: true,
		},
		Security: SecurityConfig{
			XSSProtection:    "1",
			ContentTypeNosniff: "1",
			XFrameOptions:    "DENY",
			CSRFProtection:   true,
			RateLimit:        100,
			BcryptCost:       10,
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
		Monitoring: MonitoringConfig{
			Enabled:          false,
			PrometheusPort:   9090,
			HealthCheckPath:  "/health",
			MetricsPath:      "/metrics",
			TraceSamplingRate: 0.1,
		},
		Swagger: SwaggerConfig{
			Enabled:          true,
			Path:             "/swagger",
			Title:            "RBAC Admin API",
			Description:      "RBAC权限管理系统API文档",
			Version:          "1.0.0",
			TermsOfService:   "",
			ContactName:      "Admin",
			ContactURL:       "",
			ContactEmail:     "admin@example.com",
			LicenseName:      "MIT",
			LicenseURL:       "https://opensource.org/licenses/MIT",
		},
		Performance: PerformanceConfig{
			MaxUploadSize:    "10MB",
			RequestRateLimit: 100,
			WorkerPoolSize:   10,
			CacheTTL:         3600,
		},
		App: AppConfig{
			Name:        "RBAC Admin",
			Version:     "1.0.0",
			Description: "RBAC权限管理系统",
			Copyright:   "© 2023 All Rights Reserved",
			Timezone:    "Asia/Shanghai",
			Language:    "zh-CN",
			Debug:       true,
		},
		Upload: UploadConfig{
			MaxFileSize:  10,
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
			SavePath:     "./uploads",
			UseHashName:  true,
			BackupEnabled: false,
			BackupPath:   "./backups",
		},
	}
}
