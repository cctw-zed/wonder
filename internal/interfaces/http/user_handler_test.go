package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/domain/user/mocks"
	"github.com/cctw-zed/wonder/internal/testutil/builder"
	apperrors "github.com/cctw-zed/wonder/pkg/errors"
)

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	// Setup expected behavior
	expectedUser := builder.NewUserBuilderForTesting().
		ValidUserWithEmail("test@example.com")

	mockUserService.EXPECT().
		Register(gomock.Any(), "test@example.com", "Test User", "password123").
		Return(expectedUser, nil).
		Times(1)

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup Gin router
	router := setupGinTest()
	router.POST("/users/register", handler.Register)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	responseUser := response["user"].(map[string]interface{})
	assert.Equal(t, expectedUser.ID, responseUser["id"])
	assert.Equal(t, expectedUser.Email, responseUser["email"])
	assert.Equal(t, expectedUser.Name, responseUser["name"])
}

func TestUserHandler_Register_ValidationErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"name": "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request data",
		},
		{
			name: "missing name",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request data",
		},
		{
			name: "invalid email format",
			requestBody: RegisterRequest{
				Email:    "invalid-email",
				Name:     "Test User",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request data",
		},
		{
			name: "name too short",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Name:     "A",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request data",
		},
		{
			name: "name too long",
			requestBody: RegisterRequest{
				Email:    "test@example.com",
				Name:     "This is a very long name that exceeds the maximum length allowed for user names in the system",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request data",
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin router
			router := setupGinTest()
			router.POST("/users/register", handler.Register)

			// Prepare request body
			var jsonBody []byte
			if str, ok := tt.requestBody.(string); ok {
				jsonBody = []byte(str)
			} else {
				jsonBody, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.errorContains != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Check various possible error fields
				var errorMsg string
				var exists bool
				if response["error"] != nil {
					errorMsg = response["error"].(string)
					exists = true
				} else if response["message"] != nil {
					errorMsg = response["message"].(string)
					exists = true
				} else if response["details"] != nil {
					if details, ok := response["details"].(map[string]interface{}); ok {
						if validation_error, ok := details["validation_error"].(string); ok {
							errorMsg = validation_error
							exists = true
						}
					}
				}
				require.True(t, exists, "Expected error message in response: %v", response)
				assert.Contains(t, errorMsg, tt.errorContains)
			}
		})
	}
}

func TestUserHandler_Register_ServiceErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	tests := []struct {
		name           string
		serviceError   error
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "user already exists",
			serviceError:   errors.New("email already exists"),
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "An internal server error occurred",
		},
		{
			name:           "database error",
			serviceError:   errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "An internal server error occurred",
		},
		{
			name:           "validation error",
			serviceError:   errors.New("invalid email: email is required"),
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "An internal server error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup expected behavior
			mockUserService.EXPECT().
				Register(gomock.Any(), "test@example.com", "Test User", "password123").
				Return(nil, tt.serviceError).
				Times(1)

			// Setup HTTP request
			requestBody := RegisterRequest{
				Email:    "test@example.com",
				Name:     "Test User",
				Password: "password123",
			}
			jsonBody, _ := json.Marshal(requestBody)

			// Setup Gin router
			router := setupGinTest()
			router.POST("/users/register", handler.Register)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check various possible error fields
			var errorMsg string
			var exists bool
			if response["error"] != nil {
				errorMsg = response["error"].(string)
				exists = true
			} else if response["message"] != nil {
				errorMsg = response["message"].(string)
				exists = true
			} else if response["details"] != nil {
				if details, ok := response["details"].(map[string]interface{}); ok {
					if validation_error, ok := details["validation_error"].(string); ok {
						errorMsg = validation_error
						exists = true
					}
				}
			}
			require.True(t, exists, "Expected error message in response: %v", response)
			assert.Contains(t, errorMsg, tt.errorContains)
		})
	}
}

func TestUserHandler_Register_ContextPropagation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	// Setup expected behavior with context validation
	expectedUser := builder.NewUserBuilderForTesting().
		ValidUserWithEmail("test@example.com")

	mockUserService.EXPECT().
		Register(gomock.Any(), "test@example.com", "Test User", "password123").
		DoAndReturn(func(ctx context.Context, email, name, password string) (*user.User, error) {
			// Verify context is properly passed
			assert.NotNil(t, ctx)
			return expectedUser, nil
		}).
		Times(1)

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup Gin router
	router := setupGinTest()
	router.POST("/users/register", handler.Register)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockUserService, handler.userService)
}

