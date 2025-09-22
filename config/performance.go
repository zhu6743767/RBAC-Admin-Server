package config

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxUploadSize     string `yaml:"max_upload_size"`
	RequestRateLimit  int    `yaml:"request_rate_limit"`
	WorkerPoolSize    int    `yaml:"worker_pool_size"`
	CacheTTL          int    `yaml:"cache_ttl"`
}