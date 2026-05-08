package model

import "time"

// 用户表
type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
