package initialize

import (
	"context"
	"fmt"
	"time"

	"Go-AIServiceSupport/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm(ctx context.Context, cfg *config.Config,) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(cfg.MySQL.DSN),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get mysql connection pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}