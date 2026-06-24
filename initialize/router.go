package initialize

import (
	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/internal/router"
)

func InitRouter() (*gin.Engine, error) {
	return router.InitRouter()
}
