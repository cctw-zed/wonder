package errors

import (
	"net/http"
)

// HTTPError represents errors at the HTTP interface layer
type HTTPError struct {
	StatusCode   int                    `json:"status_code"`
	ErrorCode    ErrorCode              `json:"code"`
	Message      string                 `json:"message"`
	ErrorDetails map[string]interface{} `json:"details,omitempty"`
	TraceID      string                 `json:"trace_id,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Code() ErrorCode {
	return e.ErrorCode
}

func (e *HTTPError) Type() ErrorType {
	return ErrorTypeInterface
}

func (e *HTTPError) Details() map[string]interface{} {
	return e.ErrorDetails
}

func (e *HTTPError) WithContext(key string, value interface{}) BaseError {
	if e.ErrorDetails == nil {
		e.ErrorDetails = make(map[string]interface{})
	}
	e.ErrorDetails[key] = value
	return e
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(statusCode int, code ErrorCode, message string, details map[string]interface{}, traceID string) *HTTPError {
	return &HTTPError{
		StatusCode:   statusCode,
		ErrorCode:    code,
		Message:      message,
		ErrorDetails: details,
		TraceID:      traceID,
	}
}

// ErrorMapper maps domain/application/infrastructure errors to HTTP errors
type ErrorMapper struct{}

// NewErrorMapper creates a new error mapper
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapToHTTPError maps various error types to HTTP errors
func (m *ErrorMapper) MapToHTTPError(err error, traceID string) *HTTPError {
	// Use the new error classification system
	errorType := Classifier.ClassifyError(err)

	switch errorType {
	case ErrorTypeDomain:
		if baseErr, ok := err.(BaseError); ok {
			return m.mapDomainError(baseErr, traceID)
		}
	case ErrorTypeApplication:
		if baseErr, ok := err.(BaseError); ok {
			return m.mapApplicationError(baseErr, traceID)
		}
	case ErrorTypeInfrastructure:
		if baseErr, ok := err.(BaseError); ok {
			return m.mapInfrastructureError(baseErr, traceID)
		}
	}

	// Unknown error - map to generic internal server error
	return NewHTTPError(
		http.StatusInternalServerError,
		CodeInternalError,
		"An internal server error occurred",
		map[string]interface{}{"original_error": err.Error()},
		traceID,
	)
}

// mapDomainError maps domain layer errors to HTTP errors
func (m *ErrorMapper) mapDomainError(err BaseError, traceID string) *HTTPError {
	switch err.Code() {
	case CodeValidationError, CodeRequiredField, CodeInvalidFormat, CodeInvalidValue, CodeOutOfRange:
		return NewHTTPError(
			http.StatusBadRequest,
			err.Code(),
			"Validation failed",
			err.Details(),
			traceID,
		)
	case CodeDomainRuleViolation, CodeBusinessRuleError, CodeInvariantViolation:
		return NewHTTPError(
			http.StatusUnprocessableEntity,
			err.Code(),
			"Business rule violation",
			err.Details(),
			traceID,
		)
	case CodeInvalidState, CodeStateTransition, CodePreconditionError:
		return NewHTTPError(
			http.StatusConflict,
			err.Code(),
			"Invalid entity state",
			err.Details(),
			traceID,
		)
	default:
		return NewHTTPError(
			http.StatusBadRequest,
			err.Code(),
			err.Error(),
			err.Details(),
			traceID,
		)
	}
}

// mapApplicationError maps application layer errors to HTTP errors
func (m *ErrorMapper) mapApplicationError(err BaseError, traceID string) *HTTPError {
	switch err.Code() {
	case CodeEntityNotFound:
		return NewHTTPError(
			http.StatusNotFound,
			err.Code(),
			"Resource not found",
			err.Details(),
			traceID,
		)
	case CodeResourceConflict, CodeDuplicateEntry:
		return NewHTTPError(
			http.StatusConflict,
			err.Code(),
			"Resource conflict",
			err.Details(),
			traceID,
		)
	case CodeResourceLocked:
		return NewHTTPError(
			http.StatusLocked,
			err.Code(),
			"Resource locked",
			err.Details(),
			traceID,
		)
	case CodeUnauthorized, CodeTokenExpired:
		return NewHTTPError(
			http.StatusUnauthorized,
			err.Code(),
			"Unauthorized access",
			err.Details(),
			traceID,
		)
	case CodeForbidden, CodeInsufficientRole:
		return NewHTTPError(
			http.StatusForbidden,
			err.Code(),
			"Access forbidden",
			err.Details(),
			traceID,
		)
	case CodeBusinessLogicError, CodeOperationFailed:
		return NewHTTPError(
			http.StatusUnprocessableEntity,
			err.Code(),
			"Business logic error",
			err.Details(),
			traceID,
		)
	case CodeQuotaExceeded, CodeRateLimitExceeded:
		return NewHTTPError(
			http.StatusTooManyRequests,
			err.Code(),
			"Rate limit exceeded",
			err.Details(),
			traceID,
		)
	default:
		return NewHTTPError(
			http.StatusInternalServerError,
			err.Code(),
			err.Error(),
			err.Details(),
			traceID,
		)
	}
}

// mapInfrastructureError maps infrastructure layer errors to HTTP errors
func (m *ErrorMapper) mapInfrastructureError(err BaseError, traceID string) *HTTPError {
	details := err.Details()

	// Add retryable information if available
	if retryable := IsRetryable(err); retryable {
		details["retryable"] = retryable
	}

	switch err.Code() {
	case CodeDatabaseError, CodeDatabaseConnection, CodeDatabaseTimeout, CodeDatabaseDeadlock:
		return NewHTTPError(
			http.StatusServiceUnavailable,
			err.Code(),
			"Database service unavailable",
			details,
			traceID,
		)
	case CodeNetworkError, CodeConnectionRefused, CodeConnectionTimeout, CodeServiceUnavailable:
		return NewHTTPError(
			http.StatusServiceUnavailable,
			err.Code(),
			"Network service unavailable",
			details,
			traceID,
		)
	case CodeExternalServiceError, CodeAPICallFailed, CodeExternalTimeout:
		return NewHTTPError(
			http.StatusServiceUnavailable,
			err.Code(),
			"External service unavailable",
			details,
			traceID,
		)
	case CodeConfigurationError, CodeMissingConfig, CodeInvalidConfig:
		return NewHTTPError(
			http.StatusInternalServerError,
			err.Code(),
			"Service configuration error",
			details,
			traceID,
		)
	default:
		return NewHTTPError(
			http.StatusInternalServerError,
			err.Code(),
			"Infrastructure error",
			details,
			traceID,
		)
	}
}

// GetHTTPStatusCode extracts HTTP status code from error
func GetHTTPStatusCode(err error) int {
	switch e := err.(type) {
	case *HTTPError:
		return e.StatusCode
	default:
		mapper := NewErrorMapper()
		httpErr := mapper.MapToHTTPError(err, "")
		return httpErr.StatusCode
	}
}
