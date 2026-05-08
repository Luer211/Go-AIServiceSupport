package router

import (
	"time"

	"github.com/gin-gonic/gin"

	"Go-AIServiceSupport/global"
	"Go-AIServiceSupport/internal/api/controller"
	"Go-AIServiceSupport/internal/cache"
	"Go-AIServiceSupport/internal/mq"
	"Go-AIServiceSupport/internal/repository/dao"
	"Go-AIServiceSupport/internal/service"
	"Go-AIServiceSupport/middle"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 未加入日志中间件

	cfg := global.AppConfig()
	producer := global.TaskProducer
	if producer == nil {
		producer = mq.NewNoopProducer()
	}

	userDao := dao.NewUserDao(global.DB)
	taskDao := dao.NewTaskDao(global.DB)
	taskStatusCache := cache.NewTaskStatusCache(global.Redis, cfg.Task.RedisStatusTTLSeconds)
	rateLimiter := cache.NewRateLimiter(global.Redis, time.Minute)

	authService := service.NewAuthService(userDao, cfg.JWT.Secret, cfg.JWT.ExpireSeconds)
	taskService := service.NewTaskService(taskDao, taskStatusCache, producer)

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

	return r
}
