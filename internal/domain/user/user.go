package user

import (
	"context"
	"regexp"
	"time"

	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/logger"
)

// User 用户聚合根
type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Email     string    `gorm:"uniqueIndex:idx_users_email_unique;type:varchar(255);not null" json:"email"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// UserService 用户领域服务接口
type UserService interface {
	Register(ctx context.Context, email, name string) (*User, error)
	//GetProfile(ctx context.Context, id string) (*User, error)
	//UpdateProfile(ctx context.Context, id, name string) error
}

// Validate validates the user entity
func (u *User) Validate(ctx context.Context) error {
	domainLogger := logger.NewDomainLogger(logger.NewLogger())

	// Log validation start
	domainLogger.Debug(ctx, "Starting user validation",
		logger.String("user_id", u.ID),
		logger.String("email", u.Email),
	)

	if u.ID == "" {
		domainLogger.LogBusinessRule(ctx, "user_id_required", "User", u.ID, false,
			logger.String("field", "id"),
		)
		return errors.NewRequiredFieldError("id", u.ID)
	}

	if u.Email == "" {
		domainLogger.LogBusinessRule(ctx, "user_email_required", "User", u.ID, false,
			logger.String("field", "email"),
		)
		return errors.NewRequiredFieldError("email", u.Email)
	}

	if !u.IsEmailValid() {
		domainLogger.LogBusinessRule(ctx, "user_email_format", "User", u.ID, false,
			logger.String("field", "email"),
			logger.String("invalid_email", u.Email),
		)
		return errors.NewInvalidFormatError("email", u.Email, "valid email address")
	}

	if u.Name == "" {
		domainLogger.LogBusinessRule(ctx, "user_name_required", "User", u.ID, false,
			logger.String("field", "name"),
		)
		return errors.NewRequiredFieldError("name", u.Name)
	}

	// Log successful validation
	domainLogger.LogBusinessRule(ctx, "user_validation", "User", u.ID, true)
	domainLogger.Debug(ctx, "User validation completed successfully",
		logger.String("user_id", u.ID),
		logger.String("email", u.Email),
		logger.String("name", u.Name),
	)

	return nil
}

// IsEmailValid checks if the email format is valid
func (u *User) IsEmailValid() bool {
	if u.Email == "" {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

// UpdateName updates the user's name
func (u *User) UpdateName(ctx context.Context, name string) error {
	domainLogger := logger.NewDomainLogger(logger.NewLogger())

	domainLogger.Debug(ctx, "Updating user name",
		logger.String("user_id", u.ID),
		logger.String("old_name", u.Name),
		logger.String("new_name", name),
	)

	if name == "" {
		domainLogger.LogBusinessRule(ctx, "user_name_required", "User", u.ID, false,
			logger.String("field", "name"),
			logger.String("attempted_value", name),
		)
		return errors.NewRequiredFieldError("name", name)
	}

	oldName := u.Name
	u.Name = name

	// Log domain event
	domainLogger.LogDomainEvent(ctx, "UserNameUpdated", u.ID,
		logger.String("old_name", oldName),
		logger.String("new_name", name),
	)

	domainLogger.LogAggregateChange(ctx, "User", u.ID, "UpdateName",
		logger.String("field", "name"),
		logger.String("old_value", oldName),
		logger.String("new_value", name),
	)

	return nil
}
