package initialize

import "Go-AIServiceSupport/logger"

func InitLogger() logger.Logger {
	// TODO: 替换成 Zap 初始化逻辑。
	return logger.New()
}
