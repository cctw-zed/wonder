package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/middleware"
	"github.com/cctw-zed/wonder/pkg/errors"
	"strconv"
)

type UserHandler struct {
	userService user.UserService
	errorMapper *errors.ErrorMapper
	errorLogger errors.ErrorLogger
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		errorMapper: errors.NewErrorMapper(),
		errorLogger: errors.NewDefaultErrorLogger("user-service"),
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
	// Get trace ID from context (injected by middleware)
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

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
	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Name, req.Password)
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

// GetProfile retrieves user profile by ID
func (h *UserHandler) GetProfile(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())
	userID := c.Param("id")

	if userID == "" {
		httpErr := errors.NewHTTPError(
			http.StatusBadRequest,
			errors.CodeValidationError,
			"User ID is required",
			map[string]interface{}{"field": "id"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "get_user_profile",
			"user_id":   userID,
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":     user,
		"trace_id": traceID,
	})
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())
	userID := c.Param("id")

	if userID == "" {
		httpErr := errors.NewHTTPError(
			http.StatusBadRequest,
			errors.CodeValidationError,
			"User ID is required",
			map[string]interface{}{"field": "id"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	var req user.UpdateProfileRequest
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

	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "update_user_profile",
			"user_id":   userID,
			"request":   req,
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":     updatedUser,
		"trace_id": traceID,
	})
}

// ListUsers retrieves users with pagination and filtering
func (h *UserHandler) ListUsers(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	email := c.Query("email")
	name := c.Query("name")

	req := &user.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Email:    email,
		Name:     name,
	}

	response, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "list_users",
			"request":   req,
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data":     response,
		"trace_id": traceID,
	})
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	traceID := middleware.GetTraceIDFromContext(c.Request.Context())
	userID := c.Param("id")

	if userID == "" {
		httpErr := errors.NewHTTPError(
			http.StatusBadRequest,
			errors.CodeValidationError,
			"User ID is required",
			map[string]interface{}{"field": "id"},
			traceID,
		)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.errorLogger.LogError(c.Request.Context(), err, traceID, map[string]interface{}{
			"operation": "delete_user",
			"user_id":   userID,
		})

		httpErr := h.errorMapper.MapToHTTPError(err, traceID)
		c.JSON(httpErr.StatusCode, httpErr)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "User deleted successfully",
		"trace_id": traceID,
	})
}
