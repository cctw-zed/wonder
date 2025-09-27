package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/internal/infrastructure/repository"
	"github.com/cctw-zed/wonder/pkg/jwt"
	"github.com/cctw-zed/wonder/pkg/logger"
	idMocks "github.com/cctw-zed/wonder/pkg/snowflake/id/mocks"
)

// TestAuthServiceIntegration verifies authentication service integration
// with user service, token service, and repository layers
func TestAuthServiceIntegration(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock ID generator with unique IDs for actual usage
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	gomock.InOrder(
		mockIDGen.EXPECT().Generate().Return("auth-test-1").Times(1),
		mockIDGen.EXPECT().Generate().Return("auth-test-2").Times(1),
		mockIDGen.EXPECT().Generate().Return("auth-test-3").Times(1),
		mockIDGen.EXPECT().Generate().Return("auth-test-4").Times(1),
		mockIDGen.EXPECT().Generate().Return("auth-test-5").Times(1),
	)
	// Allow additional calls with fallback IDs
	mockIDGen.EXPECT().Generate().AnyTimes().DoAndReturn(func() string {
		return fmt.Sprintf("auth-test-fallback-%d", time.Now().UnixNano())
	})

	// Skip if no test database
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}

	// Setup repository and services
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo, mockIDGen)

	// Setup JWT token service with test configuration
	tokenService := jwt.NewTokenService("test-secret-key-for-integration-tests", time.Hour)
	authService := service.NewAuthService(userService, tokenService)

	ctx := context.Background()

	t.Run("Complete authentication flow integration", func(t *testing.T) {
		// Test data
		testEmail := "auth.integration@test.com"
		testPassword := "integration123"
		testName := "Auth Integration User"

		// Step 1: Create user through user service
		createdUser, err := userService.Register(ctx, testEmail, testName, testPassword)
		require.NoError(t, err)
		require.NotNil(t, createdUser)

		// Step 2: Login through auth service
		loginResp, err := authService.Login(ctx, testEmail, testPassword)
		require.NoError(t, err)
		require.NotNil(t, loginResp)
		require.NotEmpty(t, loginResp.AccessToken)
		assert.Equal(t, createdUser.ID, loginResp.User.ID)
		assert.Equal(t, testEmail, loginResp.User.Email)
		assert.Equal(t, testName, loginResp.User.Name)

		// Step 3: Validate token through auth service
		claims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, createdUser.ID, claims.UserID)

		// Step 4: Logout through auth service
		err = authService.Logout(ctx, loginResp.AccessToken)
		require.NoError(t, err)

		// Step 5: Verify logout behavior (token still valid due to current implementation)
		// NOTE: Current implementation doesn't actually invalidate tokens on logout
		// (see TODO comment in auth_service.go line 102-103)
		// Token remains valid until natural expiration
		claims, err = authService.ValidateToken(ctx, loginResp.AccessToken)
		require.NoError(t, err) // Token is still valid (current implementation)
		assert.Equal(t, createdUser.ID, claims.UserID)
	})

	t.Run("Authentication error scenarios integration", func(t *testing.T) {
		t.Run("login with non-existent user", func(t *testing.T) {
			loginResp, err := authService.Login(ctx, "nonexistent@test.com", "password123")
			assert.Error(t, err)
			assert.Nil(t, loginResp)
		})

		t.Run("login with wrong password", func(t *testing.T) {
			// Create user first
			testEmail := "wrong.password@test.com"
			_, err := userService.Register(ctx, testEmail, "Wrong Password User", "correct123")
			require.NoError(t, err)

			// Try login with wrong password
			loginResp, err := authService.Login(ctx, testEmail, "wrong123")
			assert.Error(t, err)
			assert.Nil(t, loginResp)
		})

		t.Run("validate invalid token", func(t *testing.T) {
			claims, err := authService.ValidateToken(ctx, "invalid.jwt.token")
			assert.Error(t, err)
			assert.Nil(t, claims)
		})

		t.Run("validate expired token", func(t *testing.T) {
			// Create token service with very short expiry
			shortTokenService := jwt.NewTokenService("test-secret-key", time.Millisecond)
			shortAuthService := service.NewAuthService(userService, shortTokenService)

			// Create user and login
			testEmail := "expired.token@test.com"
			_, err := userService.Register(ctx, testEmail, "Expired Token User", "password123")
			require.NoError(t, err)

			loginResp, err := shortAuthService.Login(ctx, testEmail, "password123")
			require.NoError(t, err)

			// Wait for token to expire
			time.Sleep(10 * time.Millisecond)

			// Try to validate expired token
			claims, err := shortAuthService.ValidateToken(ctx, loginResp.AccessToken)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	})

	t.Run("Token service integration", func(t *testing.T) {
		// Create user for token testing
		testEmail := "token.integration@test.com"
		createdUser, err := userService.Register(ctx, testEmail, "Token Integration User", "password123")
		require.NoError(t, err)

		// Login to get token
		loginResp, err := authService.Login(ctx, testEmail, "password123")
		require.NoError(t, err)

		// Test token validation with different scenarios
		t.Run("valid token returns correct claims", func(t *testing.T) {
			claims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, createdUser.ID, claims.UserID)
			assert.True(t, claims.ExpiresAt.After(time.Now()))
		})

		t.Run("token contains correct user information", func(t *testing.T) {
			claims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, createdUser.ID, claims.UserID)
		})
	})

	t.Run("Concurrent authentication operations", func(t *testing.T) {
		// Create user for concurrent testing
		testEmail := "concurrent.auth@test.com"
		_, err := userService.Register(ctx, testEmail, "Concurrent Auth User", "password123")
		require.NoError(t, err)

		// Test concurrent logins
		numConcurrent := 10
		results := make(chan error, numConcurrent)

		for i := 0; i < numConcurrent; i++ {
			go func() {
				loginResp, err := authService.Login(ctx, testEmail, "password123")
				if err != nil {
					results <- err
					return
				}

				// Validate token
				_, err = authService.ValidateToken(ctx, loginResp.AccessToken)
				if err != nil {
					results <- err
					return
				}

				// Logout
				err = authService.Logout(ctx, loginResp.AccessToken)
				results <- err
			}()
		}

		// Collect results
		for i := 0; i < numConcurrent; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent operation %d should succeed", i)
		}
	})
}

