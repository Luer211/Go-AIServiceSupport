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

var ErrInvalidLogin = errors.New("invalid username or password")

type AuthService struct {
	users            *dao.UserDao
	jwtSecret        string
	jwtExpireSeconds int64
}

func NewAuthService(users *dao.UserDao, jwtSecret string, jwtExpireSeconds int64) *AuthService {
	return &AuthService{
		users:            users,
		jwtSecret:        jwtSecret,
		jwtExpireSeconds: jwtExpireSeconds,
	}
}

func (s *AuthService) Register(ctx context.Context, req request.RegisterRequest) (*response.RegisterResponse, error) {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return &response.RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidLogin
	}

	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpireSeconds)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:     token,
		ExpiresIn: s.jwtExpireSeconds,
	}, nil
}
