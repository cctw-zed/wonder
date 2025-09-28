package errors_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cctw-zed/wonder/pkg/errors"
)

func TestHTTPError(t *testing.T) {
	t.Run("Create HTTP error", func(t *testing.T) {
		details := map[string]interface{}{
			"field": "email",
			"value": "invalid-email",
		}
		err := errors.NewHTTPError(http.StatusBadRequest, errors.CodeValidationError, "Invalid email format", details, "trace-123")

		assert.Equal(t, "Invalid email format", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
		assert.Equal(t, errors.CodeValidationError, err.Code())
		assert.Equal(t, "Invalid email format", err.Message)
		assert.Equal(t, details, err.Details())
		assert.Equal(t, "trace-123", err.TraceID)
	})
}

func TestErrorMapper_MapDomainErrors(t *testing.T) {
	mapper := errors.NewErrorMapper()
	traceID := "test-trace-123"

	tests := []struct {
		name           string
		domainError    errors.BaseError
		expectedStatus int
		expectedCode   errors.ErrorCode
	}{
		{
			name:           "RequiredFieldError maps to 400 Bad Request",
			domainError:    errors.NewRequiredFieldError("email", ""),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   errors.CodeRequiredField,
		},
		{
			name:           "BusinessRuleError maps to 422 Unprocessable Entity",
			domainError:    errors.NewBusinessRuleError("uniqueness", "email must be unique"),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedCode:   errors.CodeBusinessRuleError,
		},
		{
			name:           "StateTransitionError maps to 409 Conflict",
			domainError:    errors.NewStateTransitionError("User", "inactive", "active", "user is inactive"),
			expectedStatus: http.StatusConflict,
			expectedCode:   errors.CodeStateTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpErr := mapper.MapToHTTPError(tt.domainError, traceID)

			assert.Equal(t, tt.expectedStatus, httpErr.StatusCode)
			assert.Equal(t, tt.expectedCode, httpErr.Code())
			assert.Equal(t, traceID, httpErr.TraceID)
			assert.NotEmpty(t, httpErr.Message)
			assert.NotEmpty(t, httpErr.Details())
		})
	}
}

func TestErrorMapper_MapApplicationErrors(t *testing.T) {
	mapper := errors.NewErrorMapper()
	traceID := "test-trace-123"

	tests := []struct {
		name             string
		applicationError errors.BaseError
		expectedStatus   int
		expectedCode     errors.ErrorCode
	}{
		{
			name:             "EntityNotFoundError maps to 404 Not Found",
			applicationError: errors.NewEntityNotFoundError("User", "user-123"),
			expectedStatus:   http.StatusNotFound,
			expectedCode:     errors.CodeEntityNotFound,
		},
		{
			name:             "ConflictError maps to 409 Conflict",
			applicationError: errors.NewConflictError("user", "email exists", "existing-id"),
			expectedStatus:   http.StatusConflict,
			expectedCode:     errors.CodeResourceConflict,
		},
		{
			name:             "UnauthorizedError maps to 401 Unauthorized",
			applicationError: errors.NewUnauthorizedError("delete", "user-123", "no permission"),
			expectedStatus:   http.StatusUnauthorized,
			expectedCode:     errors.CodeUnauthorized,
		},
		{
			name:             "BusinessLogicError maps to 422 Unprocessable Entity",
			applicationError: errors.NewBusinessLogicError("registration", "age restriction"),
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedCode:     errors.CodeBusinessLogicError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpErr := mapper.MapToHTTPError(tt.applicationError, traceID)

			assert.Equal(t, tt.expectedStatus, httpErr.StatusCode)
			assert.Equal(t, tt.expectedCode, httpErr.Code())
			assert.Equal(t, traceID, httpErr.TraceID)
			assert.NotEmpty(t, httpErr.Message)
			assert.NotEmpty(t, httpErr.Details())
		})
	}
}

func TestErrorMapper_MapInfrastructureErrors(t *testing.T) {
	mapper := errors.NewErrorMapper()
	traceID := "test-trace-123"

	tests := []struct {
		name                string
		infrastructureError errors.BaseError
		expectedStatus      int
		expectedCode        errors.ErrorCode
		expectedRetryable   bool
	}{
		{
			name:                "DatabaseError maps to 503 Service Unavailable",
			infrastructureError: errors.NewDatabaseError("create", "users", nil, true),
			expectedStatus:      http.StatusServiceUnavailable,
			expectedCode:        errors.CodeDatabaseError,
			expectedRetryable:   true,
		},
		{
			name:                "NetworkError maps to 503 Service Unavailable",
			infrastructureError: errors.NewNetworkError("user-service", "/users", "POST", nil, true),
			expectedStatus:      http.StatusServiceUnavailable,
			expectedCode:        errors.CodeNetworkError,
			expectedRetryable:   true,
		},
		{
			name:                "ConfigurationError maps to 500 Internal Server Error",
			infrastructureError: errors.NewConfigurationError("database", "host", "", "missing host"),
			expectedStatus:      http.StatusInternalServerError,
			expectedCode:        errors.CodeConfigurationError,
			expectedRetryable:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpErr := mapper.MapToHTTPError(tt.infrastructureError, traceID)

			assert.Equal(t, tt.expectedStatus, httpErr.StatusCode)
			assert.Equal(t, tt.expectedCode, httpErr.Code())
			assert.Equal(t, traceID, httpErr.TraceID)
			assert.NotEmpty(t, httpErr.Message)
			assert.NotEmpty(t, httpErr.Details())

			// Check retryable flag in details
			retryable, exists := httpErr.Details()["retryable"]
			assert.True(t, exists)
			assert.Equal(t, tt.expectedRetryable, retryable)
		})
	}
}

func TestErrorMapper_MapUnknownError(t *testing.T) {
	mapper := errors.NewErrorMapper()
	traceID := "test-trace-123"

	// Test with a standard Go error (not our custom error types)
	unknownErr := assert.AnError
	httpErr := mapper.MapToHTTPError(unknownErr, traceID)

	assert.Equal(t, http.StatusInternalServerError, httpErr.StatusCode)
	assert.Equal(t, errors.CodeInternalError, httpErr.Code())
	assert.Equal(t, "An internal server error occurred", httpErr.Message)
	assert.Equal(t, traceID, httpErr.TraceID)

	details := httpErr.Details()
	assert.Contains(t, details, "original_error")
	assert.Equal(t, unknownErr.Error(), details["original_error"])
}

func TestGetHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "HTTPError returns its status code",
			err:            errors.NewHTTPError(http.StatusBadRequest, errors.CodeValidationError, "test", nil, ""),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "RequiredFieldError returns 400",
			err:            errors.NewRequiredFieldError("field", "value"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "EntityNotFoundError returns 404",
			err:            errors.NewEntityNotFoundError("User", "123"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "DatabaseError returns 503",
			err:            errors.NewDatabaseError("create", "users", nil, true),
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "Unknown error returns 500",
			err:            assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := errors.GetHTTPStatusCode(tt.err)
			assert.Equal(t, tt.expectedStatus, statusCode)
		})
	}
}
