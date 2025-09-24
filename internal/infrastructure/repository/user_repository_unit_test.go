package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/pkg/logger"
)

func TestUserRepository_InputValidation(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

	// Create repository with nil db to test parameter validation
	// We need to bypass the constructor validation for this specific test
	repo := &userRepository{
		db:  nil,
		log: logger.Get().WithLayer("infrastructure").WithComponent("user_repository"),
	}
	ctx := context.Background()

	t.Run("Create method input validation", func(t *testing.T) {
		tests := []struct {
			name    string
			user    *user.User
			wantErr string
		}{
			{
				name:    "nil user",
				user:    nil,
				wantErr: "database error in operation 'create' on table 'users'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.Create(ctx, tt.user)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})

	t.Run("GetByID method input validation", func(t *testing.T) {
		tests := []struct {
			name    string
			id      string
			wantErr string
		}{
			{
				name:    "empty ID",
				id:      "",
				wantErr: "validation failed for field 'id': id is required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := repo.GetByID(ctx, tt.id)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})

	t.Run("GetByEmail method input validation", func(t *testing.T) {
		tests := []struct {
			name    string
			email   string
			wantErr string
		}{
			{
				name:    "empty email",
				email:   "",
				wantErr: "validation failed for field 'email': email is required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := repo.GetByEmail(ctx, tt.email)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})

	t.Run("Update method input validation", func(t *testing.T) {
		tests := []struct {
			name    string
			user    *user.User
			wantErr string
		}{
			{
				name:    "nil user",
				user:    nil,
				wantErr: "user cannot be nil",
			},
			{
				name: "user with empty ID",
				user: &user.User{
					ID:    "",
					Email: "test@example.com",
					Name:  "Test User",
				},
				wantErr: "user ID cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.Update(ctx, tt.user)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})

	t.Run("Delete method input validation", func(t *testing.T) {
		tests := []struct {
			name    string
			id      string
			wantErr string
		}{
			{
				name:    "empty ID",
				id:      "",
				wantErr: "user ID cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.Delete(ctx, tt.id)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})
}
