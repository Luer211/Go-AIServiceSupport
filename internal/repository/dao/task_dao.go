package dao

import (
	"context"

	"Go-AIServiceSupport/internal/model"
)

type TaskDao struct {
	db any
}

func NewTaskDao(db any) *TaskDao {
	return &TaskDao{db: db}
}

func (d *TaskDao) Create(ctx context.Context, task *model.Task) error {
	// TODO: 使用 GORM 写入 tasks 表。
	return ErrNotImplemented
}

func (d *TaskDao) FindByTaskID(ctx context.Context, taskID string) (*model.Task, error) {
	// TODO: 使用 GORM 按 task_id 查询任务，并用于校验任务归属。
	return nil, ErrNotImplemented
}

func (d *TaskDao) UpdateStatus(ctx context.Context, taskID string, status string, result string) error {
	// TODO: 消费端完成任务后更新 MySQL 状态和结果。
	return ErrNotImplemented
}
