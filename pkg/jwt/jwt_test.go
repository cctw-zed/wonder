package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenService(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	expiry := 24 * time.Hour

	service := NewTokenService(signingKey, expiry)
	require.NotNil(t, service)

	jwtService, ok := service.(*JWTService)
	require.True(t, ok)
	assert.Equal(t, []byte(signingKey), jwtService.signingKey)
	assert.Equal(t, expiry, jwtService.expiry)
}

func TestJWTService_GenerateToken(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	expiry := 24 * time.Hour
	service := NewTokenService(signingKey, expiry)

	tests := []struct {
		name      string
		userID    string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "successful token generation",
			userID:  "user123",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: true,
			errMsg:  "user_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateToken(tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)

				// Token should be in JWT format (three parts separated by dots)
				parts := len(token) > 0
				assert.True(t, parts, "Token should not be empty")
			}
		})
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	expiry := 24 * time.Hour
	service := NewTokenService(signingKey, expiry)

	// Generate a valid token first
	userID := "user123"
	validToken, err := service.GenerateToken(userID)
	require.NoError(t, err)
	require.NotEmpty(t, validToken)

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		errMsg    string
		checkUserID bool
	}{
		{
			name:        "valid token",
			token:       validToken,
			wantErr:     false,
			checkUserID: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errMsg:  "token is required",
		},
		{
			name:    "invalid token format",
			token:   "invalid.token",
			wantErr: true,
			errMsg:  "invalid token",
		},
		{
			name:    "malformed JWT",
			token:   "not.a.jwt.token",
			wantErr: true,
			errMsg:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)

				if tt.checkUserID {
					assert.Equal(t, userID, claims.UserID)
					assert.Equal(t, "wonder-api", claims.Issuer)
					assert.Equal(t, userID, claims.Subject)
					assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
					assert.True(t, claims.IssuedAt.Time.Before(time.Now().Add(time.Second)))
				}
			}
		})
	}
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	// Create a service with very short expiry
	shortExpiry := 1 * time.Millisecond
	service := NewTokenService(signingKey, shortExpiry)

	userID := "user123"
	token, err := service.GenerateToken(userID)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	claims, err := service.ValidateToken(token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	assert.Nil(t, claims)
}

func TestJWTService_ValidateToken_WrongSigningKey(t *testing.T) {
	signingKey1 := "test-signing-key-32-chars-minimum-1"
	signingKey2 := "test-signing-key-32-chars-minimum-2"

	service1 := NewTokenService(signingKey1, 24*time.Hour)
	service2 := NewTokenService(signingKey2, 24*time.Hour)

	userID := "user123"
	token, err := service1.GenerateToken(userID)
	require.NoError(t, err)

	// Try to validate with different signing key
	claims, err := service2.ValidateToken(token)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	assert.Nil(t, claims)
}

func TestJWTService_GetSigningKey(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	service := NewTokenService(signingKey, 24*time.Hour)

	retrievedKey := service.GetSigningKey()
	assert.Equal(t, []byte(signingKey), retrievedKey)
}

func TestJWTService_TokenLifecycle(t *testing.T) {
	signingKey := "test-signing-key-32-chars-minimum"
	expiry := 1 * time.Hour
	service := NewTokenService(signingKey, expiry)

	userID := "user123"

	// Generate token
	token, err := service.GenerateToken(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Validate immediately
	claims, err := service.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)

	// Check token expiry is set correctly
	expectedExpiry := time.Now().Add(expiry)
	actualExpiry := claims.ExpiresAt.Time

	// Allow for small time differences (within 10 seconds)
	timeDiff := actualExpiry.Sub(expectedExpiry)
	assert.True(t, timeDiff > -10*time.Second && timeDiff < 10*time.Second,
		"Token expiry should be close to expected time")
}