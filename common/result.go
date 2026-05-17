package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common/e"
)

// 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    e.CodeSuccess,
		Message: e.Message(e.CodeSuccess),
		Data:    data,
	})
}

// 失败响应
func Error(c *gin.Context, err error) {
	// 假如是定义的AppError的话
	if appErr, ok := AsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, Response{
			Code: 	 appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// 假如是普通的error的话
	c.JSON(http.StatusInternalServerError, Response{
		Code: 	 e.CodeInternalError,
		Message: e.Message(e.CodeInternalError),
	})
}

// 控制层发生错误的话用Fail和FailWithMessage

func Fail(c *gin.Context, code int) {
	Error(c, NewAppError(code))
}

func FailWithMessage(c *gin.Context, code int, message string) {
	Error(c, NewAppErrorWithMessage(code, message))
}
