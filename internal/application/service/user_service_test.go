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
	"github.com/cctw-zed/wonder/pkg/logger"
	idMocks "github.com/cctw-zed/wonder/pkg/snowflake/id/mocks"
)

func TestUserService_Register(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

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
			errMsg:  "validation failed for field 'email': email is required",
		},
		{
			name:     "invalid email format",
			email:    "invalid-email",
			userName: "Test User",
			setupMock: func() {
				// No mock expectations as validation should fail early
			},
			wantErr: true,
			errMsg:  "invalid format for email, expected: valid email address",
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
			errMsg:  "email 'existing@example.com' already exists",
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
			errMsg:  "database error",
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
	// Initialize logger for tests
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)

	service := NewUserServiceWithLogger(mockRepo, mockIDGen, logger.Get().WithLayer("application").WithComponent("user_service")).(*userService)

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
			errMsg:  "invalid format for email, expected: valid email address",
		},
		{
			name:    "email without domain",
			email:   "test@",
			wantErr: true,
			errMsg:  "invalid format for email, expected: valid email address",
		},
		{
			name:    "email without local part",
			email:   "@example.com",
			wantErr: true,
			errMsg:  "invalid format for email, expected: valid email address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateEmail(context.Background(), tt.email)

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

func TestUserService_GetProfile(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	service := NewUserService(mockRepo, mockIDGen)

	testUser := createTestUser()

	tests := []struct {
		name          string
		userID        string
		mockBehavior  func()
		expectedUser  *user.User
		expectedError string
	}{
		{
			name:   "successful get profile",
			userID: "test-id-123",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
			},
			expectedUser: testUser,
		},
		{
			name:   "user not found",
			userID: "nonexistent-id",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "nonexistent-id").
					Return(nil, nil).
					Times(1)
			},
			expectedError: "not found",
		},
		{
			name:   "empty user ID",
			userID: "",
			mockBehavior: func() {
				// No mock calls expected
			},
			expectedError: "id is required",
		},
		{
			name:   "repository error",
			userID: "test-id-123",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := service.GetProfile(context.Background(), tt.userID)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUser, result)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	service := NewUserService(mockRepo, mockIDGen)

	testUser := createTestUser()

	tests := []struct {
		name          string
		userID        string
		request       *user.UpdateProfileRequest
		mockBehavior  func()
		expectedError string
	}{
		{
			name:   "successful update name only",
			userID: "test-id-123",
			request: &user.UpdateProfileRequest{
				Name: "Updated Name",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "successful update email only",
			userID: "test-id-123",
			request: &user.UpdateProfileRequest{
				Email: "updated@example.com",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "updated@example.com").
					Return(nil, nil).
					Times(1)
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "update both name and email",
			userID: "test-id-123",
			request: &user.UpdateProfileRequest{
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "updated@example.com").
					Return(nil, nil).
					Times(1)
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent-id",
			request: &user.UpdateProfileRequest{
				Name: "Updated Name",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "nonexistent-id").
					Return(nil, nil).
					Times(1)
			},
			expectedError: "not found",
		},
		{
			name:   "email already exists for another user",
			userID: "test-id-123",
			request: &user.UpdateProfileRequest{
				Email: "existing@example.com",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					GetByEmail(gomock.Any(), "existing@example.com").
					Return(&user.User{ID: "another-user-id"}, nil).
					Times(1)
			},
			expectedError: "already exists",
		},
		{
			name:   "empty user ID",
			userID: "",
			request: &user.UpdateProfileRequest{
				Name: "Updated Name",
			},
			mockBehavior: func() {
				// No mock calls expected
			},
			expectedError: "id is required",
		},
		{
			name:    "nil request",
			userID:  "test-id-123",
			request: nil,
			mockBehavior: func() {
				// No mock calls expected
			},
			expectedError: "request is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := service.UpdateProfile(context.Background(), tt.userID, tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	service := NewUserService(mockRepo, mockIDGen)

	testUsers := []*user.User{createTestUser()}
	testResponse := &user.ListUsersResponse{
		Users:      testUsers,
		Total:      1,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	tests := []struct {
		name          string
		request       *user.ListUsersRequest
		mockBehavior  func()
		expectedError string
	}{
		{
			name: "successful list users",
			request: &user.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(testResponse, nil).
					Times(1)
			},
		},
		{
			name:    "nil request",
			request: nil,
			mockBehavior: func() {
				// No mock calls expected
			},
			expectedError: "request is required",
		},
		{
			name: "default pagination values",
			request: &user.ListUsersRequest{
				Page:     0, // Should be set to 1
				PageSize: 0, // Should be set to 10
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(testResponse, nil).
					Times(1)
			},
		},
		{
			name: "repository error",
			request: &user.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			mockBehavior: func() {
				mockRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			result, err := service.ListUsers(context.Background(), tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	service := NewUserService(mockRepo, mockIDGen)

	testUser := createTestUser()

	tests := []struct {
		name          string
		userID        string
		mockBehavior  func()
		expectedError string
	}{
		{
			name:   "successful deletion",
			userID: "test-id-123",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					Delete(gomock.Any(), "test-id-123").
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent-id",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "nonexistent-id").
					Return(nil, nil).
					Times(1)
			},
			expectedError: "not found",
		},
		{
			name:   "empty user ID",
			userID: "",
			mockBehavior: func() {
				// No mock calls expected
			},
			expectedError: "id is required",
		},
		{
			name:   "repository delete error",
			userID: "test-id-123",
			mockBehavior: func() {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), "test-id-123").
					Return(testUser, nil).
					Times(1)
				mockRepo.EXPECT().
					Delete(gomock.Any(), "test-id-123").
					Return(errors.New("database error")).
					Times(1)
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			err := service.DeleteUser(context.Background(), tt.userID)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
