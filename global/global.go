package global

import (
	"rbac.admin/config"
)

// Config 全局配置实例
// 所有包都可以通过global.Config访问配置
// 由main.go初始化并设置
var Config *config.Config
