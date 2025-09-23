package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/pkg/errors"
)

type UserHandler struct {
	userService  user.UserService
	errorMapper  *errors.ErrorMapper
	errorLogger  errors.ErrorLogger
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		errorMapper: errors.NewErrorMapper(),
		errorLogger: errors.NewDefaultErrorLogger("user-service"),
	}
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=2,max=50"`
}

func (h *UserHandler) Register(c *gin.Context) {
	// Generate trace ID for request tracking
	traceID := h.generateTraceID()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle validation errors from Gin binding
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

	// Call application service
	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		// Log the error with structured logging
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "user_registration",
			"email":     req.Email,
			"name":      req.Name,
		})

		// Map service layer error to HTTP error
		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	// Success response
	c.JSON(http.StatusCreated, map[string]interface{}{
		"user":     user,
		"trace_id": traceID,
	})
}

// generateTraceID generates a unique trace ID for request tracking
func (h *UserHandler) generateTraceID() string {
	return uuid.New().String()
}
