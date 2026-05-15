package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type TaskStatusCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTaskStatusCache(client *redis.Client, ttlSeconds int64) *TaskStatusCache {
	if ttlSeconds <= 0 {
		ttlSeconds = 1800
	}
	return &TaskStatusCache{
		client: client,
		ttl:    time.Duration(ttlSeconds) * time.Second,
	}
}

// 根据 key 获取 value
func (c *TaskStatusCache) Get(ctx context.Context, taskID string) (string, bool, error) {
	if c == nil || c.client == nil {
		return "", false, nil
	}

	status, err := c.client.Get(ctx, TaskStatusKey(taskID)).Result()
	if err != nil {
		// key 不存在
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		// 其他错误
		return "", false, err
	}

	return status, true, nil
}

// 设置 key-value 写入 Redis，并设置过期时间
func (c *TaskStatusCache) Set(ctx context.Context, taskID string, status string) error {
	if c == nil || c.client ==nil {
		return nil
	}

	// Todo: 这里的错误可以变成统一错误处理
	return c.client.Set(ctx, TaskStatusKey(taskID), status, c.ttl).Err()
}

// 生成 Redis key 的辅助函数
func TaskStatusKey(taskID string) string {
	return "task:status:" + taskID
}
