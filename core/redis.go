package core

import (
	"context"
	"fmt"
	"time"

	"rbac_admin_server/core/init_redis"
	"rbac_admin_server/global"
)

// RedisCtx 全局Redis上下文
var RedisCtx = context.Background()

// InitRedis 初始化Redis连接（已废弃，请使用init_redis包中的InitRedis函数）
// 注意：此函数从版本1.2.0开始已废弃
// 建议使用core/init_redis包中的InitRedis函数进行Redis初始化
func InitRedis() error {
	global.Logger.Warn("⚠️ core/InitRedis已废弃，请使用init_redis包中的InitRedis函数")
	if global.Config == nil || global.Config.Redis.Addr == "" {
		global.Logger.Warn("Redis配置为空，跳过初始化")
		return nil
	}

	// 调用新的初始化函数
	client, err := init_redis.InitRedis()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	// 设置全局Redis客户端
	global.Redis = client

	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if global.Redis != nil {
		global.Logger.Info("🔄 正在关闭Redis连接...")
		return global.Redis.Close()
	}
	return nil
}

// RedisIsConnected 检查Redis是否连接成功
func RedisIsConnected() bool {
	if global.Redis == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := global.Redis.Ping(ctx).Result()
	return err == nil
}
