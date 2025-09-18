package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// DB 全局数据库连接
// 通过global.DB可以访问数据库
var DB *gorm.DB
