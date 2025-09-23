package builder

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cctw-zed/wonder/internal/infrastructure/config"
	"github.com/cctw-zed/wonder/pkg/logger"
)

func TestUserBuilder_Default(t *testing.T) {
	user := NewUserBuilder().Build()

	assert.Equal(t, "test-user-123", user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestUserBuilder_WithCustomValues(t *testing.T) {
	customTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	user := NewUserBuilder().
		WithID("custom-id").
		WithEmail("custom@example.com").
		WithName("Custom User").
		WithTimestamps(customTime, customTime).
		Build()

	assert.Equal(t, "custom-id", user.ID)
	assert.Equal(t, "custom@example.com", user.Email)
	assert.Equal(t, "Custom User", user.Name)
	assert.Equal(t, customTime, user.CreatedAt)
	assert.Equal(t, customTime, user.UpdatedAt)
}

func TestUserBuilder_Valid(t *testing.T) {
	// Initialize logger for tests
	cfg := &config.Config{
		Log: &config.LogConfig{
			Level:       "debug",
			Format:      "text",
			ServiceName: "wonder-test",
		},
	}
	logger.InitializeGlobalLogger(cfg)
	user := NewUserBuilder().Valid().Build()

	err := user.Validate(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "valid-id-123", user.ID)
	assert.Equal(t, "valid@example.com", user.Email)
	assert.Equal(t, "Valid User", user.Name)
}

func TestUserBuilder_Invalid(t *testing.T) {
	user := NewUserBuilder().Invalid().Build()

	err := user.Validate(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "id is required")
}

func TestUserBuilder_WithInvalidEmail(t *testing.T) {
	user := NewUserBuilder().WithInvalidEmail().Build()

	err := user.Validate(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format for email, expected: valid email address")
}

func TestUserBuilder_BuildMany(t *testing.T) {
	count := 3
	users := NewUserBuilder().BuildMany(count)

	require.Len(t, users, count)

	for i, user := range users {
		expectedID := "test-user-123-" + string(rune('0'+i))
		expectedEmail := "test@example.com" + string(rune('0'+i))
		expectedName := "Test User " + string(rune('0'+i))

		assert.Equal(t, expectedID, user.ID)
		assert.Equal(t, expectedEmail, user.Email)
		assert.Equal(t, expectedName, user.Name)
	}
}

func TestUserBuilder_IsolatedInstances(t *testing.T) {
	builder := NewUserBuilder()

	user1 := builder.WithName("User 1").Build()
	user2 := builder.WithName("User 2").Build()

	// Verify that changing one doesn't affect the other
	assert.Equal(t, "User 1", user1.Name)
	assert.Equal(t, "User 2", user2.Name)
}

func TestUserBuilderForTesting_ValidUser(t *testing.T) {
	testBuilder := NewUserBuilderForTesting()
	user := testBuilder.ValidUser()

	err := user.Validate(context.Background())
	require.NoError(t, err)

	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Name)
}

func TestUserBuilderForTesting_ValidUserWithEmail(t *testing.T) {
	testBuilder := NewUserBuilderForTesting()
	email := "specific@test.com"

	user := testBuilder.ValidUserWithEmail(email)

	err := user.Validate(context.Background())
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
}

func TestUserBuilderForTesting_ValidUserWithID(t *testing.T) {
	testBuilder := NewUserBuilderForTesting()
	id := "specific-id-456"

	user := testBuilder.ValidUserWithID(id)

	err := user.Validate(context.Background())
	require.NoError(t, err)
	assert.Equal(t, id, user.ID)
}

func TestUserBuilderForTesting_InvalidUsers(t *testing.T) {
	testBuilder := NewUserBuilderForTesting()

	t.Run("empty fields", func(t *testing.T) {
		user := testBuilder.InvalidUserEmptyFields()
		err := user.Validate(context.Background())
		require.Error(t, err)
	})

	t.Run("bad email", func(t *testing.T) {
		user := testBuilder.InvalidUserBadEmail()
		err := user.Validate(context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format for email, expected: valid email address")
	})
}

func TestUserBuilderForTesting_UserFromRegistration(t *testing.T) {
	testBuilder := NewUserBuilderForTesting()
	email := "register@test.com"
	name := "Registered User"

	user := testBuilder.UserFromRegistration(email, name)

	assert.Equal(t, "generated-id-123", user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())

	err := user.Validate(context.Background())
	require.NoError(t, err)
}