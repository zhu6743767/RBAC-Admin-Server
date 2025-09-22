package core

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"rbac.admin/global"
)

// RedisCtx 全局Redis上下文
var RedisCtx = context.Background()

// InitRedis 初始化Redis连接
// 支持单节点模式
// 配置连接池参数，Addr为空时跳过初始化
func InitRedis() error {
	// 如果全局配置不存在或者Redis地址为空，则跳过初始化
	if global.Config == nil || global.Config.Redis.Addr == "" {
		global.Logger.Warn("Redis配置为空，跳过初始化")
		return nil
	}

	cfg := global.Config.Redis

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,             // 服务器地址
		Password: cfg.Password,         // 密码
		DB:       cfg.DB,               // 数据库索引

		// 连接池配置（使用配置文件中的值）
		PoolSize:            cfg.PoolSize,                      // 连接池最大连接数
		MinIdleConns:        cfg.MinIdleConns,                  // 最小空闲连接数
		MaxConnAge:          time.Duration(cfg.MaxConnAge) * time.Second,          // 连接的最大存活时间
		PoolTimeout:         time.Duration(cfg.PoolTimeout) * time.Second,         // 从连接池获取连接的超时时间
		IdleTimeout:         time.Duration(cfg.IdleTimeout) * time.Second,         // 空闲连接的超时时间
		IdleCheckFrequency:  time.Duration(cfg.IdleCheckFrequency) * time.Second,  // 空闲连接检查频率
		ReadTimeout:         time.Duration(cfg.ReadTimeout) * time.Second,         // 读取超时
		WriteTimeout:        time.Duration(cfg.WriteTimeout) * time.Second,        // 写入超时
		DialTimeout:         time.Duration(cfg.DialTimeout) * time.Second,         // 连接超时

		// 连接重试配置
		MaxRetries:          cfg.MaxRetries,          // 最大重试次数
		MinRetryBackoff:     time.Duration(cfg.MinRetryBackoff) * time.Millisecond,     // 最小重试间隔
		MaxRetryBackoff:     time.Duration(cfg.MaxRetryBackoff) * time.Millisecond,     // 最大重试间隔

		// TLS配置
		TLSConfig:           nil,                     // TLS配置
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	// 设置全局Redis客户端
	global.Redis = client

	global.Logger.Info(fmt.Sprintf("Redis连接成功: %s, DB: %d", cfg.Addr, cfg.DB))

	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if global.Redis != nil {
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
