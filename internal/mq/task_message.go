package mq

type TaskMessage struct {
	TaskID string `json:"task_id"`
	UserID uint64 `json:"user_id"`
	Prompt string `json:"prompt"`
}