// Benchmark test for the register endpoint
func BenchmarkUserHandler_Register(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	expectedUser := builder.NewUserBuilderForTesting().
		ValidUserWithEmail("bench@example.com")

	// Setup mock expectations for all iterations
	mockUserService.EXPECT().
		Register(gomock.Any(), "bench@example.com", "Bench User", "benchpass123").
		Return(expectedUser, nil).
		AnyTimes()

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email:    "bench@example.com",
		Name:     "Bench User",
		Password: "benchpass123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Setup Gin router
	router := setupGinTest()
	router.POST("/users/register", handler.Register)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	}
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	// Setup expected behavior
	expectedUser := builder.NewUserBuilderForTesting().
		ValidUserWithEmail("test@example.com")

	mockUserService.EXPECT().
		GetProfile(gomock.Any(), "test-user-id").
		Return(expectedUser, nil).
		Times(1)

	// Setup Gin router
	router := setupGinTest()
	router.GET("/users/:id", handler.GetProfile)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/test-user-id", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "trace_id")

	userData := response["user"].(map[string]interface{})
	assert.Equal(t, expectedUser.ID, userData["id"])
	assert.Equal(t, expectedUser.Email, userData["email"])
}

func TestUserHandler_GetProfile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	mockUserService.EXPECT().
		GetProfile(gomock.Any(), "nonexistent-id").
		Return(nil, apperrors.NewEntityNotFoundError("user", "nonexistent-id")).
		Times(1)

	router := setupGinTest()
	router.GET("/users/:id", handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	updatedUser := builder.NewUserBuilderForTesting().
		ValidUserWithEmail("updated@example.com")

	mockUserService.EXPECT().
		UpdateProfile(gomock.Any(), "test-user-id", gomock.Any()).
		Return(updatedUser, nil).
		Times(1)

	requestBody := user.UpdateProfileRequest{
		Email: "updated@example.com",
		Name:  "Updated Name",
	}
	jsonBody, _ := json.Marshal(requestBody)

	router := setupGinTest()
	router.PUT("/users/:id", handler.UpdateProfile)

	req := httptest.NewRequest(http.MethodPut, "/users/test-user-id", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "trace_id")
}

func TestUserHandler_UpdateProfile_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	router := setupGinTest()
	router.PUT("/users/:id", handler.UpdateProfile)

	req := httptest.NewRequest(http.MethodPut, "/users/test-user-id", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	testUsers := []*user.User{
		builder.NewUserBuilderForTesting().ValidUserWithEmail("user1@example.com"),
		builder.NewUserBuilderForTesting().ValidUserWithEmail("user2@example.com"),
	}

	expectedResponse := &user.ListUsersResponse{
		Users:      testUsers,
		Total:      2,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	mockUserService.EXPECT().
		ListUsers(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil).
		Times(1)

	router := setupGinTest()
	router.GET("/users", handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "trace_id")

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(1), data["page"])
}

func TestUserHandler_ListUsers_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	expectedResponse := &user.ListUsersResponse{
		Users:      []*user.User{},
		Total:      0,
		Page:       1,
		PageSize:   10,
		TotalPages: 0,
	}

	mockUserService.EXPECT().
		ListUsers(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil).
		Times(1)

	router := setupGinTest()
	router.GET("/users", handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&page_size=5&email=test&name=john", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	mockUserService.EXPECT().
		DeleteUser(gomock.Any(), "test-user-id").
		Return(nil).
		Times(1)

	router := setupGinTest()
	router.DELETE("/users/:id", handler.DeleteUser)

	req := httptest.NewRequest(http.MethodDelete, "/users/test-user-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Contains(t, response, "trace_id")
	assert.Equal(t, "User deleted successfully", response["message"])
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	mockUserService.EXPECT().
		DeleteUser(gomock.Any(), "nonexistent-id").
		Return(apperrors.NewEntityNotFoundError("user", "nonexistent-id")).
		Times(1)

	router := setupGinTest()
	router.DELETE("/users/:id", handler.DeleteUser)

	req := httptest.NewRequest(http.MethodDelete, "/users/nonexistent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_GetProfile_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	router := setupGinTest()
	router.GET("/users/:id", handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// The request should result in a 404 because the route doesn't match
	assert.Equal(t, http.StatusNotFound, w.Code)
}
