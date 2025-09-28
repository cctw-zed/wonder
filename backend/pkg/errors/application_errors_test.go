package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cctw-zed/wonder/pkg/errors"
)

func TestEntityNotFoundError(t *testing.T) {
	t.Run("Create entity not found error", func(t *testing.T) {
		err := errors.NewEntityNotFoundError("User", "user-123")

		assert.Error(t, err)
		assert.Equal(t, "User with ID 'user-123' not found", err.Error())
		assert.Equal(t, errors.CodeEntityNotFound, err.Code())
		assert.True(t, errors.IsApplicationError(err))

		details := err.Details()
		assert.Equal(t, "User", details["entity_type"])
		assert.Equal(t, "user-123", details["entity_id"])
	})

	t.Run("Create entity not found error with context", func(t *testing.T) {
		context := map[string]interface{}{
			"operation": "get_user_profile",
			"source":    "user_service",
		}
		err := errors.NewEntityNotFoundError("User", "user-123", context)

		details := err.Details()
		assert.Equal(t, "User", details["entity_type"])
		assert.Equal(t, "user-123", details["entity_id"])
		assert.Equal(t, "get_user_profile", details["operation"])
		assert.Equal(t, "user_service", details["source"])
	})
}

func TestConflictError(t *testing.T) {
	t.Run("Create conflict error", func(t *testing.T) {
		err := errors.NewConflictError("user", "email already registered", "existing-user-456")

		assert.Error(t, err)
		assert.Equal(t, "conflict with existing user: email already registered", err.Error())
		assert.Equal(t, errors.CodeResourceConflict, err.Code())
		assert.True(t, errors.IsApplicationError(err))

		details := err.Details()
		assert.Equal(t, "user", details["resource"])
		assert.Equal(t, "email already registered", details["reason"])
		assert.Equal(t, "existing-user-456", details["existing_id"])
	})

	t.Run("Create conflict error with context", func(t *testing.T) {
		context := map[string]interface{}{
			"email":        "test@example.com",
			"attempted_at": "2023-10-01T10:00:00Z",
		}
		err := errors.NewConflictError("user", "email already registered", "existing-user-456", context)

		details := err.Details()
		assert.Equal(t, "user", details["resource"])
		assert.Equal(t, "email already registered", details["reason"])
		assert.Equal(t, "existing-user-456", details["existing_id"])
		assert.Equal(t, "test@example.com", details["email"])
		assert.Equal(t, "2023-10-01T10:00:00Z", details["attempted_at"])
	})
}

func TestUnauthorizedError(t *testing.T) {
	t.Run("Create unauthorized error", func(t *testing.T) {
		err := errors.NewUnauthorizedError("delete_user", "user-123", "insufficient permissions")

		assert.Error(t, err)
		assert.Equal(t, "unauthorized to perform operation 'delete_user': insufficient permissions", err.Error())
		assert.Equal(t, errors.CodeUnauthorized, err.Code())
		assert.True(t, errors.IsApplicationError(err))

		details := err.Details()
		assert.Equal(t, "delete_user", details["operation"])
		assert.Equal(t, "user-123", details["user_id"])
		assert.Equal(t, "insufficient permissions", details["reason"])
	})

	t.Run("Create unauthorized error with context", func(t *testing.T) {
		context := map[string]interface{}{
			"required_role": "admin",
			"current_role":  "user",
			"resource_id":   "resource-789",
		}
		err := errors.NewUnauthorizedError("delete_user", "user-123", "insufficient permissions", context)

		details := err.Details()
		assert.Equal(t, "delete_user", details["operation"])
		assert.Equal(t, "user-123", details["user_id"])
		assert.Equal(t, "insufficient permissions", details["reason"])
		assert.Equal(t, "admin", details["required_role"])
		assert.Equal(t, "user", details["current_role"])
		assert.Equal(t, "resource-789", details["resource_id"])
	})
}

func TestBusinessLogicError(t *testing.T) {
	t.Run("Create business logic error", func(t *testing.T) {
		err := errors.NewBusinessLogicError("user_registration", "cannot register user under 18")

		assert.Error(t, err)
		assert.Equal(t, "business logic error in operation 'user_registration': cannot register user under 18", err.Error())
		assert.Equal(t, errors.CodeBusinessLogicError, err.Code())
		assert.True(t, errors.IsApplicationError(err))

		details := err.Details()
		assert.Equal(t, "user_registration", details["operation"])
		assert.Equal(t, "cannot register user under 18", details["reason"])
	})

	t.Run("Create business logic error with details", func(t *testing.T) {
		details := map[string]interface{}{
			"user_age":    17,
			"minimum_age": 18,
			"country":     "US",
		}
		err := errors.NewBusinessLogicError("user_registration", "cannot register user under 18", details)

		errorDetails := err.Details()
		assert.Equal(t, "user_registration", errorDetails["operation"])
		assert.Equal(t, "cannot register user under 18", errorDetails["reason"])
		assert.Equal(t, 17, errorDetails["user_age"])
		assert.Equal(t, 18, errorDetails["minimum_age"])
		assert.Equal(t, "US", errorDetails["country"])
	})
}

func TestApplicationErrorInterface(t *testing.T) {
	tests := []struct {
		name string
		err  errors.BaseError
	}{
		{
			name: "EntityNotFoundError implements BaseError",
			err:  errors.NewEntityNotFoundError("User", "user-123"),
		},
		{
			name: "ConflictError implements BaseError",
			err:  errors.NewConflictError("user", "email already exists", "existing-user"),
		},
		{
			name: "UnauthorizedError implements BaseError",
			err:  errors.NewUnauthorizedError("delete", "user-123", "insufficient permissions"),
		},
		{
			name: "BusinessLogicError implements BaseError",
			err:  errors.NewBusinessLogicError("registration", "age restriction"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all application errors implement the BaseError interface
			assert.NotEmpty(t, tt.err.Error())
			assert.NotEmpty(t, tt.err.Code())
			assert.NotNil(t, tt.err.Details())
			assert.True(t, errors.IsApplicationError(tt.err))
		})
	}
}
