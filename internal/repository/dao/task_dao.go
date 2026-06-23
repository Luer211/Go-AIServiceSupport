package dao

import (
	"context"

	"Go-AIServiceSupport/internal/model"

	"gorm.io/gorm"
)

type TaskDao struct {
	db *gorm.DB
}

func NewTaskDao(db *gorm.DB) *TaskDao {
	return &TaskDao{db: db}
}

// 把任务写入 tasks 表。
func (d *TaskDao) Create(ctx context.Context, task *model.Task) error {
	err := d.db.WithContext(ctx).Create(task).Error
	return wrapDBError("create task", err)
}

// 按 task_id 去 MySQL 查询任务
func (d *TaskDao) FindByTaskID(ctx context.Context, taskID string) (*model.Task, error) {
	var task model.Task
	err := d.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		First(&task).Error
	if err != nil {
		return nil, wrapDBError("find task by task id", err)
	}
	return &task, nil
}

// 消费端完成任务后更新: MySQL 结果 + redis 状态
func (d *TaskDao) UpdateStatus(ctx context.Context, taskID string, status string, result string) error {
	// 这是消费端应该做的，此处只是说明
	return ErrNotImplemented
}
