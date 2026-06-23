package controller

import (
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

// 任务创建接口
func (ctl *TaskController) CreateTask(ctx *gin.Context) {
	var req request.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	// 调用服务层创建任务：传入上下文、当前登录用户ID、请求参数
	resp, err := ctl.taskService.CreateTask(ctx.Request.Context(), middle.CurrentUserID(ctx), req)
	if err != nil {
		common.ErrorResponse(ctx, err)
		return
	}

	common.Success(ctx, resp)
}

// 查询任务状态接口
func (ctl *TaskController) GetTaskStatus(ctx *gin.Context) {
	var req request.GetTaskStatusRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	// 调用服务层查询任务状态：传入上下文、当前登录用户ID、任务ID
	resp, err := ctl.taskService.GetTaskStatus(ctx.Request.Context(), middle.CurrentUserID(ctx), req.TaskID)
	if err != nil {
		common.ErrorResponse(ctx, err)
		return
	}

	common.Success(ctx, resp)
}
