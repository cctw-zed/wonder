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
// Note: This endpoint is protected by auth middleware, so token is already validated
func (h *AuthHandler) Logout(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	// Extract token from Authorization header (already validated by middleware)
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

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
// Note: This endpoint is protected by auth middleware, so user ID is already available in context
func (h *AuthHandler) GetMe(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	// Get user ID from context (injected by auth middleware)
	userID := middleware.GetUserIDFromGinContext(c)

	// Return user ID from middleware context
	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":  userID,
		"trace_id": traceID,
	})
}
