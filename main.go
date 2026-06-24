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

func main() {
	if err := run(); err != nil {
		log.Printf("application stopped: %v", err)
		os.Exit(1)
	}
}

func run() (err error) {
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

	server := &http.Server{
		Addr:    ":" + global.AppConfig().Server.Port,
		Handler: engine,
	}

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

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