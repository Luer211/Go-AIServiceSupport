package dao

import (
	"context"

	"Go-AIServiceSupport/internal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

// 创建新用户
func (d *UserDao) Create(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

// 按照用户名称查询用户
func (d *UserDao) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := d.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
			return nil, err
		}
	return &user, nil
}
