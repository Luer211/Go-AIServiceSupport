package global

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/internal/mq"
	"Go-AIServiceSupport/logger"
)

var (
	Config       *config.Config
	DB           any
	Redis        any
	Log          logger.Logger
	TaskProducer mq.Producer
)

func AppConfig() *config.Config {
	if Config == nil {
		return config.Default()
	}
	return Config
}
