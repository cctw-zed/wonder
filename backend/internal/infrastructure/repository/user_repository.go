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
	"github.com/cctw-zed/wonder/pkg/logger"
)

type userRepository struct {
	db  *gorm.DB
	log logger.Logger
}

// NewUserRepository creates a new UserRepository implementation
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return NewUserRepositoryWithLogger(db, logger.Get().WithLayer("infrastructure").WithComponent("user_repository"))
}

// NewUserRepositoryWithLogger creates a new UserRepository implementation with explicit logger
func NewUserRepositoryWithLogger(db *gorm.DB, log logger.Logger) user.UserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}
	if log == nil {
		panic("logger cannot be nil")
	}

	return &userRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		r.log.Error(ctx, "user cannot be nil")
		return wonderErrors.NewDatabaseError("create", "users", nil, false, map[string]interface{}{
			"reason": "user cannot be nil",
		})
	}

	if r.log.DebugEnabled() {
		r.log.Debug(ctx, "creating user", "user_id", u.ID, "email", u.Email)
	}

	// Domain validation is handled by aggregate, but we double-check here
	if err := u.Validate(ctx); err != nil {
		return err
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
		if isDuplicateKeyError(err) {
			r.log.Warn(ctx, "duplicate email", "email", u.Email)
			return wonderErrors.NewConflictError("user", "email already exists", "", map[string]interface{}{
				"email": u.Email,
			})
		}
		r.log.Error(ctx, "database create failed", "error", err, "retryable", isRetryableError(err))
		return wonderErrors.NewDatabaseError("create", "users", err, isRetryableError(err), map[string]interface{}{
			"user_id": u.ID,
			"email":   u.Email,
		})
	}

	r.log.Info(ctx, "user created", "user_id", u.ID)
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

	if r.log.DebugEnabled() {
		r.log.Debug(ctx, "querying user by email", "email", email)
	}

	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found (application layer will handle)
		}
		r.log.Error(ctx, "email query failed", "error", err, "email", email)
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
	if err := u.Validate(ctx); err != nil {
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

// List retrieves users with pagination and filtering
func (r *userRepository) List(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	if req == nil {
		return nil, wonderErrors.NewRequiredFieldError("request", "nil")
	}

	// Set default pagination values
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	if r.log.DebugEnabled() {
		r.log.Debug(ctx, "listing users", "page", page, "page_size", pageSize, "email_filter", req.Email, "name_filter", req.Name)
	}

	// Build query with filters
	query := r.db.WithContext(ctx).Model(&user.User{})

	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}

	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	// Get total count
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		r.log.Error(ctx, "failed to count users", "error", err)
		return nil, wonderErrors.NewDatabaseError("count", "users", err, isRetryableError(err), map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		})
	}

	// Get users with pagination
	var users []*user.User
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		r.log.Error(ctx, "failed to list users", "error", err)
		return nil, wonderErrors.NewDatabaseError("list", "users", err, isRetryableError(err), map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		})
	}

	// Calculate total pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	r.log.Info(ctx, "users listed successfully", "total", total, "page", page, "page_size", pageSize, "returned_count", len(users))

	return &user.ListUsersResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
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
