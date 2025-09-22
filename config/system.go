package config

import "fmt"

// SystemConfig 系统配置
type SystemConfig struct {
	Mode        string `yaml:"mode"`
	IP          string `yaml:"ip"`
	Port        int    `yaml:"port"`
	ReadTimeout int    `yaml:"read_timeout"`
	WriteTimeout int   `yaml:"write_timeout"`
	MaxHeaderBytes int `yaml:"max_header_bytes"`
}

// Addr 返回服务器地址
func (s *SystemConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}