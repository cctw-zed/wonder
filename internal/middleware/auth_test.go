package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	serviceMocks "github.com/cctw-zed/wonder/internal/application/service/mocks"
	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/jwt"
)

func setupAuthMiddlewareTest(t *testing.T) (*AuthMiddleware, *serviceMocks.MockAuthService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuthService := serviceMocks.NewMockAuthService(ctrl)
	middleware := NewAuthMiddleware(mockAuthService)
	return middleware, mockAuthService, ctrl
}

func createTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TraceIDMiddleware()) // Add trace ID middleware

	// Protected route with middleware
	router.GET("/protected", middleware, func(c *gin.Context) {
		userID := GetUserIDFromGinContext(c)
		c.JSON(http.StatusOK, gin.H{
			"message":       "success",
			"user_id":       userID,
			"authenticated": IsAuthenticatedGin(c),
		})
	})

	return router
}

func TestNewAuthMiddleware(t *testing.T) {
	t.Run("should create middleware with valid auth service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAuthService := serviceMocks.NewMockAuthService(ctrl)
		middleware := NewAuthMiddleware(mockAuthService)

		assert.NotNil(t, middleware)
		assert.Equal(t, mockAuthService, middleware.authService)
	})

	t.Run("should panic with nil auth service", func(t *testing.T) {
		assert.Panics(t, func() {
			NewAuthMiddleware(nil)
		})
	})
}

func TestRequireAuth_Success(t *testing.T) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Setup mock expectations
	expectedClaims := &jwt.Claims{
		UserID: "user123",
	}
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedClaims, nil).
		Times(1)

	// Create test router
	router := createTestRouter(middleware.RequireAuth())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["message"])
	assert.Equal(t, "user123", response["user_id"])
	assert.True(t, response["authenticated"].(bool))
	assert.Equal(t, "user123", w.Header().Get(UserIDHeader))
}

func TestRequireAuth_MissingToken(t *testing.T) {
	middleware, _, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Create test router
	router := createTestRouter(middleware.RequireAuth())

	// Create test request without Authorization header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response errors.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "UNAUTHORIZED", string(response.Code()))
	assert.Equal(t, "Unauthorized access", response.Message)
}

func TestRequireAuth_InvalidTokenFormat(t *testing.T) {
	middleware, _, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	testCases := []struct {
		name            string
		authHeader      string
		expectedMessage string
	}{
		{
			name:            "missing Bearer prefix",
			authHeader:      "invalid-token",
			expectedMessage: "Authorization header must use Bearer token format",
		},
		{
			name:            "empty token after Bearer",
			authHeader:      "Bearer ",
			expectedMessage: "Bearer token cannot be empty",
		},
		{
			name:            "only Bearer keyword",
			authHeader:      "Bearer",
			expectedMessage: "Bearer token cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test router
			router := createTestRouter(middleware.RequireAuth())

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.authHeader)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var response errors.HTTPError
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, "UNAUTHORIZED", string(response.Code()))
			assert.Equal(t, "Unauthorized access", response.Message)
		})
	}
}

func TestRequireAuth_TokenValidationFailure(t *testing.T) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Setup mock expectations for token validation failure
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token").
		Return(nil, errors.NewUnauthorizedError("token_validation", "", "token expired")).
		Times(1)

	// Create test router
	router := createTestRouter(middleware.RequireAuth())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response errors.HTTPError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", string(response.Code()))
	assert.Equal(t, "Unauthorized access", response.Message)
}

func TestOptionalAuth_WithValidToken(t *testing.T) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Setup mock expectations
	expectedClaims := &jwt.Claims{
		UserID: "user456",
	}
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedClaims, nil).
		Times(1)

	// Create test router
	router := createTestRouter(middleware.OptionalAuth())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["message"])
	assert.Equal(t, "user456", response["user_id"])
	assert.True(t, response["authenticated"].(bool))
	assert.Equal(t, "user456", w.Header().Get(UserIDHeader))
}

func TestOptionalAuth_WithoutToken(t *testing.T) {
	middleware, _, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Create test router
	router := createTestRouter(middleware.OptionalAuth())

	// Create test request without Authorization header
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions - should continue processing without authentication
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["message"])
	assert.Equal(t, "", response["user_id"])
	assert.False(t, response["authenticated"].(bool))
	assert.Equal(t, "", w.Header().Get(UserIDHeader))
}

