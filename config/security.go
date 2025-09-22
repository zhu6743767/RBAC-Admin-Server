package config

// SecurityConfig 安全配置
type SecurityConfig struct {
	XSSProtection      string `yaml:"xss_protection"`
	ContentTypeNosniff string `yaml:"content_type_nosniff"`
	XFrameOptions      string `yaml:"x_frame_options"`
	CSRFProtection     bool   `yaml:"csrf_protection"`
	RateLimit          int    `yaml:"rate_limit"`
	BcryptCost         int    `yaml:"bcrypt_cost"`
}