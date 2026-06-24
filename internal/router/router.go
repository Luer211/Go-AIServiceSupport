package router

import (
	"time"
	"fmt"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/global"
	"Go-AIServiceSupport/internal/api/controller"
	"Go-AIServiceSupport/internal/cache"
	"Go-AIServiceSupport/internal/repository/dao"
	"Go-AIServiceSupport/internal/service"
	"Go-AIServiceSupport/middle"
)

func InitRouter() (*gin.Engine, error) {
	cfg := global.AppConfig()

	if global.DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}
	if global.Redis == nil {
		return nil, fmt.Errorf("redis is not initialized")
	}
	if global.TaskProducer == nil {
		return nil, fmt.Errorf("task producer is not initialized")
	}
	if global.Log == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	r := gin.New()

	r.Use(middle.RequestLog())
	r.Use(middle.RecoveryLog())

	userDao := dao.NewUserDao(global.DB)
	taskDao := dao.NewTaskDao(global.DB)

	taskStatusCache := cache.NewTaskStatusCache(
		global.Redis,
		cfg.Task.RedisStatusTTLSeconds,
	)

	rateLimiter := cache.NewRateLimiter(
		global.Redis,
		time.Minute,
	)

	authService := service.NewAuthService(
		userDao,
		cfg.JWT.Secret,
		cfg.JWT.ExpireSeconds,
	)

	taskService, err := service.NewTaskService(
		taskDao,
		taskStatusCache,
		global.TaskProducer,
	)
	if err != nil {
		return nil, fmt.Errorf("initialize task service: %w", err)
	}

	authController := controller.NewAuthController(authService)
	taskController := controller.NewTaskController(taskService)
	
	api := r.Group("/api/v1")

	// 公共接口：注册、登录
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// 私有接口：发出任务请求、轮询任务结果（需要 JWT）
	tasks := api.Group("/tasks")
	tasks.Use(middle.VerifyJWT(cfg.JWT.Secret))
	tasks.Use(middle.RateLimit(rateLimiter, cfg.RateLimit.IPLimitPerMinute, cfg.RateLimit.UserLimitPerMinute))
	{
		tasks.POST("", taskController.CreateTask)
		tasks.GET("/:task_id", taskController.GetTaskStatus)
	}

	return r, nil
}
