package common

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common/e"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    e.CodeSuccess,
		Message: e.Message(e.CodeSuccess),
		Data:    data,
	})
}

func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: e.Message(code),
	})
}

func FailWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
