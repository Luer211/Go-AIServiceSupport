package service

import (
	"context"
	"errors"

	"Go-AIServiceSupport/common/utils"
	"Go-AIServiceSupport/internal/api/request"
	"Go-AIServiceSupport/internal/api/response"
	"Go-AIServiceSupport/internal/model"
	"Go-AIServiceSupport/internal/repository/dao"
)

// 全局错误变量：用户名不存在/密码错误 时返回这个错误
var ErrInvalidLogin = errors.New("invalid username or password")

type AuthService struct {
	users            *dao.UserDao
	jwtSecret        string			// JWT 签名密钥
	jwtExpireSeconds int64			// JWT 过期时间
}

func NewAuthService(users *dao.UserDao, jwtSecret string, jwtExpireSeconds int64) *AuthService {
	return &AuthService{
		users:            users,
		jwtSecret:        jwtSecret,
		jwtExpireSeconds: jwtExpireSeconds,
	}
}

// 用户注册服务
func (s *AuthService) Register(ctx context.Context, req request.RegisterRequest) (*response.RegisterResponse, error) {
	// 密码加密
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 构造用户对象
	user := &model.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	// 存储到数据库
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	// 返回结果
	return &response.RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// 用户登录服务
func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	// 根据用户名查询用户
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// 校验密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidLogin
	}

	// 生成 JWT 令牌
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpireSeconds)
	if err != nil {
		return nil, err
	}

	// 返回 token 和过期时间
	return &response.LoginResponse{
		Token:     token,
		ExpiresIn: s.jwtExpireSeconds,
	}, nil
}
