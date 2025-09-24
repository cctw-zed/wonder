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
	log := logger.Get().WithLayer("domain").WithComponent("user")

	if log.DebugEnabled() {
		log.Debug(ctx, "validating user", "user_id", u.ID, "email", u.Email)
	}

	if u.ID == "" {
		return errors.NewRequiredFieldError("id", u.ID)
	}

	if u.Email == "" {
		return errors.NewRequiredFieldError("email", u.Email)
	}

	if !u.IsEmailValid() {
		return errors.NewInvalidFormatError("email", u.Email, "valid email address")
	}

	if u.Name == "" {
		return errors.NewRequiredFieldError("name", u.Name)
	}

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
	log := logger.Get().WithLayer("domain").WithComponent("user")

	if log.DebugEnabled() {
		log.Debug(ctx, "updating user name", "user_id", u.ID, "old_name", u.Name, "new_name", name)
	}

	if name == "" {
		return errors.NewRequiredFieldError("name", name)
	}

	oldName := u.Name
	u.Name = name

	log.Info(ctx, "user name updated", "user_id", u.ID, "old_name", oldName, "new_name", name)
	return nil
}
