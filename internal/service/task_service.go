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

func (s *TaskService) CreateTask(ctx context.Context, userID uint64, req request.CreateTaskRequest) (*response.CreateTaskResponse, error) {
	taskID := newTaskID()
	task := &model.Task{
		TaskID: taskID,
		UserID: userID,
		Prompt: req.Prompt,
		Status: model.TaskStatusPending,
	}

	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}

	if err := s.producer.PublishTask(ctx, mq.TaskMessage{
		TaskID: taskID,
		UserID: userID,
		Prompt: req.Prompt,
	}); err != nil {
		return nil, err
	}

	return &response.CreateTaskResponse{
		TaskID: taskID,
		Status: model.TaskStatusPending,
	}, nil
}

func (s *TaskService) GetTaskStatus(ctx context.Context, userID uint64, taskID string) (*response.GetTaskStatusResponse, error) {
	task, err := s.tasks.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrTaskForbidden
	}

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

func newTaskID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("task_%d", time.Now().UnixNano())
	}
	return "task_" + hex.EncodeToString(buf)
}
