package config

// JWTConfig JWT配置
type JWTConfig struct {
	Secret            string `yaml:"secret"`
	ExpireHours       int    `yaml:"expire_hours"`
	RefreshExpireHours int   `yaml:"refresh_expire_hours"`
	Issuer            string `yaml:"issuer"`
	Audience          string `yaml:"audience"`
	SigningMethod     string `yaml:"signing_method"`
	TokenName         string `yaml:"token_name"`
}