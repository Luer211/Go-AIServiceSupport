package middle

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/internal/cache"
)

// Todo：全局错误的实现
func RateLimit(limiter *cache.RateLimiter, ipLimit int, userLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		// IP 限流
		allowed, err := limiter.Allow(c.Request.Context(), cache.IPRateKey(c.ClientIP()), ipLimit)
		if err != nil {
			common.FailWithMessage(c, e.CodeInternalError, err.Error())
			c.Abort()
			return
		}
		if !allowed {
			common.Fail(c, e.CodeTooManyReq)
			c.Abort()
			return
		}

		// 用户ID 限流
		userID := CurrentUserID(c)
		if userID != 0 {
			allowed, err = limiter.Allow(c.Request.Context(), cache.UserRateKey(strconv.FormatUint(userID, 10)), userLimit)
			if err != nil {
				common.FailWithMessage(c, e.CodeInternalError, err.Error())
				c.Abort()
				return
			}
			if !allowed {
				common.Fail(c, e.CodeTooManyReq)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
