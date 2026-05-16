package initialize

import (
	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/global"
)

// 按顺序把全局依赖安装好
// Todo：我们可以看到这里日志做了错误处理，其他的是不是也应该做呢？
func GlobalInit() error {
	cfg, err := config.Load("config/application-dev.yaml")
	if err != nil {
		return err
	}

	global.Config = cfg

	log, err := InitLogger()
	if err != nil {
		return err
	}
	global.Log = log

	global.DB = InitGorm(cfg)

	global.Redis = InitRedis(cfg)

	global.TaskProducer = InitMQ(cfg)

	return nil
}
