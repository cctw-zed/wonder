package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cctw-zed/wonder/pkg/errors"
)

func TestValidationError(t *testing.T) {
	t.Run("Create validation error", func(t *testing.T) {
		err := errors.NewInvalidFormatError("email", "invalid-email", "valid email format")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format for email")
		assert.Equal(t, errors.CodeInvalidFormat, err.Code())
		assert.True(t, errors.IsDomainError(err))

		details := err.Details()
		assert.Equal(t, "email", details["field"])
		assert.Equal(t, "invalid-email", details["value"])
		assert.Contains(t, details["message"].(string), "valid email format")
	})

	t.Run("Create validation error with range", func(t *testing.T) {
		err := errors.NewOutOfRangeError("password", 3, 6, 50)

		details := err.Details()
		assert.Equal(t, "password", details["field"])
		assert.Equal(t, 3, details["value"])
		assert.Contains(t, details["message"].(string), "must be between 6 and 50")
	})
}

func TestDomainRuleError(t *testing.T) {
	t.Run("Create domain rule error", func(t *testing.T) {
		err := errors.NewBusinessRuleError("user_uniqueness", "user with this email already exists")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "domain rule violation")
		assert.Equal(t, errors.CodeBusinessRuleError, err.Code())
		assert.True(t, errors.IsDomainError(err))

		details := err.Details()
		assert.Equal(t, "user_uniqueness", details["rule"])
		assert.Equal(t, "user with this email already exists", details["message"])
	})

	t.Run("Create domain rule error with context", func(t *testing.T) {
		context := map[string]interface{}{
			"email":   "test@example.com",
			"user_id": "existing-user-123",
		}
		err := errors.NewBusinessRuleError("user_uniqueness", "user with this email already exists", context)

		details := err.Details()
		assert.Equal(t, "user_uniqueness", details["rule"])
		assert.Equal(t, "user with this email already exists", details["message"])
		assert.Equal(t, "test@example.com", details["email"])
		assert.Equal(t, "existing-user-123", details["user_id"])
	})
}

func TestInvalidStateError(t *testing.T) {
	t.Run("Create invalid state error", func(t *testing.T) {
		err := errors.NewStateTransitionError("User", "inactive", "active", "cannot perform operation on inactive user")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transition from inactive to active")
		assert.Equal(t, errors.CodeStateTransition, err.Code())
		assert.True(t, errors.IsDomainError(err))

		details := err.Details()
		assert.Equal(t, "User", details["entity"])
		assert.Equal(t, "inactive", details["state"])
		assert.Contains(t, details["message"].(string), "cannot transition from inactive to active")
	})

	t.Run("Create invalid state error with context", func(t *testing.T) {
		err := errors.NewStateTransitionError("User", "pending_activation", "active", "user must be active").WithContext("user_id", "user-123")

		details := err.Details()
		assert.Equal(t, "User", details["entity"])
		assert.Equal(t, "pending_activation", details["state"])
		assert.Contains(t, details["message"].(string), "cannot transition from pending_activation to active")
		assert.Equal(t, "user-123", details["user_id"])
	})
}

func TestDomainErrorInterface(t *testing.T) {
	tests := []struct {
		name string
		err  errors.BaseError
	}{
		{
			name: "RequiredFieldError implements BaseError",
			err:  errors.NewRequiredFieldError("field", "value"),
		},
		{
			name: "BusinessRuleError implements BaseError",
			err:  errors.NewBusinessRuleError("rule", "message"),
		},
		{
			name: "StateTransitionError implements BaseError",
			err:  errors.NewStateTransitionError("entity", "from", "to", "reason"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all domain errors implement the BaseError interface
			assert.NotEmpty(t, tt.err.Error())
			assert.NotEmpty(t, tt.err.Code())
			assert.NotNil(t, tt.err.Details())
			assert.True(t, errors.IsDomainError(tt.err))
		})
	}
}
