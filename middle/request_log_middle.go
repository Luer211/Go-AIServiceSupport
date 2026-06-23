package middle

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Go-AIServiceSupport/global"
)

// 在每次请求处理完成后记录一条访问日志
func RequestLog() gin.HandlerFunc {
	return  func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		userID := CurrentUserID(c)

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
		}

		// 补充可选字段：假如有用户id（用户已登录），追加 user_id 字段
		if userID != 0 {
			fields = append(fields, zap.Uint64("user_id", userID))
		}

		// Todo: 按不同情况输出不同级别的日志，目前这些是示例
		// 后面可能想要细分成就是说，不同的业务错误码，需要不同的错误？
		// 然后还有一个就是说，我们的HTTP状态码目前设计有问题，全部响应是200，这是不对的

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
			global.Log.Error("request completed with errors", fields...)
			return
		}

		global.Log.Info("request completed", fields...)
	}
}