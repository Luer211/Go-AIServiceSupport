package global

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/internal/mq"
	"Go-AIServiceSupport/logger"

	"gorm.io/gorm"
)

// Todo: 这里要明确一下，我们已经严肃设计了DB和Redis了，就不应该是any了
var (
	Config       *config.Config
	DB           *gorm.DB
	Redis        any
	Log          logger.Logger
	TaskProducer mq.Producer
)

// Todo: 即便发生了 panic 我们也应该有一个日志中间件捕获错误统一错误和打印日志
// 这个也就是我们的 recovery_log_middle 要做的
func AppConfig() *config.Config {
	if Config == nil {
		panic("global config not initialized")
	}
	return Config
}
