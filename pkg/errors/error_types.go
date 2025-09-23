package errors

// ErrorType represents the layer where the error originates
type ErrorType string

const (
	ErrorTypeDomain         ErrorType = "domain"
	ErrorTypeApplication    ErrorType = "application"
	ErrorTypeInfrastructure ErrorType = "infrastructure"
	ErrorTypeInterface      ErrorType = "interface"
	ErrorTypeSystem         ErrorType = "system"
)

// BaseError provides common functionality for all custom errors
type BaseError interface {
	error
	Code() ErrorCode
	Type() ErrorType
	Details() map[string]interface{}
	WithContext(key string, value interface{}) BaseError
}

// ErrorClassifier helps identify error types without redundant methods
type ErrorClassifier struct{}

// ClassifyError determines the error type and provides type-safe checking
func (c *ErrorClassifier) ClassifyError(err error) ErrorType {
	switch err.(type) {
	case *ValidationError, *DomainRuleError, *InvalidStateError:
		return ErrorTypeDomain
	case *EntityNotFoundError, *ConflictError, *UnauthorizedError, *BusinessLogicError:
		return ErrorTypeApplication
	case *DatabaseError, *NetworkError, *ExternalServiceError, *ConfigurationError:
		return ErrorTypeInfrastructure
	case *HTTPError:
		return ErrorTypeInterface
	default:
		return ErrorTypeSystem
	}
}

// IsDomainError checks if error is from domain layer
func (c *ErrorClassifier) IsDomainError(err error) bool {
	return c.ClassifyError(err) == ErrorTypeDomain
}

// IsApplicationError checks if error is from application layer
func (c *ErrorClassifier) IsApplicationError(err error) bool {
	return c.ClassifyError(err) == ErrorTypeApplication
}

// IsInfrastructureError checks if error is from infrastructure layer
func (c *ErrorClassifier) IsInfrastructureError(err error) bool {
	return c.ClassifyError(err) == ErrorTypeInfrastructure
}

// IsRetryable determines if an error should be retried
func (c *ErrorClassifier) IsRetryable(err error) bool {
	switch e := err.(type) {
	case *DatabaseError:
		return e.IsRetryable
	case *NetworkError:
		return e.IsRetryable
	case *ExternalServiceError:
		return e.IsRetryable
	case *ConfigurationError:
		return false // Configuration errors are never retryable
	default:
		return false
	}
}

// Global classifier instance
var Classifier = &ErrorClassifier{}

// Convenience functions for error classification
func IsDomainError(err error) bool {
	return Classifier.IsDomainError(err)
}

func IsApplicationError(err error) bool {
	return Classifier.IsApplicationError(err)
}

func IsInfrastructureError(err error) bool {
	return Classifier.IsInfrastructureError(err)
}

func IsRetryable(err error) bool {
	return Classifier.IsRetryable(err)
}