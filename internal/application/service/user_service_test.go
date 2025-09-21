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
	idMocks "github.com/cctw-zed/wonder/pkg/snowflake/id/mocks"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)

	service := NewUserService(mockRepo, mockIDGen)

	tests := []struct {
		name      string
		email     string
		userName  string
		setupMock func()
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			userName: "Test User",
			setupMock: func() {
				// Expect email validation (check if exists)
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, nil).
					Times(1)

				// Expect ID generation
				mockIDGen.EXPECT().
					Generate().
					Return("test-id-123").
					Times(1)

				// Expect user creation
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *user.User) error {
						assert.Equal(t, "test-id-123", u.ID)
						assert.Equal(t, "test@example.com", u.Email)
						assert.Equal(t, "Test User", u.Name)
						assert.False(t, u.CreatedAt.IsZero())
						assert.False(t, u.UpdatedAt.IsZero())
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:     "empty email",
			email:    "",
			userName: "Test User",
			setupMock: func() {
				// No mock expectations as validation should fail early
			},
			wantErr: true,
			errMsg:  "invalid email: email is required",
		},
		{
			name:     "invalid email format",
			email:    "invalid-email",
			userName: "Test User",
			setupMock: func() {
				// No mock expectations as validation should fail early
			},
			wantErr: true,
			errMsg:  "invalid email: invalid email format",
		},
		{
			name:     "user already exists",
			email:    "existing@example.com",
			userName: "Existing User",
			setupMock: func() {
				existingUser := &user.User{
					ID:    "existing-id",
					Email: "existing@example.com",
					Name:  "Existing User",
				}
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "existing@example.com").
					Return(existingUser, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name:     "repository create fails",
			email:    "test@example.com",
			userName: "Test User",
			setupMock: func() {
				// Expect email validation (check if exists)
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, nil).
					Times(1)

				// Expect ID generation
				mockIDGen.EXPECT().
					Generate().
					Return("test-id-123").
					Times(1)

				// Expect user creation to fail
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to create user: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := service.Register(context.Background(), tt.email, tt.userName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "test-id-123", result.ID)
				assert.Equal(t, tt.email, result.Email)
				assert.Equal(t, tt.userName, result.Name)
				assert.False(t, result.CreatedAt.IsZero())
				assert.False(t, result.UpdatedAt.IsZero())
			}
		})
	}
}

func TestUserService_validateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)

	service := &userService{
		repo:  mockRepo,
		idGen: mockIDGen,
	}

	tests := []struct {
		name    string
		email   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name:    "invalid email format",
			email:   "invalid-email",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "email without domain",
			email:   "test@",
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name:    "email without local part",
			email:   "@example.com",
			wantErr: true,
			errMsg:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateEmail(tt.email)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)

	service := NewUserService(mockRepo, mockIDGen)

	assert.NotNil(t, service)

	// Verify that the service implements the interface
	var _ user.UserService = service
}

// Integration test helper
func createTestUser() *user.User {
	return &user.User{
		ID:        "test-id-123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
