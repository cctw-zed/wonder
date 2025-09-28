package user

import (
	"context"
	"regexp"
	"time"

	"github.com/cctw-zed/wonder/pkg/errors"
	"github.com/cctw-zed/wonder/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// User 用户聚合根
type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Email        string    `gorm:"uniqueIndex:idx_users_email_unique;type:varchar(255);not null" json:"email"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
}

// UserService 用户领域服务接口
type UserService interface {
	Register(ctx context.Context, email, name, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetProfile(ctx context.Context, id string) (*User, error)
	UpdateProfile(ctx context.Context, id string, req *UpdateProfileRequest) (*User, error)
	ChangePassword(ctx context.Context, id string, oldPassword, newPassword string) error
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ListUsersRequest represents the request to list users with pagination
type ListUsersRequest struct {
	Page     int    `json:"page" binding:"min=1"`
	PageSize int    `json:"page_size" binding:"min=1,max=100"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
}

// ListUsersResponse represents the response for list users
type ListUsersResponse struct {
	Users      []*User `json:"users"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
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

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(ctx context.Context, email string) error {
	log := logger.Get().WithLayer("domain").WithComponent("user")

	if log.DebugEnabled() {
		log.Debug(ctx, "updating user email", "user_id", u.ID, "old_email", u.Email, "new_email", email)
	}

	if email == "" {
		return errors.NewRequiredFieldError("email", email)
	}

	// Create temporary user to validate email format
	tempUser := &User{Email: email}
	if !tempUser.IsEmailValid() {
		return errors.NewInvalidFormatError("email", email, "valid email address")
	}

	oldEmail := u.Email
	u.Email = email

	log.Info(ctx, "user email updated", "user_id", u.ID, "old_email", oldEmail, "new_email", email)
	return nil
}

// SetPassword sets the user's password (hashes it)
func (u *User) SetPassword(ctx context.Context, password string) error {
	log := logger.Get().WithLayer("domain").WithComponent("user")

	if log.DebugEnabled() {
		log.Debug(ctx, "setting user password", "user_id", u.ID)
	}

	if password == "" {
		return errors.NewRequiredFieldError("password", password)
	}

	if len(password) < 6 {
		return errors.NewInvalidFormatError("password", password, "at least 6 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(ctx, "failed to hash password", "error", err, "user_id", u.ID)
		return errors.NewBusinessLogicError("password_hashing", "password hashing failed")
	}

	u.PasswordHash = string(hashedPassword)
	log.Info(ctx, "user password updated", "user_id", u.ID)
	return nil
}

// CheckPassword verifies the password against the stored hash
func (u *User) CheckPassword(ctx context.Context, password string) error {
	log := logger.Get().WithLayer("domain").WithComponent("user")

	if log.DebugEnabled() {
		log.Debug(ctx, "checking user password", "user_id", u.ID)
	}

	if password == "" {
		return errors.NewRequiredFieldError("password", password)
	}

	if u.PasswordHash == "" {
		log.Warn(ctx, "user has no password set", "user_id", u.ID)
		return errors.NewInvalidStateError(errors.CodeBusinessLogicError, "user", u.ID, "password not set")
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		log.Warn(ctx, "password verification failed", "user_id", u.ID)
		return errors.NewUnauthorizedError("password_verification", u.ID, "invalid password")
	}

	log.Debug(ctx, "password verification successful", "user_id", u.ID)
	return nil
}
