package config

// DBConfig 数据库配置
type DBConfig struct {
	Mode       string `yaml:"mode"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DbNAME     string `yaml:"dbname"`
	Path       string `yaml:"path"`
	MaxIdleConns int  `yaml:"max_idle_conns"`
	MaxOpenConns int  `yaml:"max_open_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
}