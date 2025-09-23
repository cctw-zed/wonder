package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/cctw-zed/wonder/internal/domain/user"
	wonderErrors "github.com/cctw-zed/wonder/pkg/errors"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return wonderErrors.NewDatabaseError("create", "users", nil, false, map[string]interface{}{
			"reason": "user cannot be nil",
		})
	}

	// Domain validation is handled by aggregate, but we double-check here
	if err := u.Validate(); err != nil {
		return err // Return the domain validation error directly
	}

	// Set timestamps if not already set
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	// Create user in database
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		// Check for unique constraint violation (email already exists)
		if isDuplicateKeyError(err) {
			return wonderErrors.NewConflictError("user", "email already exists", "", map[string]interface{}{
				"email": u.Email,
			})
		}
		// Other database errors
		return wonderErrors.NewDatabaseError("create", "users", err, isRetryableError(err), map[string]interface{}{
			"user_id": u.ID,
			"email":   u.Email,
		})
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, wonderErrors.NewRequiredFieldError("id", id)
	}

	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found (application layer will handle)
		}
		return nil, wonderErrors.NewDatabaseError("get_by_id", "users", err, isRetryableError(err), map[string]interface{}{
			"user_id": id,
		})
	}

	return &u, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, wonderErrors.NewRequiredFieldError("email", email)
	}

	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found (application layer will handle)
		}
		return nil, wonderErrors.NewDatabaseError("get_by_email", "users", err, isRetryableError(err), map[string]interface{}{
			"email": email,
		})
	}

	return &u, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}

	if u.ID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Validate user before updating
	if err := u.Validate(); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	// Update timestamp
	u.UpdatedAt = time.Now()

	// Update user in database
	result := r.db.WithContext(ctx).Save(u)
	if result.Error != nil {
		// Check for unique constraint violation
		if isDuplicateKeyError(result.Error) {
			return fmt.Errorf("user with email %s already exists", u.Email)
		}
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	// Check if user exists
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", u.ID)
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	result := r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	// Check if user was found and deleted
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}

// isDuplicateKeyError checks if the error is a duplicate key constraint violation
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// PostgreSQL duplicate key error codes
	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "duplicate key value violates unique constraint") ||
		strings.Contains(errorStr, "unique constraint failed") ||
		strings.Contains(errorStr, "23505") // PostgreSQL unique violation error code
}

// isRetryableError determines if a database error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())

	// Network or connection errors (retryable)
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"deadlock",
		"lock timeout",
		"temporary failure",
		"connection lost",
		"server is not available",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errorStr, pattern) {
			return true
		}
	}

	// Non-retryable errors
	nonRetryablePatterns := []string{
		"syntax error",
		"permission denied",
		"constraint violation",
		"invalid",
		"malformed",
	}

	for _, pattern := range nonRetryablePatterns {
		if strings.Contains(errorStr, pattern) {
			return false
		}
	}

	// Default to retryable for unknown database errors
	return true
}
