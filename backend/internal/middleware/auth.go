package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/jwt"
)

const (
	// AuthorizationHeader is the HTTP header name for authorization
	AuthorizationHeader = "Authorization"
	// UserIDKey is the context key for storing user ID
	UserIDKey = "user_id"
	// UserIDHeader is the HTTP header name for user ID (injected into request)
	UserIDHeader = "X-User-ID"
	// BearerPrefix is the prefix for Bearer tokens
	BearerPrefix = "Bearer "
)

// AuthMiddleware provides JWT authentication functionality
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	if authService == nil {
		panic("auth service cannot be nil")
	}
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth creates middleware that requires valid JWT authentication
// Returns 401 Unauthorized if token is missing, invalid, or expired
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.validateTokenFromRequest(c)
		if err != nil {
			m.handleAuthError(c, err)
			c.Abort()
			return
		}

		// Inject user ID into context and request headers
		m.injectUserContext(c, claims.UserID)
		c.Next()
	}
}

// OptionalAuth creates middleware that optionally validates JWT authentication
// Continues processing even if no token is provided, but validates if present
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := m.validateTokenFromRequest(c)
		if err == nil && claims != nil {
			// Token is valid, inject user context
			m.injectUserContext(c, claims.UserID)
		}
		// Continue processing regardless of token validity
		c.Next()
	}
}

// validateTokenFromRequest extracts and validates JWT token from request
func (m *AuthMiddleware) validateTokenFromRequest(c *gin.Context) (*jwt.Claims, error) {
	// Extract token from Authorization header
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return nil, errors.NewUnauthorizedError(
			"auth_middleware",
			"missing_token",
			"Authorization header is required",
		)
	}

	// Check Bearer prefix
	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return nil, errors.NewUnauthorizedError(
			"auth_middleware",
			"invalid_token_format",
			"Authorization header must use Bearer token format",
		)
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return nil, errors.NewUnauthorizedError(
			"auth_middleware",
			"empty_token",
			"Bearer token cannot be empty",
		)
	}

	// Validate token using auth service
	claims, err := m.authService.ValidateToken(c.Request.Context(), token)
	if err != nil {
		return nil, err // Auth service already returns proper error types
	}

	return claims, nil
}

// injectUserContext injects user ID into both request context and headers
func (m *AuthMiddleware) injectUserContext(c *gin.Context, userID string) {
	// Inject user ID into request context
	ctx := context.WithValue(c.Request.Context(), UserIDKey, userID)
	c.Request = c.Request.WithContext(ctx)

	// Inject user ID into request headers for easy access in handlers
	c.Header(UserIDHeader, userID)
}

// handleAuthError handles authentication errors with proper HTTP responses
func (m *AuthMiddleware) handleAuthError(c *gin.Context, err error) {
	traceID := GetTraceIDFromContext(c.Request.Context())

	// Map application/domain errors to HTTP errors
	errorMapper := errors.NewErrorMapper()
	httpErr := errorMapper.MapToHTTPError(err, traceID)

	c.JSON(httpErr.StatusCode, httpErr)
}

// GetUserIDFromContext extracts user ID from context
// Returns empty string if no user ID is found in context
func GetUserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		if str, ok := userID.(string); ok {
			return str
		}
	}

	return ""
}

// GetUserIDFromGinContext extracts user ID from Gin context
// This is a convenience function for Gin handlers
func GetUserIDFromGinContext(c *gin.Context) string {
	return GetUserIDFromContext(c.Request.Context())
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(ctx context.Context) bool {
	return GetUserIDFromContext(ctx) != ""
}

// IsAuthenticatedGin checks if the current Gin request is authenticated
func IsAuthenticatedGin(c *gin.Context) bool {
	return GetUserIDFromGinContext(c) != ""
}
