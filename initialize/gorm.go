package initialize

import (
	"Go-AIServiceSupport/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm(cfg *config.Config) any {
	// TODO: 我们这里的错误处理尚未接入全局日志
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