// TestAuthMiddlewareIntegration tests middleware integration with auth service
func TestAuthMiddlewareIntegration(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock ID generator with unique IDs for middleware test
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	gomock.InOrder(
		mockIDGen.EXPECT().Generate().Return("middleware-test-1").Times(1),
		mockIDGen.EXPECT().Generate().Return("middleware-test-2").Times(1),
	)
	// Allow additional calls with fallback IDs
	mockIDGen.EXPECT().Generate().AnyTimes().DoAndReturn(func() string {
		return fmt.Sprintf("middleware-fallback-%d", time.Now().UnixNano())
	})

	// Skip if no test database
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}

	// Setup services
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo, mockIDGen)
	tokenService := jwt.NewTokenService("test-secret-middleware", time.Hour)
	authService := service.NewAuthService(userService, tokenService)

	ctx := context.Background()

	t.Run("Middleware validates tokens from auth service", func(t *testing.T) {
		// Create user and login to get valid token
		testEmail := "middleware.integration@test.com"
		createdUser, err := userService.Register(ctx, testEmail, "Middleware Integration User", "password123")
		require.NoError(t, err)

		loginResp, err := authService.Login(ctx, testEmail, "password123")
		require.NoError(t, err)

		// Test token validation through auth service (simulating middleware)
		claims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, claims.UserID)

		// Test logout functionality
		err = authService.Logout(ctx, loginResp.AccessToken)
		require.NoError(t, err)

		// NOTE: Current implementation doesn't actually invalidate tokens on logout
		// (see TODO comment in auth_service.go line 102-103)
		// Token remains valid until natural expiration
		// This test validates the current behavior - when token blacklist is implemented,
		// this should be changed to expect an error
		claims, err = authService.ValidateToken(ctx, loginResp.AccessToken)
		require.NoError(t, err) // Token is still valid (current implementation)
		assert.Equal(t, createdUser.ID, claims.UserID)
	})

	t.Run("Integration with user service operations", func(t *testing.T) {
		// Create user
		testEmail := "user.ops@test.com"
		createdUser, err := userService.Register(ctx, testEmail, "User Ops User", "password123")
		require.NoError(t, err)

		// Login
		loginResp, err := authService.Login(ctx, testEmail, "password123")
		require.NoError(t, err)

		// Use token to validate user access (simulating protected endpoint)
		claims, err := authService.ValidateToken(ctx, loginResp.AccessToken)
		require.NoError(t, err)

		// Use claims to get user info (simulating middleware injecting user context)
		user, err := userService.GetProfile(ctx, claims.UserID)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, testEmail, user.Email)
	})
}
