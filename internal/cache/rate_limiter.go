package cache

import (
	"context"
	"time"
)

type RateLimiter struct {
	client any
	window time.Duration
}

func NewRateLimiter(client any, window time.Duration) *RateLimiter {
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{
		client: client,
		window: window,
	}
}

func (l *RateLimiter) Allow(ctx context.Context, key string, limit int) (bool, error) {
	// TODO: 使用 Redis INCR + EXPIRE 实现固定窗口限流。
	return true, nil
}

func IPRateKey(ip string) string {
	return "rate:ip:" + ip
}

func UserRateKey(userID string) string {
	return "rate:user:" + userID
}
