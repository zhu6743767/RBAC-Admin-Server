package config

// AppConfig 应用配置
type AppConfig struct {
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	Description   string `yaml:"description"`
	Copyright     string `yaml:"copyright"`
	Timezone      string `yaml:"timezone"`
	Language      string `yaml:"language"`
	Debug         bool   `yaml:"debug"`
}