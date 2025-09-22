package config

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
	MinIdleConns int `yaml:"min_idle_conns"`
	DialTimeout int `yaml:"dial_timeout"`
	ReadTimeout int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
	IdleTimeout int `yaml:"idle_timeout"`
	MaxConnAge int `yaml:"max_conn_age"`
	PoolTimeout int `yaml:"pool_timeout"`
	IdleCheckFrequency int `yaml:"idle_check_frequency"`
	MaxRetries int `yaml:"max_retries"`
	MinRetryBackoff int `yaml:"min_retry_backoff"`
	MaxRetryBackoff int `yaml:"max_retry_backoff"`
}