package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cctw-zed/wonder/pkg/logger"
)

func TestUser_Creation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		email     string
		userName  string
		createdAt time.Time
		updatedAt time.Time
		want      *User
	}{
		{
			name:      "valid user creation",
			id:        "user123",
			email:     "test@example.com",
			userName:  "Test User",
			createdAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &User{
				ID:        "user123",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "user with empty name",
			id:        "user456",
			email:     "empty@example.com",
			userName:  "",
			createdAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &User{
				ID:        "user456",
				Email:     "empty@example.com",
				Name:      "",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				ID:        tt.id,
				Email:     tt.email,
				Name:      tt.userName,
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
			}

			assert.Equal(t, tt.want.ID, user.ID)
			assert.Equal(t, tt.want.Email, user.Email)
			assert.Equal(t, tt.want.Name, user.Name)
			assert.Equal(t, tt.want.CreatedAt, user.CreatedAt)
			assert.Equal(t, tt.want.UpdatedAt, user.UpdatedAt)
		})
	}
}

func TestUser_Validation(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: &User{
				ID:        "user123",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			user: &User{
				ID:        "",
				Email:     "test@example.com",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "empty email",
			user: &User{
				ID:        "user123",
				Email:     "",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid email format",
			user: &User{
				ID:        "user123",
				Email:     "invalid-email",
				Name:      "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid format for email, expected: valid email address",
		},
		{
			name: "empty name",
			user: &User{
				ID:        "user123",
				Email:     "test@example.com",
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errMsg:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	user := &User{
		ID:        "user123",
		Email:     "test@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	originalUpdatedAt := user.UpdatedAt
	newName := "New Name"

	err := user.UpdateName(context.Background(), newName)
	require.NoError(t, err)

	assert.Equal(t, newName, user.Name)
	// UpdateName no longer updates timestamp - that's done by repository
	assert.Equal(t, originalUpdatedAt, user.UpdatedAt)
}

func TestUser_UpdateName_Invalid(t *testing.T) {
	user := &User{
		ID:        "user123",
		Email:     "test@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		newName string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty name",
			newName: "",
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "valid name",
			newName: "Valid Name",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdateName(context.Background(), tt.newName)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, user.Name)
			}
		})
	}
}

func TestUser_IsEmailValid(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "invalid email - no @",
			email: "invalid-email",
			want:  false,
		},
		{
			name:  "invalid email - no domain",
			email: "test@",
			want:  false,
		},
		{
			name:  "invalid email - no local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "empty email",
			email: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Email: tt.email}
			got := user.IsEmailValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	// Initialize logger for testing
	logger.InitializeWithConfig(logger.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "stdout",
	})

	tests := []struct {
		name        string
		user        *User
		newEmail    string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid email update",
			user: &User{
				ID:    "user123",
				Email: "old@example.com",
				Name:  "Test User",
			},
			newEmail: "new@example.com",
			wantErr:  false,
		},
		{
			name: "empty email should fail",
			user: &User{
				ID:    "user123",
				Email: "old@example.com",
				Name:  "Test User",
			},
			newEmail:    "",
			wantErr:     true,
			expectedErr: "email is required",
		},
		{
			name: "invalid email format should fail",
			user: &User{
				ID:    "user123",
				Email: "old@example.com",
				Name:  "Test User",
			},
			newEmail:    "invalid-email",
			wantErr:     true,
			expectedErr: "invalid format for email",
		},
		{
			name: "email with special characters",
			user: &User{
				ID:    "user123",
				Email: "old@example.com",
				Name:  "Test User",
			},
			newEmail: "test+special@example.co.uk",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			oldEmail := tt.user.Email

			err := tt.user.UpdateEmail(ctx, tt.newEmail)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				// Email should not be changed on error
				assert.Equal(t, oldEmail, tt.user.Email)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newEmail, tt.user.Email)
			}
		})
	}
}

func TestUser_UpdateNameAdvanced(t *testing.T) {
	// Initialize logger for testing
	logger.InitializeWithConfig(logger.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "stdout",
	})

	tests := []struct {
		name        string
		user        *User
		newName     string
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid name update",
			user: &User{
				ID:    "user123",
				Email: "test@example.com",
				Name:  "Old Name",
			},
			newName: "New Name",
			wantErr: false,
		},
		{
			name: "empty name should fail",
			user: &User{
				ID:    "user123",
				Email: "test@example.com",
				Name:  "Old Name",
			},
			newName:     "",
			wantErr:     true,
			expectedErr: "name is required",
		},
		{
			name: "name with special characters",
			user: &User{
				ID:    "user123",
				Email: "test@example.com",
				Name:  "Old Name",
			},
			newName: "José María García-López",
			wantErr: false,
		},
		{
			name: "very long name",
			user: &User{
				ID:    "user123",
				Email: "test@example.com",
				Name:  "Old Name",
			},
			newName: "This is a very long name that might be used for testing purposes and should be handled properly",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			oldName := tt.user.Name

			err := tt.user.UpdateName(ctx, tt.newName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				// Name should not be changed on error
				assert.Equal(t, oldName, tt.user.Name)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.newName, tt.user.Name)
			}
		})
	}
}
