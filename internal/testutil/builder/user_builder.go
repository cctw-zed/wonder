package builder

import (
	"time"

	"github.com/cctw-zed/wonder/internal/domain/user"
)

// UserBuilder helps create test User instances with default values
type UserBuilder struct {
	user *user.User
}

// NewUserBuilder creates a new UserBuilder with default values
func NewUserBuilder() *UserBuilder {
	now := time.Now()
	return &UserBuilder{
		user: &user.User{
			ID:        "test-user-123",
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// WithID sets the user ID
func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.user.ID = id
	return b
}

// WithEmail sets the user email
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

// WithName sets the user name
func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

// WithCreatedAt sets the creation timestamp
func (b *UserBuilder) WithCreatedAt(createdAt time.Time) *UserBuilder {
	b.user.CreatedAt = createdAt
	return b
}

// WithUpdatedAt sets the update timestamp
func (b *UserBuilder) WithUpdatedAt(updatedAt time.Time) *UserBuilder {
	b.user.UpdatedAt = updatedAt
	return b
}

// WithTimestamps sets both created and updated timestamps
func (b *UserBuilder) WithTimestamps(createdAt, updatedAt time.Time) *UserBuilder {
	b.user.CreatedAt = createdAt
	b.user.UpdatedAt = updatedAt
	return b
}

// Valid creates a user that passes validation
func (b *UserBuilder) Valid() *UserBuilder {
	return b.WithID("valid-id-123").
		WithEmail("valid@example.com").
		WithName("Valid User")
}

// Invalid creates a user that fails validation (empty required fields)
func (b *UserBuilder) Invalid() *UserBuilder {
	return b.WithID("").
		WithEmail("").
		WithName("")
}

// WithInvalidEmail creates a user with invalid email format
func (b *UserBuilder) WithInvalidEmail() *UserBuilder {
	return b.WithEmail("invalid-email-format")
}

// Build returns the constructed User
func (b *UserBuilder) Build() *user.User {
	// Create a copy to avoid sharing state between tests
	return &user.User{
		ID:        b.user.ID,
		Email:     b.user.Email,
		Name:      b.user.Name,
		CreatedAt: b.user.CreatedAt,
		UpdatedAt: b.user.UpdatedAt,
	}
}

// BuildPtr returns a pointer to the constructed User
func (b *UserBuilder) BuildPtr() *user.User {
	return b.Build()
}

// BuildMany creates multiple users with incremented IDs and emails
func (b *UserBuilder) BuildMany(count int) []*user.User {
	users := make([]*user.User, count)
	for i := 0; i < count; i++ {
		builder := NewUserBuilder().
			WithID(b.user.ID+"-"+string(rune('0'+i))).
			WithEmail(b.user.Email+string(rune('0'+i))).
			WithName(b.user.Name+" "+string(rune('0'+i))).
			WithTimestamps(b.user.CreatedAt, b.user.UpdatedAt)
		users[i] = builder.Build()
	}
	return users
}

// UserBuilderForTesting provides common test scenarios
type UserBuilderForTesting struct{}

// NewUserBuilderForTesting creates testing-specific builders
func NewUserBuilderForTesting() *UserBuilderForTesting {
	return &UserBuilderForTesting{}
}

// ValidUser creates a valid user for testing
func (t *UserBuilderForTesting) ValidUser() *user.User {
	return NewUserBuilder().Valid().Build()
}

// ValidUserWithEmail creates a valid user with specific email
func (t *UserBuilderForTesting) ValidUserWithEmail(email string) *user.User {
	return NewUserBuilder().Valid().WithEmail(email).Build()
}

// ValidUserWithID creates a valid user with specific ID
func (t *UserBuilderForTesting) ValidUserWithID(id string) *user.User {
	return NewUserBuilder().Valid().WithID(id).Build()
}

// InvalidUserEmptyFields creates a user with empty required fields
func (t *UserBuilderForTesting) InvalidUserEmptyFields() *user.User {
	return NewUserBuilder().Invalid().Build()
}

// InvalidUserBadEmail creates a user with invalid email format
func (t *UserBuilderForTesting) InvalidUserBadEmail() *user.User {
	return NewUserBuilder().WithInvalidEmail().Build()
}

// UserFromRegistration creates a user as if it came from registration
func (t *UserBuilderForTesting) UserFromRegistration(email, name string) *user.User {
	now := time.Now()
	return NewUserBuilder().
		WithID("generated-id-123").
		WithEmail(email).
		WithName(name).
		WithTimestamps(now, now).
		Build()
}
