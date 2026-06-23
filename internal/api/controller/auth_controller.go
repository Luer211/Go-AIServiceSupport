package controller

import (
	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/internal/api/request"
	"Go-AIServiceSupport/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// 用户注册接口
func (ctl *AuthController) Register(ctx *gin.Context) {
	// 绑定并校验前端传入的 JSON 注册请求参数
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	// 调用服务层处理业务逻辑
	resp, err := ctl.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(ctx, err)
		return
	}

	common.Success(ctx, resp)
}

// 用户登录接口
func (ctl *AuthController) Login(ctx *gin.Context) {
	// 绑定并校验前端传入的 JSON 注册请求参数
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	// 调用服务层处理业务逻辑
	resp, err := ctl.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(ctx, err)
		return
	}

	common.Success(ctx, resp)
}
