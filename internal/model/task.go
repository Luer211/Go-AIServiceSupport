package model

import "time"

const (
	TaskStatusPending = "pending"
	TaskStatusSuccess = "success"
	TaskStatusFailed  = "failed"
)

// 任务表
type Task struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	TaskID    string `gorm:"size:64;uniqueIndex;not null"`
	UserID    uint64 `gorm:"index;not null"`
	Prompt    string `gorm:"type:text;not null"`
	Result    string `gorm:"type:longtext"`
	Status    string `gorm:"size:16;index;not null;default:pending"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
