package initialize

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/global"
)

// 按顺序把全局依赖安装好
func GlobalInit() error {
	cfg, err := config.Load("config/application-dev.yaml")
	if err != nil {
		return err
	}

	global.Config = cfg
	global.Log = InitLogger()
	global.DB = InitGorm(cfg)
	global.Redis = InitRedis(cfg)
	global.TaskProducer = InitMQ(cfg)

	return nil
}
