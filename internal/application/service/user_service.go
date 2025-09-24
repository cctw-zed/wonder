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
	repo  user.UserRepository
	idGen id.Generator
	log   logger.Logger
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

// GetProfile retrieves user profile by ID
func (s *userService) GetProfile(ctx context.Context, id string) (*user.User, error) {
	s.log.Info(ctx, "getting user profile", "user_id", id)

	if id == "" {
		s.log.Warn(ctx, "user ID is required")
		return nil, errors.NewRequiredFieldError("id", id)
	}

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error(ctx, "failed to get user profile", "error", err, "user_id", id)
		return nil, err
	}

	if u == nil {
		s.log.Warn(ctx, "user not found", "user_id", id)
		return nil, errors.NewEntityNotFoundError("user", id)
	}

	s.log.Info(ctx, "user profile retrieved successfully", "user_id", id)
	return u, nil
}

// UpdateProfile updates user profile information
func (s *userService) UpdateProfile(ctx context.Context, id string, req *user.UpdateProfileRequest) (*user.User, error) {
	s.log.Info(ctx, "updating user profile", "user_id", id)

	if id == "" {
		return nil, errors.NewRequiredFieldError("id", id)
	}

	if req == nil {
		return nil, errors.NewRequiredFieldError("request", "nil")
	}

	// Get existing user
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error(ctx, "failed to get user for update", "error", err, "user_id", id)
		return nil, err
	}

	if u == nil {
		s.log.Warn(ctx, "user not found for update", "user_id", id)
		return nil, errors.NewEntityNotFoundError("user", id)
	}

	// Update fields if provided
	if req.Name != "" {
		if err := u.UpdateName(ctx, req.Name); err != nil {
			s.log.Warn(ctx, "failed to update user name", "error", err, "user_id", id)
			return nil, err
		}
	}

	if req.Email != "" {
		// Check if new email already exists (but not for the same user)
		existingUser, err := s.repo.GetByEmail(ctx, req.Email)
		if err != nil {
			s.log.Error(ctx, "failed to check existing email", "error", err, "email", req.Email)
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			s.log.Warn(ctx, "email already exists for another user", "email", req.Email, "existing_user_id", existingUser.ID)
			return nil, errors.NewDuplicateEntryError("user", "email", req.Email, existingUser.ID)
		}

		if err := u.UpdateEmail(ctx, req.Email); err != nil {
			s.log.Warn(ctx, "failed to update user email", "error", err, "user_id", id)
			return nil, err
		}
	}

	// Update timestamp
	u.UpdatedAt = time.Now()

	// Validate the updated aggregate
	if err := u.Validate(ctx); err != nil {
		s.log.Warn(ctx, "user aggregate validation failed after update", "error", err, "user_id", id)
		return nil, err
	}

	// Persist the updated user
	if err := s.repo.Update(ctx, u); err != nil {
		s.log.Error(ctx, "failed to persist user update", "error", err, "user_id", id)
		return nil, err
	}

	s.log.Info(ctx, "user profile updated successfully", "user_id", id)
	return u, nil
}

// ListUsers retrieves a list of users with pagination and filtering
func (s *userService) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	if req == nil {
		return nil, errors.NewRequiredFieldError("request", "nil")
	}

	if s.log.DebugEnabled() {
		s.log.Debug(ctx, "listing users", "page", req.Page, "page_size", req.PageSize)
	}

	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	response, err := s.repo.List(ctx, req)
	if err != nil {
		s.log.Error(ctx, "failed to list users", "error", err)
		return nil, err
	}

	s.log.Info(ctx, "users listed successfully", "total", response.Total, "page", response.Page, "returned", len(response.Users))
	return response, nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	s.log.Info(ctx, "deleting user", "user_id", id)

	if id == "" {
		return errors.NewRequiredFieldError("id", id)
	}

	// Check if user exists before deleting
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error(ctx, "failed to get user for deletion", "error", err, "user_id", id)
		return err
	}

	if u == nil {
		s.log.Warn(ctx, "user not found for deletion", "user_id", id)
		return errors.NewEntityNotFoundError("user", id)
	}

	// Delete the user
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error(ctx, "failed to delete user", "error", err, "user_id", id)
		return err
	}

	s.log.Info(ctx, "user deleted successfully", "user_id", id)
	return nil
}
