package service

import (
	"context"
	"time"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/logger"
	"github.com/cctw-zed/wonder/pkg/snowflake/id"
)

type userService struct {
	repo user.UserRepository
	idGen id.Generator
	log  logger.Logger
}

func NewUserService(repo user.UserRepository, idGen id.Generator) user.UserService {
	return NewUserServiceWithLogger(repo, idGen, logger.Get().WithLayer("application").WithComponent("user_service"))
}

func NewUserServiceWithLogger(repo user.UserRepository, idGen id.Generator, log logger.Logger) user.UserService {
	if repo == nil {
		panic("user repository cannot be nil")
	}
	if idGen == nil {
		panic("ID generator cannot be nil")
	}
	if log == nil {
		panic("logger cannot be nil")
	}

	return &userService{
		repo:  repo,
		idGen: idGen,
		log:   log,
	}
}

func (s *userService) Register(ctx context.Context, email, name string) (*user.User, error) {
	s.log.Info(ctx, "registering user", "email", email, "name", name)

	// Business rule validation
	if err := s.validateEmail(ctx, email); err != nil {
		s.log.Warn(ctx, "email validation failed", "error", err, "email", email)
		return nil, err
	}

	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error(ctx, "failed to check existing email", "error", err, "email", email)
		return nil, err
	}
	if existingUser != nil {
		s.log.Warn(ctx, "email already exists", "email", email, "existing_user_id", existingUser.ID)
		return nil, errors.NewDuplicateEntryError("user", "email", email, existingUser.ID)
	}

	// Create user aggregate
	userID := s.idGen.Generate()
	u := &user.User{
		ID:        userID,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if s.log.DebugEnabled() {
		s.log.Debug(ctx, "created user aggregate", "user_id", userID, "email", email, "name", name)
	}

	// Validate the aggregate before persisting
	if err := u.Validate(ctx); err != nil {
		s.log.Warn(ctx, "user aggregate validation failed", "error", err, "user_id", userID)
		return nil, err
	}

	// Persist the user
	if err := s.repo.Create(ctx, u); err != nil {
		s.log.Error(ctx, "failed to persist user", "error", err, "user_id", userID)
		return nil, err
	}

	s.log.Info(ctx, "user registered successfully", "user_id", userID, "email", email)
	return u, nil
}

func (s *userService) validateEmail(ctx context.Context, email string) error {
	if s.log.DebugEnabled() {
		s.log.Debug(ctx, "validating email", "email", email)
	}

	if email == "" {
		return errors.NewRequiredFieldError("email", email)
	}

	u := &user.User{Email: email}
	if !u.IsEmailValid() {
		return errors.NewInvalidFormatError("email", email, "valid email address")
	}

	return nil
}
