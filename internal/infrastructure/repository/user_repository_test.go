package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/testutil/builder"
	"github.com/cctw-zed/wonder/pkg/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Initialize logger for tests
	logger.Initialize()

	// Skip integration tests if no test database is available
	dsn := "host=localhost port=5432 user=test password=test dbname=wonder_test sslmode=disable timezone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("No test database available, skipping integration tests")
		return nil
	}

	// Clean up any existing data
	db.Exec("DROP TABLE IF EXISTS users")

	// Auto-migrate the schema
	err = db.AutoMigrate(&user.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user creation",
			user: builder.NewUserBuilder().
				WithID("test-id-123").
				WithEmail("test@example.com").
				WithName("Test User").
				Build(),
			wantErr: false,
		},
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
			errMsg:  "database error in operation 'create' on table 'users'",
		},
		{
			name: "invalid user - empty ID",
			user: builder.NewUserBuilder().
				WithID("").
				WithEmail("test@example.com").
				WithName("Test User").
				Build(),
			wantErr: true,
			errMsg:  "validation failed for field",
		},
		{
			name: "invalid user - invalid email",
			user: builder.NewUserBuilder().
				WithID("test-id-456").
				WithEmail("invalid-email").
				WithName("Test User").
				Build(),
			wantErr: true,
			errMsg:  "validation failed for field",
		},
		{
			name: "duplicate email",
			user: builder.NewUserBuilder().
				WithID("test-id-duplicate").
				WithEmail("test@example.com"). // Same email as first test
				WithName("Duplicate User").
				Build(),
			wantErr: true,
			errMsg:  "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)

				// Verify user was created with timestamps
				assert.False(t, tt.user.CreatedAt.IsZero())
				assert.False(t, tt.user.UpdatedAt.IsZero())
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test user
	testUser := builder.NewUserBuilder().
		WithID("test-get-id").
		WithEmail("get@example.com").
		WithName("Get User").
		Build()

	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		exists  bool
		wantErr bool
		errMsg  string
	}{
		{
			name:    "existing user",
			id:      "test-get-id",
			exists:  true,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      "non-existing-id",
			exists:  false,
			wantErr: false,
		},
		{
			name:    "empty ID",
			id:      "",
			exists:  false,
			wantErr: true,
			errMsg:  "validation failed for field 'id': id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				if !tt.exists {
					require.Nil(t, user)
				} else {
					require.NotNil(t, user)
					assert.Equal(t, tt.id, user.ID)
					assert.Equal(t, testUser.Email, user.Email)
					assert.Equal(t, testUser.Name, user.Name)
				}
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test user
	testUser := builder.NewUserBuilder().
		WithID("test-get-email").
		WithEmail("email@example.com").
		WithName("Email User").
		Build()

	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		email   string
		exists  bool
		wantErr bool
		errMsg  string
	}{
		{
			name:    "existing email",
			email:   "email@example.com",
			exists:  true,
			wantErr: false,
		},
		{
			name:    "non-existing email",
			email:   "notfound@example.com",
			exists:  false,
			wantErr: false,
		},
		{
			name:    "empty email",
			email:   "",
			exists:  false,
			wantErr: true,
			errMsg:  "validation failed for field 'email': email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByEmail(ctx, tt.email)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				if !tt.exists {
					require.Nil(t, user)
				} else {
					require.NotNil(t, user)
					assert.Equal(t, testUser.ID, user.ID)
					assert.Equal(t, tt.email, user.Email)
					assert.Equal(t, testUser.Name, user.Name)
				}
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test user
	testUser := builder.NewUserBuilder().
		WithID("test-update-id").
		WithEmail("update@example.com").
		WithName("Original Name").
		Build()

	err := repo.Create(ctx, testUser)
	require.NoError(t, err)

	originalUpdatedAt := testUser.UpdatedAt

	tests := []struct {
		name    string
		setupFn func() *user.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update",
			setupFn: func() *user.User {
				u := *testUser // Copy
				u.Name = "Updated Name"
				return &u
			},
			wantErr: false,
		},
		{
			name: "nil user",
			setupFn: func() *user.User {
				return nil
			},
			wantErr: true,
			errMsg:  "user cannot be nil",
		},
		{
			name: "empty ID",
			setupFn: func() *user.User {
				u := *testUser
				u.ID = ""
				return &u
			},
			wantErr: true,
			errMsg:  "user ID cannot be empty",
		},
		{
			name: "non-existing user",
			setupFn: func() *user.User {
				return builder.NewUserBuilder().
					WithID("non-existing-id").
					WithEmail("nonexist@example.com").
					WithName("Non Existing").
					Build()
			},
			wantErr: false,
		},
		{
			name: "invalid user data",
			setupFn: func() *user.User {
				u := *testUser
				u.Email = "invalid-email"
				return &u
			},
			wantErr: true,
			errMsg:  "validation failed for field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := tt.setupFn()
			err := repo.Update(ctx, user)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)

				// Verify UpdatedAt was updated
				assert.True(t, user.UpdatedAt.After(originalUpdatedAt))

				// Verify user was actually updated in database
				updated, err := repo.GetByID(ctx, user.ID)
				require.NoError(t, err)
				assert.Equal(t, user.Name, updated.Name)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	testUser1 := builder.NewUserBuilder().
		WithID("test-delete-1").
		WithEmail("delete1@example.com").
		WithName("Delete User 1").
		Build()

	testUser2 := builder.NewUserBuilder().
		WithID("test-delete-2").
		WithEmail("delete2@example.com").
		WithName("Delete User 2").
		Build()

	err := repo.Create(ctx, testUser1)
	require.NoError(t, err)
	err = repo.Create(ctx, testUser2)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "existing user",
			id:      "test-delete-1",
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      "non-existing-id",
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
			errMsg:  "user ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)

				// Verify user was actually deleted
				deleted, err := repo.GetByID(ctx, tt.id)
				assert.NoError(t, err)
				assert.Nil(t, deleted)
			}
		})
	}
}

