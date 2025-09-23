package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cctw-zed/wonder/pkg/errors"
)

func TestErrorClassification(t *testing.T) {
	classifier := &errors.ErrorClassifier{}

	tests := []struct {
		name     string
		err      error
		expected errors.ErrorType
	}{
		{
			name:     "ValidationError is domain error",
			err:      errors.NewRequiredFieldError("email", ""),
			expected: errors.ErrorTypeDomain,
		},
		{
			name:     "BusinessRuleError is domain error",
			err:      errors.NewBusinessRuleError("uniqueness", "email must be unique"),
			expected: errors.ErrorTypeDomain,
		},
		{
			name:     "EntityNotFoundError is application error",
			err:      errors.NewEntityNotFoundError("User", "123"),
			expected: errors.ErrorTypeApplication,
		},
		{
			name:     "DuplicateEntryError is application error",
			err:      errors.NewDuplicateEntryError("User", "email", "test@example.com", "existing-123"),
			expected: errors.ErrorTypeApplication,
		},
		{
			name:     "DatabaseError is infrastructure error",
			err:      errors.NewDatabaseError("create", "users", nil, true),
			expected: errors.ErrorTypeInfrastructure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorType := classifier.ClassifyError(tt.err)
			assert.Equal(t, tt.expected, errorType)

			// Test convenience functions
			switch tt.expected {
			case errors.ErrorTypeDomain:
				assert.True(t, errors.IsDomainError(tt.err))
				assert.False(t, errors.IsApplicationError(tt.err))
				assert.False(t, errors.IsInfrastructureError(tt.err))
			case errors.ErrorTypeApplication:
				assert.False(t, errors.IsDomainError(tt.err))
				assert.True(t, errors.IsApplicationError(tt.err))
				assert.False(t, errors.IsInfrastructureError(tt.err))
			case errors.ErrorTypeInfrastructure:
				assert.False(t, errors.IsDomainError(tt.err))
				assert.False(t, errors.IsApplicationError(tt.err))
				assert.True(t, errors.IsInfrastructureError(tt.err))
			}
		})
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected errors.ErrorCode
	}{
		{
			name:     "Required field error has correct code",
			err:      errors.NewRequiredFieldError("email", ""),
			expected: errors.CodeRequiredField,
		},
		{
			name:     "Invalid format error has correct code",
			err:      errors.NewInvalidFormatError("email", "invalid", "email format"),
			expected: errors.CodeInvalidFormat,
		},
		{
			name:     "Duplicate entry error has correct code",
			err:      errors.NewDuplicateEntryError("User", "email", "test@example.com", "123"),
			expected: errors.CodeDuplicateEntry,
		},
		{
			name:     "Entity not found error has correct code",
			err:      errors.NewEntityNotFoundError("User", "123"),
			expected: errors.CodeEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseErr, ok := tt.err.(errors.BaseError)
			assert.True(t, ok, "Error should implement BaseError interface")
			assert.Equal(t, tt.expected, baseErr.Code())
		})
	}
}

func TestErrorCodeValidation(t *testing.T) {
	validCodes := []errors.ErrorCode{
		errors.CodeValidationError,
		errors.CodeRequiredField,
		errors.CodeInvalidFormat,
		errors.CodeEntityNotFound,
		errors.CodeDatabaseError,
		errors.CodeNetworkError,
	}

	for _, code := range validCodes {
		t.Run(string(code), func(t *testing.T) {
			assert.True(t, code.IsValid(), "Code %s should be valid", code)
		})
	}

	// Test invalid code
	invalidCode := errors.ErrorCode("INVALID_CODE")
	assert.False(t, invalidCode.IsValid(), "Invalid code should return false")
}

func TestErrorContext(t *testing.T) {
	err := errors.NewRequiredFieldError("email", "")

	// Test WithContext method
	enrichedErr := err.WithContext("operation", "user_registration")
	enrichedErr = enrichedErr.WithContext("timestamp", "2023-10-01T10:00:00Z")

	details := enrichedErr.Details()
	assert.Equal(t, "user_registration", details["operation"])
	assert.Equal(t, "2023-10-01T10:00:00Z", details["timestamp"])
	assert.Equal(t, "email", details["field"])
	assert.Equal(t, "", details["value"])
}

func TestConvenienceConstructors(t *testing.T) {
	t.Run("Domain error constructors", func(t *testing.T) {
		// Required field error
		err1 := errors.NewRequiredFieldError("email", "")
		assert.Equal(t, errors.CodeRequiredField, err1.Code())
		assert.Contains(t, err1.Error(), "email is required")

		// Invalid format error
		err2 := errors.NewInvalidFormatError("email", "invalid", "valid email")
		assert.Equal(t, errors.CodeInvalidFormat, err2.Code())
		assert.Contains(t, err2.Error(), "invalid format for email")

		// Business rule error
		err3 := errors.NewBusinessRuleError("uniqueness", "email must be unique")
		assert.Equal(t, errors.CodeBusinessRuleError, err3.Code())
		assert.Contains(t, err3.Error(), "uniqueness")

		// State transition error
		err4 := errors.NewStateTransitionError("User", "inactive", "active", "not verified")
		assert.Equal(t, errors.CodeStateTransition, err4.Code())
		assert.Contains(t, err4.Error(), "cannot transition from inactive to active")
	})

	t.Run("Application error constructors", func(t *testing.T) {
		// Entity not found error
		err1 := errors.NewEntityNotFoundError("User", "123")
		assert.Equal(t, errors.CodeEntityNotFound, err1.Code())
		assert.Contains(t, err1.Error(), "User with ID '123' not found")

		// Duplicate entry error
		err2 := errors.NewDuplicateEntryError("User", "email", "test@example.com", "existing-123")
		assert.Equal(t, errors.CodeDuplicateEntry, err2.Code())
		assert.Contains(t, err2.Error(), "email 'test@example.com' already exists")
	})
}

func TestRetryableErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "Retryable database error",
			err:       errors.NewDatabaseError("create", "users", nil, true),
			retryable: true,
		},
		{
			name:      "Non-retryable database error",
			err:       errors.NewDatabaseError("create", "users", nil, false),
			retryable: false,
		},
		{
			name:      "Configuration error is not retryable",
			err:       errors.NewConfigurationError("database", "host", "", "missing host"),
			retryable: false,
		},
		{
			name:      "Domain error is not retryable",
			err:       errors.NewRequiredFieldError("email", ""),
			retryable: false,
		},
		{
			name:      "Application error is not retryable",
			err:       errors.NewEntityNotFoundError("User", "123"),
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.retryable, errors.IsRetryable(tt.err))
		})
	}
}