// 安装全局依赖

package initialize

import (
	"context"
	"errors"
	"database/sql"
	"fmt"

	"Go-AIServiceSupport/config"
	"Go-AIServiceSupport/global"
)

// 按顺序把全局依赖安装好
func GlobalInit(ctx context.Context) error {
	// 1. 加载配置文件
	cfg, err := config.Load("config/application-dev.yaml")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. 初始化日志系统
	appLog, err := InitLogger()
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	// 3. 初始化 MySQL 数据库
	db, err := InitGorm(ctx, cfg)
	if err != nil {
		// 关闭初始化了的日志系统
		_ = appLog.Sync()
		return fmt.Errorf("initialize mysql: %w", err)
	}

	// 4. 初始化 Redis
	redisClient, err := InitRedis(ctx, cfg)
	if err != nil {
		// 关闭数据库，关闭日志系统
		closeGorm(db)
		_ = appLog.Sync()
		return fmt.Errorf("initialize redis: %w", err)
	}

	// 5. 初始化 MQ
	producer, err := InitMQ(cfg)
	if err != nil {
		// 关闭redis，关闭数据库，关闭日志系统
		_ = redisClient.Close()
		closeGorm(db)
		_ = appLog.Sync()
		return fmt.Errorf("initialize mq: %w", err)
	}

	// 所有依赖都初始化成功后，才写入全局变量
	global.Config = cfg
	global.Log = appLog
	global.DB = db
	global.Redis = redisClient
	global.TaskProducer = producer

	return nil
}

// GlobalClose 用于关闭应用运行期间打开的全局资源。
// 通常在服务退出、收到信号、main 函数 defer 时调用。
func GlobalClose() error {
	var errs []error

	// 1. 关闭消息队列生产者
	if global.TaskProducer != nil {
		if err := global.TaskProducer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close mq: %w", err))
		}
	}

	// 2. 关闭 Redis 客户端
	if global.Redis != nil {
		if err := global.Redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close redis: %w", err))
		}
	}

	// 3. 关闭 MySQL 连接池
	if global.DB != nil {
		if err := closeGorm(global.DB); err != nil {
			errs = append(errs, err)
		}
	}

	// 4. 同步日志缓冲区
	if global.Log != nil {
		// zap 的 Sync 在 stdout/stderr 上有时会返回
		// inappropriate ioctl 之类的错误。
		// 这种错误通常不影响程序退出，所以这里忽略。
		_ = global.Log.Sync()
	}

	// errors.Join 会把多个错误合并成一个 error。
	// 如果 errs 为空，则返回 nil。
	return errors.Join(errs...)
}

// closeGorm 用于关闭 GORM 底层的数据库连接池。
func closeGorm(db interface {
	DB() (*sql.DB, error)
}) error {
	// 从 GORM 对象中取出底层的 database/sql 连接池。
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get mysql connection pool: %w", err)
	}

	// 关闭数据库连接池。
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close mysql: %w", err)
	}

	return nil
}
