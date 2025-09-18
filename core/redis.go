package core

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"rbac.admin/config"
)

var (
	// RedisClient 全局Redis客户端
	RedisClient *redis.Client
	// RedisCtx 全局Redis上下文
	RedisCtx = context.Background()
)

// InitRedis 初始化Redis连接
// 支持单节点模式
// 配置连接池参数，Addr为空时跳过初始化
func InitRedis(cfg *config.RedisConfig) error {
	// 如果Redis地址为空，则跳过初始化
	if cfg.Addr == "" {
		fmt.Println("⚠️ Redis配置为空，跳过初始化")
		return nil
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,             // 服务器地址
		Password: cfg.Password,         // 密码（使用正确的字段名）
		DB:       cfg.DB,               // 数据库索引

		// 连接池配置（使用配置文件中的值）
		PoolSize:            cfg.PoolSize,            // 连接池最大连接数
		MinIdleConns:        cfg.MinIdleConns,        // 最小空闲连接数
		MaxConnAge:          cfg.MaxConnAge,          // 连接的最大存活时间
		PoolTimeout:         cfg.PoolTimeout,         // 从连接池获取连接的超时时间
		IdleTimeout:         cfg.IdleTimeout,         // 空闲连接的超时时间
		IdleCheckFrequency:  cfg.IdleCheckFrequency,  // 空闲连接检查频率
		ReadTimeout:         cfg.ReadTimeout,         // 读取超时
		WriteTimeout:        cfg.WriteTimeout,        // 写入超时
		DialTimeout:         cfg.DialTimeout,         // 连接超时

		// 连接重试配置
		MaxRetries:          cfg.MaxRetries,          // 最大重试次数
		MinRetryBackoff:     cfg.MinRetryBackoff,     // 最小重试间隔
		MaxRetryBackoff:     cfg.MaxRetryBackoff,     // 最大重试间隔

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

	RedisClient = client

	fmt.Printf("✅ Redis连接成功: %s, DB: %d\n", cfg.Addr, cfg.DB)

	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return RedisClient
}

// RedisIsConnected 检查Redis是否连接成功
func RedisIsConnected() bool {
	if RedisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	return err == nil
}
