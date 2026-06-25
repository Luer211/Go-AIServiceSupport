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
		
		// 把 requestID 写进 gin.Context 里面
		requestID := ensureRequestID(c)

		c.Next()

		latency := time.Since(start)
		userID := CurrentUserID(c)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("route", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.String("client_ip", c.ClientIP()),
		}

		// 补充可选字段：假如有用户id（用户已登录），追加 user_id 字段
		if userID != 0 {
			fields = append(fields, zap.Uint64("user_id", userID))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			global.Log.Error("request completed", fields...)
		case status >= 400:
			global.Log.Warn("request completed", fields...)
		default:
			global.Log.Info("request completed", fields...)
		}
	}
}