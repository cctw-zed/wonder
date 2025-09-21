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
		Register(gomock.Any(), "test@example.com", "Test User").
		Return(expectedUser, nil).
		Times(1)

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email: "test@example.com",
		Name:  "Test User",
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

	var responseUser user.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	require.NoError(t, err)

	assert.Equal(t, expectedUser.ID, responseUser.ID)
	assert.Equal(t, expectedUser.Email, responseUser.Email)
	assert.Equal(t, expectedUser.Name, responseUser.Name)
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
			errorContains:  "Email",
		},
		{
			name: "missing name",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Name",
		},
		{
			name: "invalid email format",
			requestBody: RegisterRequest{
				Email: "invalid-email",
				Name:  "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "email",
		},
		{
			name: "name too short",
			requestBody: RegisterRequest{
				Email: "test@example.com",
				Name:  "A",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "min",
		},
		{
			name: "name too long",
			requestBody: RegisterRequest{
				Email: "test@example.com",
				Name:  "This is a very long name that exceeds the maximum length allowed for user names in the system",
			},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "max",
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

				errorMsg, exists := response["error"]
				require.True(t, exists)
				assert.Contains(t, errorMsg.(string), tt.errorContains)
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
			errorContains:  "email already exists",
		},
		{
			name:           "database error",
			serviceError:   errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "database connection failed",
		},
		{
			name:           "validation error",
			serviceError:   errors.New("invalid email: email is required"),
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup expected behavior
			mockUserService.EXPECT().
				Register(gomock.Any(), "test@example.com", "Test User").
				Return(nil, tt.serviceError).
				Times(1)

			// Setup HTTP request
			requestBody := RegisterRequest{
				Email: "test@example.com",
				Name:  "Test User",
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

			errorMsg, exists := response["error"]
			require.True(t, exists)
			assert.Contains(t, errorMsg.(string), tt.errorContains)
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
		Register(gomock.Any(), "test@example.com", "Test User").
		DoAndReturn(func(ctx context.Context, email, name string) (*user.User, error) {
			// Verify context is properly passed
			assert.NotNil(t, ctx)
			return expectedUser, nil
		}).
		Times(1)

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email: "test@example.com",
		Name:  "Test User",
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
		Register(gomock.Any(), "bench@example.com", "Bench User").
		Return(expectedUser, nil).
		AnyTimes()

	// Setup HTTP request
	requestBody := RegisterRequest{
		Email: "bench@example.com",
		Name:  "Bench User",
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