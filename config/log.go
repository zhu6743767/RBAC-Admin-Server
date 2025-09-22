package config

// LogConfig 日志配置
type LogConfig struct {
	Level       string `yaml:"level"`
	Dir         string `yaml:"dir"`
	Filename    string `yaml:"filename"`
	Format      string `yaml:"format"`
	MaxSize     int    `yaml:"max_size"`
	MaxAge      int    `yaml:"max_age"`
	MaxBackups  int    `yaml:"max_backups"`
	Compress    bool   `yaml:"compress"`
	Stdout      bool   `yaml:"stdout"`
	EnableCaller bool  `yaml:"enable_caller"`
}