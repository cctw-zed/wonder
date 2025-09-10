package service

import (
	"context"
	"fmt"
	"github.com/cctw-zed/wonder/internal/domain/user"
	"time"
)

type userService struct {
	repo  user.UserRepository
	idGen snowflake.Generator
}

func NewUserService(repo user.UserRepository, idGen snowflake.Generator) user.UserService {
	return &userService{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *userService) Register(ctx context.Context, email, name string) (*user.User, error) {
	// 业务规则验证
	if err := s.validateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// 检查邮箱是否已存在
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// 创建用户
	u := &user.User{
		ID:        s.idGen.Generate(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (s *userService) validateEmail(email string) error {
	// TODO 待实现
	return nil
}