func TestOptionalAuth_WithInvalidToken(t *testing.T) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Setup mock expectations for token validation failure
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token").
		Return(nil, errors.NewUnauthorizedError("token_validation", "", "token expired")).
		Times(1)

	// Create test router
	router := createTestRouter(middleware.OptionalAuth())

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions - should continue processing without authentication
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["message"])
	assert.Equal(t, "", response["user_id"])
	assert.False(t, response["authenticated"].(bool))
	assert.Equal(t, "", w.Header().Get(UserIDHeader))
}

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("should return user ID from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "user789")
		userID := GetUserIDFromContext(ctx)
		assert.Equal(t, "user789", userID)
	})

	t.Run("should return empty string for nil context", func(t *testing.T) {
		userID := GetUserIDFromContext(nil)
		assert.Equal(t, "", userID)
	})

	t.Run("should return empty string when no user ID in context", func(t *testing.T) {
		ctx := context.Background()
		userID := GetUserIDFromContext(ctx)
		assert.Equal(t, "", userID)
	})

	t.Run("should return empty string when user ID is wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, 123)
		userID := GetUserIDFromContext(ctx)
		assert.Equal(t, "", userID)
	})
}

func TestIsAuthenticated(t *testing.T) {
	t.Run("should return true when user ID exists", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "user123")
		assert.True(t, IsAuthenticated(ctx))
	})

	t.Run("should return false when no user ID", func(t *testing.T) {
		ctx := context.Background()
		assert.False(t, IsAuthenticated(ctx))
	})

	t.Run("should return false for nil context", func(t *testing.T) {
		assert.False(t, IsAuthenticated(nil))
	})
}

func TestGetUserIDFromGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return user ID from Gin context", func(t *testing.T) {
		// Create a test context with user ID
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserIDKey, "gin-user123")
		req = req.WithContext(ctx)
		c.Request = req

		userID := GetUserIDFromGinContext(c)
		assert.Equal(t, "gin-user123", userID)
	})

	t.Run("should return empty string when no user ID in Gin context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req

		userID := GetUserIDFromGinContext(c)
		assert.Equal(t, "", userID)
	})
}

func TestIsAuthenticatedGin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return true when user ID exists in Gin context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserIDKey, "gin-user456")
		req = req.WithContext(ctx)
		c.Request = req

		assert.True(t, IsAuthenticatedGin(c))
	})

	t.Run("should return false when no user ID in Gin context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		c.Request = req

		assert.False(t, IsAuthenticatedGin(c))
	})
}

func TestMiddleware_TraceIDIntegration(t *testing.T) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(t)
	defer ctrl.Finish()

	// Setup mock expectations
	expectedClaims := &jwt.Claims{
		UserID: "trace-user",
	}
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedClaims, nil).
		Times(1)

	// Create test router with trace ID middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TraceIDMiddleware()) // This should run first
	router.Use(middleware.RequireAuth())

	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceIDFromContext(c.Request.Context())
		userID := GetUserIDFromGinContext(c)

		c.JSON(http.StatusOK, gin.H{
			"trace_id": traceID,
			"user_id":  userID,
		})
	})

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify trace ID and user ID are both present
	assert.NotEmpty(t, response["trace_id"])
	assert.Equal(t, "trace-user", response["user_id"])

	// Verify trace ID is also in response header
	assert.Equal(t, response["trace_id"], w.Header().Get(TraceIDHeader))
	assert.Equal(t, "trace-user", w.Header().Get(UserIDHeader))
}

func BenchmarkRequireAuth_ValidToken(b *testing.B) {
	middleware, mockAuthService, ctrl := setupAuthMiddlewareTest(&testing.T{})
	defer ctrl.Finish()

	// Setup mock expectations
	expectedClaims := &jwt.Claims{
		UserID: "benchmark-user",
	}
	mockAuthService.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(expectedClaims, nil).
		Times(b.N)

	// Create test router
	router := createTestRouter(middleware.RequireAuth())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
