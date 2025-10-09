package init_redis

import (
	"context"
	"fmt"
	"rbac_admin_server/global"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() (*redis.Client, error) {
	config := global.Config.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		// 连接池设置
		PoolSize:            config.PoolSize,
		MinIdleConns:        config.MinIdleConns,
		DialTimeout:         time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:         time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:        time.Duration(config.WriteTimeout) * time.Second,
		PoolTimeout:         time.Duration(config.PoolTimeout) * time.Second,
		MaxRetries:          config.MaxRetries,
		MinRetryBackoff:     time.Duration(config.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff:     time.Duration(config.MaxRetryBackoff) * time.Millisecond,
	})

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	logrus.Info("✅ Redis连接初始化成功")
	return client, nil
}

// CloseRedis 关闭Redis连接
func CloseRedis(client *redis.Client) error {
	if client != nil {
		if err := client.Close(); err != nil {
			return fmt.Errorf("Redis关闭失败: %w", err)
		}
	}
	return nil
}

// SetRedis 设置Redis键值对
func SetRedis(key string, value interface{}, expiration time.Duration) error {
	if global.Redis == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	return global.Redis.Set(ctx, key, value, expiration).Err()
}

// GetRedis 获取Redis键值
func GetRedis(key string) (string, error) {
	if global.Redis == nil {
		return "", fmt.Errorf("Redis客户端未初始化")
	}

	return global.Redis.Get(ctx, key).Result()
}

// DelRedis 删除Redis键
func DelRedis(key string) error {
	if global.Redis == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	return global.Redis.Del(ctx, key).Err()
}