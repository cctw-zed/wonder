package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/domain/user/mocks"
	apperrors "github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/jwt"
	"github.com/cctw-zed/wonder/pkg/logger"
)

func TestNewAuthService(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)

	// Create a real JWT service for testing
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 24*time.Hour)

	authService := NewAuthService(mockUserService, tokenService)
	require.NotNil(t, authService)
}

func TestAuthService_Login(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 24*time.Hour)
	authService := NewAuthService(mockUserService, tokenService)

	tests := []struct {
		name      string
		email     string
		password  string
		setupMock func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			setupMock: func() {
				testUser := &user.User{
					ID:        "user123",
					Email:     "test@example.com",
					Name:      "Test User",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				mockUserService.EXPECT().
					Login(gomock.Any(), "test@example.com", "password123").
					Return(testUser, nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), "", "password123").
					Return(nil, apperrors.NewRequiredFieldError("email", "")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name:     "empty password",
			email:    "test@example.com",
			password: "",
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), "test@example.com", "").
					Return(nil, apperrors.NewRequiredFieldError("password", "")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name:     "user service login fails",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), "test@example.com", "wrongpassword").
					Return(nil, errors.New("invalid credentials")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMock: func() {
				mockUserService.EXPECT().
					Login(gomock.Any(), "nonexistent@example.com", "password123").
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			response, err := authService.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, "user123", response.User.ID)
				assert.Equal(t, "test@example.com", response.User.Email)
				assert.NotEmpty(t, response.AccessToken)
				assert.Equal(t, "Bearer", response.TokenType)
				assert.Equal(t, int64(24*60*60), response.ExpiresIn) // 24 hours in seconds

				// Verify the token can be validated
				claims, err := tokenService.ValidateToken(response.AccessToken)
				require.NoError(t, err)
				assert.Equal(t, "user123", claims.UserID)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 24*time.Hour)
	authService := NewAuthService(mockUserService, tokenService)

	// Generate a valid token for testing
	validToken, err := tokenService.GenerateToken("user123")
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful logout with valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "logout with empty token",
			token:   "",
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name:    "logout with invalid token",
			token:   "invalid.token",
			wantErr: true,
			errMsg:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Logout(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 24*time.Hour)
	authService := NewAuthService(mockUserService, tokenService)

	// Generate a valid token for testing
	validToken, err := tokenService.GenerateToken("user123")
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		setupMock func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:  "successful token validation",
			token: validToken,
			setupMock: func() {
				// No mock expectations needed for token validation
			},
			wantErr: false,
		},
		{
			name:  "empty token",
			token: "",
			setupMock: func() {
				// No mock expectations as validation should fail early
			},
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name:  "invalid token",
			token: "invalid.token",
			setupMock: func() {
				// No mock expectations as validation should fail early
			},
			wantErr: true,
			errMsg:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			claims, err := authService.ValidateToken(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
				assert.Equal(t, "user123", claims.UserID)
			}
		})
	}
}

func TestAuthService_Integration(t *testing.T) {
	// Test the complete authentication flow
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 24*time.Hour)
	authService := NewAuthService(mockUserService, tokenService)

	testUser := &user.User{
		ID:        "user123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock expectations
	mockUserService.EXPECT().
		Login(gomock.Any(), "test@example.com", "password123").
		Return(testUser, nil).
		Times(1)

	ctx := context.Background()

	// Step 1: Login
	loginResponse, err := authService.Login(ctx, "test@example.com", "password123")
	require.NoError(t, err)
	require.NotNil(t, loginResponse)
	assert.NotEmpty(t, loginResponse.AccessToken)

	// Step 2: Use token to validate and get claims
	claims, err := authService.ValidateToken(ctx, loginResponse.AccessToken)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, testUser.ID, claims.UserID)

	// Step 3: Logout
	err = authService.Logout(ctx, loginResponse.AccessToken)
	require.NoError(t, err)
}

func TestAuthService_TokenExpiry(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	// Create a token service with very short expiry for testing
	tokenService := jwt.NewTokenService("test-signing-key-32-chars-minimum", 1*time.Millisecond)
	authService := NewAuthService(mockUserService, tokenService)

	testUser := &user.User{
		ID:        "user123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserService.EXPECT().
		Login(gomock.Any(), "test@example.com", "password123").
		Return(testUser, nil).
		Times(1)

	// Login to get a token
	loginResponse, err := authService.Login(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
	require.NotNil(t, loginResponse)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to validate expired token
	claims, err := authService.ValidateToken(context.Background(), loginResponse.AccessToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	assert.Nil(t, claims)
}
