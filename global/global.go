package global

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"rbac_admin_server/config"
)

// Config 全局配置实例
// 所有包都可以通过global.Config访问配置
var Config *config.Config

// Logger 全局日志实例
// 通过global.Logger可以访问统一的日志系统
var Logger *logrus.Logger

// DB 全局数据库连接
// 通过global.DB可以访问数据库
var DB *gorm.DB

// Redis 全局Redis客户端
// 通过global.Redis可以访问Redis缓存
var Redis *redis.Client

// Casbin 全局Casbin权限管理器
// 通过global.Casbin可以进行权限验证
var Casbin *casbin.CachedEnforcer

// RedisCtx Redis上下文
// 用于Redis操作的上下文对象
var RedisCtx = context.Background()
