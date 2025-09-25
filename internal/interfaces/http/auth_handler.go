package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/internal/middleware"
	"github.com/cctw-zed/wonder/pkg/errors"
)

type AuthHandler struct {
	authService service.AuthService
	errorMapper *errors.ErrorMapper
	errorLogger errors.ErrorLogger
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		errorMapper: errors.NewErrorMapper(),
		errorLogger: errors.NewDefaultErrorLogger("auth-service"),
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login authenticates user and returns JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpErr := errors.NewHTTPError(
			http.StatusBadRequest,
			errors.CodeValidationError,
			"Invalid request data",
			map[string]interface{}{"validation_error": err.Error()},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Authenticate user
	response, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "user_login",
			"email":     req.Email,
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Success response
	c.JSON(http.StatusOK, map[string]interface{}{
		"data":     response,
		"trace_id": traceID,
	})
}

// Logout invalidates the current user's token
func (h *AuthHandler) Logout(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		httpErr := errors.NewHTTPError(
			http.StatusUnauthorized,
			errors.CodeUnauthorized,
			"Authorization header required",
			map[string]interface{}{"header": "Authorization"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		httpErr := errors.NewHTTPError(
			http.StatusUnauthorized,
			errors.CodeUnauthorized,
			"Invalid authorization header format",
			map[string]interface{}{"expected": "Bearer <token>"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	token := tokenParts[1]

	// Logout user
	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "user_logout",
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Success response
	c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Logout successful",
		"trace_id": traceID,
	})
}

// GetMe returns current user information based on JWT token
func (h *AuthHandler) GetMe(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		httpErr := errors.NewHTTPError(
			http.StatusUnauthorized,
			errors.CodeUnauthorized,
			"Authorization header required",
			map[string]interface{}{"header": "Authorization"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Parse Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		httpErr := errors.NewHTTPError(
			http.StatusUnauthorized,
			errors.CodeUnauthorized,
			"Invalid authorization header format",
			map[string]interface{}{"expected": "Bearer <token>"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	token := tokenParts[1]

	// Validate token and get user ID
	claims, err := h.authService.ValidateToken(c.Request.Context(), token)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "token_validation",
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Return user ID from token claims
	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":  claims.UserID,
		"trace_id": traceID,
	})
}