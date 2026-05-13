package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

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

func NewTaskService(tasks *dao.TaskDao, statusCache *cache.TaskStatusCache, producer mq.Producer) *TaskService {
	if producer == nil {
		producer = mq.NewNoopProducer()
	}
	return &TaskService{
		tasks:       tasks,
		statusCache: statusCache,
		producer:    producer,
	}
}

// 创建任务
// 流程：生成任务ID → 保存到数据库 → 发送到消息队列 → 存入缓存中running → 返回任务信息
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
		return nil, err
	}

	// 发送到消息队列，异步处理
	if err := s.producer.PublishTask(ctx, mq.TaskMessage{
		TaskID: taskID,
		UserID: userID,
		Prompt: req.Prompt,
	}); err != nil {
		return nil, err
	}

	// Todo: 将任务加入缓存中，设置状态为 running

	return &response.CreateTaskResponse{
		TaskID: taskID,
		Status: model.TaskStatusPending,
	}, nil
}

// 查询任务状态
// 流程：查询任务 → 校验权限 → 从缓存获取状态 → 返回结果
func (s *TaskService) GetTaskStatus(ctx context.Context, userID uint64, taskID string) (*response.GetTaskStatusResponse, error) {
	// 根据任务 ID 查询任务信息
	task, err := s.tasks.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// 权限校验：只能查询自己的任务
	if task.UserID != userID {
		return nil, ErrTaskForbidden
	}

	// 从缓存中获取状态
	status, ok, err := s.statusCache.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if !ok {
		status = model.TaskStatusPending
	}

	return &response.GetTaskStatusResponse{
		TaskID: taskID,
		Status: status,
	}, nil
}

// 创建任务ID的工具函数
func newTaskID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("task_%d", time.Now().UnixNano())
	}
	return "task_" + hex.EncodeToString(buf)
}
