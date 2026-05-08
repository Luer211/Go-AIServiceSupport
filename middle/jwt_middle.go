package middle

import (
	"strings"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/common"
	"Go-AIServiceSupport/common/e"
	"Go-AIServiceSupport/common/utils"
)

const (
	ContextUserIDKey   = "user_id"
	ContextUsernameKey = "username"
)

func VerifyJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.Fail(c, e.CodeUnauthorized)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader || token == "" {
			common.Fail(c, e.CodeUnauthorized)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token, secret)
		if err != nil {
			common.Fail(c, e.CodeUnauthorized)
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Next()
	}
}

func VerifyJWTAdmin() gin.HandlerFunc {
	return VerifyJWT("dev-secret")
}

func CurrentUserID(c *gin.Context) uint64 {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}

	switch userID := value.(type) {
	case uint64:
		return userID
	case uint:
		return uint64(userID)
	case int64:
		return uint64(userID)
	case int:
		return uint64(userID)
	default:
		return 0
	}
}
