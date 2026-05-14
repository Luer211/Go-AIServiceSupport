package cache

import (
	"context"
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
	// TODO: 查询 Redis key: task:status:{task_id}。
	return "", false, nil
}

// 设置 key-value
func (c *TaskStatusCache) Set(ctx context.Context, taskID string, status string) error {
	// TODO: 写入 Redis key: task:status:{task_id}，并设置 c.ttl。
	return nil
}

// 生成 Redis key 的辅助函数
func TaskStatusKey(taskID string) string {
	return "task:status:" + taskID
}
