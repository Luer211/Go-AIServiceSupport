package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/internal/api/request"
	"Go-AIServiceSupport/internal/api/response"
	"Go-AIServiceSupport/internal/cache"
	"Go-AIServiceSupport/internal/model"
	"Go-AIServiceSupport/internal/mq"
	"Go-AIServiceSupport/internal/repository/dao"
)

// 全局错误变量：若是无权访问任务错误，则返回ErrTaskForbidden
var ErrTaskForbidden = errors.New("task does not belong to current user")

type TaskService struct {
	tasks       *dao.TaskDao
	statusCache *cache.TaskStatusCache
	producer    mq.Producer
}

func NewTaskService(tasks *dao.TaskDao, statusCache *cache.TaskStatusCache, producer mq.Producer) (*TaskService, error) {
	if tasks == nil {
		return nil, fmt.Errorf("task dao is required")
	}
	if statusCache == nil {
		return nil, fmt.Errorf("task status cache is required")
	}
	if producer == nil {
		return nil, fmt.Errorf("task producer is required")
	}

	return &TaskService{
		tasks:       tasks,
		statusCache: statusCache,
		producer:    producer,
	}, nil
}

// 创建任务
// 流程：生成任务ID → 保存到数据库 → 发送到消息队列 → 存入缓存中pending → 返回任务信息
func (s *TaskService) CreateTask(ctx context.Context, userID uint64, req request.CreateTaskRequest) (*response.CreateTaskResponse, error) {
	// 创建任务
	taskID := newTaskID()
	task := &model.Task{
		TaskID: taskID,
		UserID: userID,
		Prompt: req.Prompt,
		Status: model.TaskStatusPending,
	}

	// 存入数据库
	if err := s.tasks.Create(ctx, task); err != nil {
		if errors.Is(err, dao.ErrAlreadyExists) {
			return nil, common.WrapAppError(e.CodeTaskSubmitFailed, err)
		}
		return nil, common.WrapAppError(e.CodeInternalError, err)
	}

	// 发送到消息队列，异步处理
	if err := s.producer.PublishTask(ctx, mq.TaskMessage{
		TaskID: taskID,
		UserID: userID,
		Prompt: req.Prompt,
	}); err != nil {
		return nil, common.WrapAppError(e.CodeTaskSubmitFailed, err)
	}

	// 将任务加入缓存中，设置状态为 pending
	if err := s.statusCache.Set(
		ctx,
		taskID,
		model.TaskStatusPending,
	); err != nil {
		return nil, common.WrapAppError(e.CodeTaskWriteInRedisFailed, err)
	}

	return &response.CreateTaskResponse{
		TaskID: taskID,
		Status: model.TaskStatusPending,
	}, nil
}

// 查询任务状态
// 流程：查询任务 → 校验权限 → 从缓存获取状态 → 若缓存没有则从数据库中取状态 → 返回结果
func (s *TaskService) GetTaskStatus(ctx context.Context, userID uint64, taskID string) (*response.GetTaskStatusResponse, error) {
	// 根据任务 ID 查询任务信息
	task, err := s.tasks.FindByTaskID(ctx, taskID)
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			return nil, common.WrapAppError(e.CodeTaskNotFound, err)
		}
		return nil, common.WrapAppError(e.CodeInternalError, err)
	}

	// 权限校验：只能查询自己的任务
	if task.UserID != userID {
		return nil, common.WrapAppError(e.CodeForbidden, ErrTaskForbidden)
	}

	// 从缓存中获取状态
	status, ok, err := s.statusCache.Get(ctx, taskID)
	if err != nil {
		return nil, common.WrapAppError(e.CodeInternalError, err)
	}
	// 如果没有，先用 MySQL 里面查询出来的
	if !ok {
		status = task.Status
	}

	return &response.GetTaskStatusResponse{
		TaskID: taskID,
		Status: status,
	}, nil
}

// 创建任务ID的工具函数：随机生成ID
func newTaskID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("task_%d", time.Now().UnixNano())
	}
	return "task_" + hex.EncodeToString(buf)
}
