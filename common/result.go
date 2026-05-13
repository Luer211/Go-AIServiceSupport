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

// 失败响应：使用预设错误码
func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: e.Message(code),
	})
}

// 失败响应：自定义错误信息
func FailWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
