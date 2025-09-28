package errors

// ErrorCode represents standardized error codes across the system
type ErrorCode string

// Domain error codes
const (
	// Validation errors
	CodeValidationError ErrorCode = "VALIDATION_ERROR"
	CodeRequiredField   ErrorCode = "REQUIRED_FIELD"
	CodeInvalidFormat   ErrorCode = "INVALID_FORMAT"
	CodeInvalidValue    ErrorCode = "INVALID_VALUE"
	CodeOutOfRange      ErrorCode = "OUT_OF_RANGE"

	// Domain rule errors
	CodeDomainRuleViolation ErrorCode = "DOMAIN_RULE_VIOLATION"
	CodeInvariantViolation  ErrorCode = "INVARIANT_VIOLATION"
	CodeBusinessRuleError   ErrorCode = "BUSINESS_RULE_ERROR"

	// State errors
	CodeInvalidState      ErrorCode = "INVALID_STATE"
	CodeStateTransition   ErrorCode = "INVALID_STATE_TRANSITION"
	CodePreconditionError ErrorCode = "PRECONDITION_ERROR"
)

// Application error codes
const (
	// Resource errors
	CodeEntityNotFound   ErrorCode = "ENTITY_NOT_FOUND"
	CodeResourceConflict ErrorCode = "RESOURCE_CONFLICT"
	CodeDuplicateEntry   ErrorCode = "DUPLICATE_ENTRY"
	CodeResourceLocked   ErrorCode = "RESOURCE_LOCKED"

	// Authorization errors
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeForbidden        ErrorCode = "FORBIDDEN"
	CodeInsufficientRole ErrorCode = "INSUFFICIENT_ROLE"
	CodeTokenExpired     ErrorCode = "TOKEN_EXPIRED"

	// Business logic errors
	CodeBusinessLogicError ErrorCode = "BUSINESS_LOGIC_ERROR"
	CodeOperationFailed    ErrorCode = "OPERATION_FAILED"
	CodeQuotaExceeded      ErrorCode = "QUOTA_EXCEEDED"
	CodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// Infrastructure error codes
const (
	// Database errors
	CodeDatabaseError       ErrorCode = "DATABASE_ERROR"
	CodeDatabaseConnection  ErrorCode = "DATABASE_CONNECTION_ERROR"
	CodeDatabaseTimeout     ErrorCode = "DATABASE_TIMEOUT"
	CodeDatabaseDeadlock    ErrorCode = "DATABASE_DEADLOCK"
	CodeTransactionRollback ErrorCode = "TRANSACTION_ROLLBACK"

	// Network errors
	CodeNetworkError       ErrorCode = "NETWORK_ERROR"
	CodeConnectionRefused  ErrorCode = "CONNECTION_REFUSED"
	CodeConnectionTimeout  ErrorCode = "CONNECTION_TIMEOUT"
	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// External service errors
	CodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
	CodeAPICallFailed        ErrorCode = "API_CALL_FAILED"
	CodeExternalTimeout      ErrorCode = "EXTERNAL_TIMEOUT"
	CodeInvalidResponse      ErrorCode = "INVALID_RESPONSE"

	// Configuration errors
	CodeConfigurationError ErrorCode = "CONFIGURATION_ERROR"
	CodeMissingConfig      ErrorCode = "MISSING_CONFIGURATION"
	CodeInvalidConfig      ErrorCode = "INVALID_CONFIGURATION"
)

// System error codes
const (
	CodeInternalError     ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeNotImplemented    ErrorCode = "NOT_IMPLEMENTED"
	CodeServiceStartup    ErrorCode = "SERVICE_STARTUP_ERROR"
	CodeDependencyMissing ErrorCode = "DEPENDENCY_MISSING"
)

// String returns the string representation of the error code
func (c ErrorCode) String() string {
	return string(c)
}

// IsValid checks if the error code is a known valid code
func (c ErrorCode) IsValid() bool {
	validCodes := map[ErrorCode]bool{
		// Domain codes
		CodeValidationError:     true,
		CodeRequiredField:       true,
		CodeInvalidFormat:       true,
		CodeInvalidValue:        true,
		CodeOutOfRange:          true,
		CodeDomainRuleViolation: true,
		CodeInvariantViolation:  true,
		CodeBusinessRuleError:   true,
		CodeInvalidState:        true,
		CodeStateTransition:     true,
		CodePreconditionError:   true,

		// Application codes
		CodeEntityNotFound:     true,
		CodeResourceConflict:   true,
		CodeDuplicateEntry:     true,
		CodeResourceLocked:     true,
		CodeUnauthorized:       true,
		CodeForbidden:          true,
		CodeInsufficientRole:   true,
		CodeTokenExpired:       true,
		CodeBusinessLogicError: true,
		CodeOperationFailed:    true,
		CodeQuotaExceeded:      true,
		CodeRateLimitExceeded:  true,

		// Infrastructure codes
		CodeDatabaseError:        true,
		CodeDatabaseConnection:   true,
		CodeDatabaseTimeout:      true,
		CodeDatabaseDeadlock:     true,
		CodeTransactionRollback:  true,
		CodeNetworkError:         true,
		CodeConnectionRefused:    true,
		CodeConnectionTimeout:    true,
		CodeServiceUnavailable:   true,
		CodeExternalServiceError: true,
		CodeAPICallFailed:        true,
		CodeExternalTimeout:      true,
		CodeInvalidResponse:      true,
		CodeConfigurationError:   true,
		CodeMissingConfig:        true,
		CodeInvalidConfig:        true,

		// System codes
		CodeInternalError:     true,
		CodeNotImplemented:    true,
		CodeServiceStartup:    true,
		CodeDependencyMissing: true,
	}

	return validCodes[c]
}
