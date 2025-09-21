package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/cctw-zed/wonder/internal/domain/user"
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
		return fmt.Errorf("user cannot be nil")
	}

	// Validate user before creating
	if err := u.Validate(); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
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
			return fmt.Errorf("user with email %s already exists", u.Email)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &u, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
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
	errorStr := err.Error()
	return contains(errorStr, "duplicate key value violates unique constraint") ||
		contains(errorStr, "UNIQUE constraint failed") ||
		contains(errorStr, "23505") // PostgreSQL unique violation error code
}

// contains checks if a string contains a substring (case-insensitive)
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr || len(substr) == 0 ||
			anySubstring(str, substr))
}

func anySubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
