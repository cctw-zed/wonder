package service

import (
	"context"
	"time"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/jwt"
	"github.com/cctw-zed/wonder/pkg/logger"
)

// AuthService provides authentication functionality
type AuthService interface {
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*jwt.Claims, error)
}

// LoginResponse represents the response for login
type LoginResponse struct {
	User        *user.User `json:"user"`
	AccessToken string     `json:"access_token"`
	TokenType   string     `json:"token_type"`
	ExpiresIn   int64      `json:"expires_in"`
}

type authService struct {
	userService  user.UserService
	tokenService jwt.TokenService
	log          logger.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(userService user.UserService, tokenService jwt.TokenService) AuthService {
	return NewAuthServiceWithLogger(userService, tokenService, logger.Get().WithLayer("application").WithComponent("auth_service"))
}

func NewAuthServiceWithLogger(userService user.UserService, tokenService jwt.TokenService, log logger.Logger) AuthService {
	if userService == nil {
		panic("user service cannot be nil")
	}
	if tokenService == nil {
		panic("token service cannot be nil")
	}
	if log == nil {
		panic("logger cannot be nil")
	}

	return &authService{
		userService:  userService,
		tokenService: tokenService,
		log:          log,
	}
}

// Login authenticates user and returns access token
func (s *authService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	s.log.Info(ctx, "processing login request", "email", email)

	// Authenticate user
	u, err := s.userService.Login(ctx, email, password)
	if err != nil {
		s.log.Warn(ctx, "login failed", "error", err, "email", email)
		return nil, err
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateToken(u.ID)
	if err != nil {
		s.log.Error(ctx, "failed to generate access token", "error", err, "user_id", u.ID)
		return nil, err
	}

	s.log.Info(ctx, "login successful", "user_id", u.ID, "email", email)

	return &LoginResponse{
		User:        u,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour.Seconds()), // TODO: Make configurable
	}, nil
}

// Logout invalidates the access token
func (s *authService) Logout(ctx context.Context, token string) error {
	s.log.Info(ctx, "processing logout request")

	if token == "" {
		return errors.NewRequiredFieldError("token", token)
	}

	// Validate token first
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		s.log.Warn(ctx, "logout with invalid token", "error", err)
		return err
	}

	s.log.Info(ctx, "logout successful", "user_id", claims.UserID)

	// TODO: Add token blacklist/invalidation mechanism
	// For now, we just log the logout - token will naturally expire
	return nil
}

// ValidateToken validates an access token and returns claims
func (s *authService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	if s.log.DebugEnabled() {
		s.log.Debug(ctx, "validating token")
	}

	if token == "" {
		return nil, errors.NewRequiredFieldError("token", token)
	}

	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		if s.log.DebugEnabled() {
			s.log.Debug(ctx, "token validation failed", "error", err)
		}
		return nil, err
	}

	if s.log.DebugEnabled() {
		s.log.Debug(ctx, "token validation successful", "user_id", claims.UserID)
	}

	return claims, nil
}