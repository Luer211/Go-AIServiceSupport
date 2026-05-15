package initialize

import (
	"context"
	"time"

	"Go-AIServiceSupport/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {
	// 创建 Redis 连接
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Todo: 这里使用 panic，暂时还没纳入统一错误处理
	// Ping 一下
	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return client
}
