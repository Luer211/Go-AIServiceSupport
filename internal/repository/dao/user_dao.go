package dao

import (
	"context"

	"Go-AIServiceSupport/internal/model"
)

type UserDao struct {
	db any
}

func NewUserDao(db any) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) Create(ctx context.Context, user *model.User) error {
	// TODO: 使用 GORM 写入 users 表。
	return ErrNotImplemented
}

func (d *UserDao) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	// TODO: 使用 GORM 按 username 查询用户。
	return nil, ErrNotImplemented
}
