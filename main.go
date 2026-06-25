package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Go-AIServiceSupport/global"
	"Go-AIServiceSupport/initialize"
	"Go-AIServiceSupport/internal/router"
)

// Todo: 这一段代码涉及了 channel、goroutine、context，是学习的好资料。

func main() {
	if err := run(); err != nil {
		log.Printf("application stopped: %v", err)
		os.Exit(1)
	}
}

func run() (err error) {
	// 创建可被系统信号取消的根 Context
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// 初始化基础组件
	if err := initialize.GlobalInit(ctx); err != nil {
		return fmt.Errorf("global initialization: %w", err)
	}

	// 保证 run 返回前关闭所有资源
	defer func() {
		err = errors.Join(err, initialize.GlobalClose())
	}()

	// 组装业务依赖和路由
	engine, err := router.InitRouter()
	if err != nil {
		return fmt.Errorf("initialize router: %w", err)
	}

	// 用标准库 http.Server 承载 Gin engine
	server := &http.Server{
		Addr:    ":" + global.AppConfig().Server.Port,
		Handler: engine,
	}

	serverErr := make(chan error, 1)

	// HTTP 服务是阻塞式的，所以放到 goroutine 中运行
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	// 等待：服务异常退出，或者收到系统终止信号
	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)

	case <-ctx.Done():
		// 根 ctx 已经取消，所以关闭服务要使用新的 ctx
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		return nil
	}
}