package initialize

import "Go-AIServiceSupport/logger"

func InitLogger() (logger.Logger, error) {
	return logger.New()
}
