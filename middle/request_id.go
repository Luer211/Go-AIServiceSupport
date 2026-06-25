// request_id 可以把相关日志串起来

package middle

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextRequestIDKey = "request_id"
	HeaderRequestID    = "X-Request-ID"
)

func ensureRequestID(c *gin.Context) string {
	requestID := c.GetHeader(HeaderRequestID)
	if requestID == "" {
		requestID = newRequestID()
	}

	// 把 requestID 写进 gin.Context 里面
	c.Set(ContextRequestIDKey, requestID)
	c.Header(HeaderRequestID, requestID)

	return requestID
}

func currentRequestID(c *gin.Context) string {
	if value, exists := c.Get(ContextRequestIDKey); exists {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}

	return time.Now().Format("20060102150405.000000000")
}
