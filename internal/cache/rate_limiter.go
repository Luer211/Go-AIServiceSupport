package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	window time.Duration
}

func NewRateLimiter(client *redis.Client, window time.Duration) *RateLimiter {
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{
		client: client,
		window: window,
	}
}

// 使用 Redis INCR + EXPIRE 实现固定窗口限流
// Todo: 未实现统一的错误处理
func (l *RateLimiter) Allow(ctx context.Context, key string, limit int) (bool, error) {
	if limit <= 0 {
		return true, nil
	}

	if l == nil || l.client == nil {
		return false, nil
	}

	// 原子自增
	count, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// 如果是第一次请求，给 key 设置过期时间
	if count == 1 {
		if err := l.client.Expire(ctx, key, l.window).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(limit), nil
}

// 生成 IP 限流 key
func IPRateKey(ip string) string {
	return "rate:ip:" + ip
}

// 生成用户限流 key
func UserRateKey(userID string) string {
	return "rate:user:" + userID
}
