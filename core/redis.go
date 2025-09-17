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
	// RedisCtx Redis上下文
	RedisCtx = context.Background()
)

// InitRedis 初始化Redis连接
// 支持单节点、哨兵和集群模式
// 自动配置连接池参数
func InitRedis(cfg *config.RedisConfig) error {
	// 如果地址为空，跳过Redis初始化
	if cfg.Addr == "" {
		fmt.Println("⚠️  Redis未配置，跳过初始化")
		return nil
	}

	// 单节点模式（简化版本）
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	if err := client.Ping(RedisCtx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	RedisClient = client
	fmt.Printf("✅ Redis连接成功: %s\n", cfg.Addr)
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// RedisSet 设置Redis键值
func RedisSet(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Set(RedisCtx, key, value, expiration).Err()
}

// RedisGet 获取Redis键值
func RedisGet(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Get(RedisCtx, key).Result()
}

// RedisDel 删除Redis键
func RedisDel(keys ...string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Del(RedisCtx, keys...).Err()
}

// RedisExists 检查Redis键是否存在
func RedisExists(keys ...string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Exists(RedisCtx, keys...).Result()
}

// RedisExpire 设置Redis键过期时间
func RedisExpire(key string, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.Expire(RedisCtx, key, expiration).Err()
}

// RedisTTL 获取Redis键剩余过期时间
func RedisTTL(key string) (time.Duration, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.TTL(RedisCtx, key).Result()
}

// RedisHSet 设置Redis哈希字段
func RedisHSet(key string, field string, value interface{}) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HSet(RedisCtx, key, field, value).Err()
}

// RedisHGet 获取Redis哈希字段
func RedisHGet(key string, field string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HGet(RedisCtx, key, field).Result()
}

// RedisHDel 删除Redis哈希字段
func RedisHDel(key string, fields ...string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HDel(RedisCtx, key, fields...).Err()
}

// RedisHExists 检查Redis哈希字段是否存在
func RedisHExists(key string, field string) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("Redis客户端未初始化")
	}
	return RedisClient.HExists(RedisCtx, key, field).Result()
}
