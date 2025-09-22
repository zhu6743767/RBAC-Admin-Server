package config

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
	Issuer      string `yaml:"issuer"`
	RenewWindow int    `yaml:"renew_window"` // 续期窗口(小时)
}