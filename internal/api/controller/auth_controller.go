package controller

import (
	"errors"

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

func (ctl *AuthController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	resp, err := ctl.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		common.FailWithMessage(ctx, e.CodeInternalError, err.Error())
		return
	}

	common.Success(ctx, resp)
}

func (ctl *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.Fail(ctx, e.CodeInvalidParams)
		return
	}

	resp, err := ctl.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLogin) {
			common.Fail(ctx, e.CodeInvalidLogin)
			return
		}
		common.FailWithMessage(ctx, e.CodeInternalError, err.Error())
		return
	}

	common.Success(ctx, resp)
}
