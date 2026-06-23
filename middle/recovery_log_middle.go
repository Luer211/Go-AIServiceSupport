package middle

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/global"
)

// 捕获请求处理过程中未处理的 panic，防止整个服务进程崩溃
func RecoveryLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				global.Log.Error("panic recovered",
					zap.String("panic", fmt.Sprint(r)),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()), // 崩溃堆栈（定位代码用）
				)

				// 返回错误给前端
				common.Fail(c, e.CodeInternalError)
				c.Abort()
			}
		}()

		c.Next()
	}
}
