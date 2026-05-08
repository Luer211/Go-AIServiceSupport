package main

import (
	"log"

	"Go-AIServiceSupport/global"
	"Go-AIServiceSupport/initialize"
	"Go-AIServiceSupport/internal/router"
)

func main() {
	// 把全局依赖安装好
	if err := initialize.GlobalInit(); err != nil {
		log.Fatal(err)
	}

	// 初始化路由
	r := router.InitRouter()
	if err := r.Run(":" + global.AppConfig().Server.Port); err != nil {
		log.Fatal(err)
	}
}
