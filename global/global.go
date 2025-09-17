package global

import (
	"github.com/sirupsen/logrus"
	"rbac.admin/config"
)

// Config 全局配置实例
// 所有包都可以通过global.Config访问配置
// 由main.go初始化并设置
var Config *config.Config

// Logger 全局日志实例
// 通过global.Logger可以访问统一的日志系统
// 支持按时间、大小、等级分片
var Logger *logrus.Logger

// ConfigManager 全局配置管理器实例
// 提供配置热重载和回写功能
// 由main.go初始化并提供
var ConfigManager interface {
	GetConfig() *config.Config
	SaveConfig(*config.Config) error
	UpdateConfig(func(*config.Config)) error
	Close() error
}

// DBManager 全局数据库管理器
var DBManager interface{}
