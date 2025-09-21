package user

import (
	"context"
	"errors"
	"regexp"
	"time"
)

// User 用户聚合根
type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Email     string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
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
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user ID is required")
	}

	if u.Email == "" {
		return errors.New("user email is required")
	}

	if !u.IsEmailValid() {
		return errors.New("invalid email format")
	}

	if u.Name == "" {
		return errors.New("user name is required")
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
func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	u.Name = name
	return nil
}
