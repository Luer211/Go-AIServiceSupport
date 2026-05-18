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

// 任务创建接口
func (ctl *TaskController) CreateTask(ctx *gin.Context) {
	var req request.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Error(ctx, err)
		return
	}

	// 调用服务层创建任务：传入上下文、当前登录用户ID、请求参数
	resp, err := ctl.taskService.CreateTask(ctx.Request.Context(), middle.CurrentUserID(ctx), req)
	if err != nil {
		common.FailWithMessage(ctx, e.CodeTaskSubmitFailed, err.Error())
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
		// 判断是否为无权限访问错误，若是则返回权限不足响应
		if errors.Is(err, service.ErrTaskForbidden) {
			common.Error(ctx, err)
			return
		}
		// 其他错误（如任务不存在），返回任务未找到响应
		common.FailWithMessage(ctx, e.CodeTaskNotFound, err.Error())
		return
	}

	common.Success(ctx, resp)
}
