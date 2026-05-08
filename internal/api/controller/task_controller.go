package controller

import (
	"errors"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/internal/api/request"
	"Go-AIServiceSupport/internal/service"
	"Go-AIServiceSupport/middle"
)

type TaskController struct {
	taskService *service.TaskService
}

func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

func (ctl *TaskController) CreateTask(ctx *gin.Context) {
	var req request.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	resp, err := ctl.taskService.CreateTask(ctx.Request.Context(), middle.CurrentUserID(ctx), req)
	if err != nil {
		common.FailWithMessage(ctx, e.CodeTaskSubmitFailed, err.Error())
		return
	}

	common.Success(ctx, resp)
}

func (ctl *TaskController) GetTaskStatus(ctx *gin.Context) {
	var req request.GetTaskStatusRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	resp, err := ctl.taskService.GetTaskStatus(ctx.Request.Context(), middle.CurrentUserID(ctx), req.TaskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskForbidden) {
			common.Fail(ctx, e.CodeForbidden)
			return
		}
		common.FailWithMessage(ctx, e.CodeTaskNotFound, err.Error())
		return
	}

	common.Success(ctx, resp)
}
