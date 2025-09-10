package user

import (
	"context"
	"time"
)

// User 用户聚合根
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