func TestUserRepository_ConcurrentOperations(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Test concurrent creates with same email (should fail)
	done := make(chan bool, 2)
	errors := make(chan error, 2)

	createUser := func(id, suffix string) {
		user := builder.NewUserBuilder().
			WithID(id).
			WithEmail("concurrent@example.com"). // Same email
			WithName("Concurrent User " + suffix).
			Build()

		err := repo.Create(ctx, user)
		errors <- err
		done <- true
	}

	go createUser("concurrent-1", "1")
	go createUser("concurrent-2", "2")

	// Wait for both operations
	<-done
	<-done
	close(errors)

	// One should succeed, one should fail
	var successCount, errorCount int
	for err := range errors {
		if err != nil {
			errorCount++
			assert.Contains(t, err.Error(), "already exists")
		} else {
			successCount++
		}
	}

	assert.Equal(t, 1, successCount, "exactly one create should succeed")
	assert.Equal(t, 1, errorCount, "exactly one create should fail")
}

func TestUserRepository_TimestampHandling(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Test automatic timestamp setting
	user := builder.NewUserBuilder().
		WithID("timestamp-test").
		WithEmail("timestamp@example.com").
		WithName("Timestamp User").
		Build()

	// Clear timestamps to test auto-setting
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

	before := time.Now()
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	after := time.Now()

	// Verify timestamps were set
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.True(t, user.CreatedAt.After(before) || user.CreatedAt.Equal(before))
	assert.True(t, user.CreatedAt.Before(after) || user.CreatedAt.Equal(after))

	// Test update timestamp
	originalUpdated := user.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference

	user.Name = "Updated Name"
	err = repo.Update(ctx, user)
	require.NoError(t, err)

	assert.True(t, user.UpdatedAt.After(originalUpdated))
}
