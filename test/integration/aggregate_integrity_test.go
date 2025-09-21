package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/infrastructure/repository"
	"github.com/cctw-zed/wonder/internal/testutil/builder"
	idMocks "github.com/cctw-zed/wonder/pkg/snowflake/id/mocks"

	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestUserAggregateIntegrity verifies that User aggregate maintains business invariants
// and data consistency across all operations
func TestUserAggregateIntegrity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock ID generator
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	mockIDGen.EXPECT().Generate().Return("integrity-test-id").AnyTimes()

	// Skip if no test database
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}

	// Setup repository and service
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo, mockIDGen)

	ctx := context.Background()

	t.Run("User aggregate enforces business invariants", func(t *testing.T) {
		// Test 1: Email uniqueness constraint
		t.Run("email uniqueness constraint", func(t *testing.T) {
			email := "unique@test.com"

			// First user creation should succeed
			user1, err := userService.Register(ctx, email, "First User")
			require.NoError(t, err)
			assert.Equal(t, email, user1.Email)

			// Second user with same email should fail
			_, err = userService.Register(ctx, email, "Second User")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "already exists")
		})

		// Test 2: Entity validation consistency
		t.Run("entity validation consistency", func(t *testing.T) {
			// Create user with valid data
			validUser := builder.NewUserBuilder().
				WithID("validation-test").
				WithEmail("validation@test.com").
				WithName("Valid User").
				Build()

			err := repo.Create(ctx, validUser)
			require.NoError(t, err)

			// Retrieve and verify data integrity
			retrieved, err := repo.GetByID(ctx, "validation-test")
			require.NoError(t, err)

			assert.Equal(t, validUser.ID, retrieved.ID)
			assert.Equal(t, validUser.Email, retrieved.Email)
			assert.Equal(t, validUser.Name, retrieved.Name)
			assert.False(t, retrieved.CreatedAt.IsZero())
			assert.False(t, retrieved.UpdatedAt.IsZero())

			// Verify entity still validates after retrieval
			err = retrieved.Validate()
			require.NoError(t, err)
		})

		// Test 3: Update operations maintain aggregate consistency
		t.Run("update operations maintain consistency", func(t *testing.T) {
			// Create initial user
			initialUser := builder.NewUserBuilder().
				WithID("update-test").
				WithEmail("update@test.com").
				WithName("Original Name").
				Build()

			err := repo.Create(ctx, initialUser)
			require.NoError(t, err)

			originalCreatedAt := initialUser.CreatedAt
			originalUpdatedAt := initialUser.UpdatedAt

			// Wait to ensure timestamp difference
			time.Sleep(time.Millisecond)

			// Update user name
			err = initialUser.UpdateName("Updated Name")
			require.NoError(t, err)

			// Save changes
			err = repo.Update(ctx, initialUser)
			require.NoError(t, err)

			// Verify aggregate consistency after update
			updated, err := repo.GetByID(ctx, "update-test")
			require.NoError(t, err)

			assert.Equal(t, "Updated Name", updated.Name)
			assert.Equal(t, "update@test.com", updated.Email) // Email unchanged
			// Compare timestamps by truncating to microseconds to handle database precision
			assert.True(t, originalCreatedAt.Truncate(time.Microsecond).Equal(updated.CreatedAt.Truncate(time.Microsecond))) // CreatedAt unchanged
			assert.True(t, updated.UpdatedAt.After(originalUpdatedAt.Truncate(time.Microsecond))) // UpdatedAt changed

			// Verify entity still validates after update
			err = updated.Validate()
			require.NoError(t, err)
		})

		// Test 4: Aggregate boundary enforcement
		t.Run("aggregate boundary enforcement", func(t *testing.T) {
			// Verify repository only operates on complete aggregates
			user := builder.NewUserBuilder().
				WithID("boundary-test").
				WithEmail("boundary@test.com").
				WithName("Boundary User").
				Build()

			// Create user
			err := repo.Create(ctx, user)
			require.NoError(t, err)

			// Verify we can only retrieve complete aggregates
			retrieved, err := repo.GetByID(ctx, "boundary-test")
			require.NoError(t, err)

			// All aggregate fields should be populated
			assert.NotEmpty(t, retrieved.ID)
			assert.NotEmpty(t, retrieved.Email)
			assert.NotEmpty(t, retrieved.Name)
			assert.False(t, retrieved.CreatedAt.IsZero())
			assert.False(t, retrieved.UpdatedAt.IsZero())

			// Delete should remove entire aggregate
			err = repo.Delete(ctx, "boundary-test")
			require.NoError(t, err)

			// Verify aggregate no longer exists
			deletedUser, err := repo.GetByID(ctx, "boundary-test")
			require.NoError(t, err)
			assert.Nil(t, deletedUser, "user should not exist after deletion")
		})

		// Test 5: Transaction consistency (simulated)
		t.Run("transaction consistency", func(t *testing.T) {
			// Test that partial failures don't leave system in inconsistent state
			invalidUser := builder.NewUserBuilder().
				WithID("transaction-test").
				WithEmail("invalid-email-format").
				WithName("Transaction User").
				Build()

			// Attempt to create invalid user should fail
			err := repo.Create(ctx, invalidUser)
			require.Error(t, err)

			// Verify no partial data was persisted
			user, err := repo.GetByID(ctx, "transaction-test")
			require.NoError(t, err)
			assert.Nil(t, user, "user should not exist after failed creation")

			user, err = repo.GetByEmail(ctx, "invalid-email-format")
			require.NoError(t, err)
			assert.Nil(t, user, "user should not exist after failed creation")
		})
	})

	t.Run("Repository adheres to DDD patterns", func(t *testing.T) {
		// Test 6: Repository only exposes aggregate operations
		t.Run("only aggregate operations exposed", func(t *testing.T) {
			// Verify repository interface only has aggregate-level operations
			// This is validated by the interface design itself

			user := builder.NewUserBuilder().
				WithID("aggregate-ops").
				WithEmail("aggregateops@test.com").
				WithName("Aggregate User").
				Build()

			// All operations work on complete User aggregate
			err := repo.Create(ctx, user)
			require.NoError(t, err)

			retrieved, err := repo.GetByID(ctx, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user.ID, retrieved.ID)

			retrieved, err = repo.GetByEmail(ctx, user.Email)
			require.NoError(t, err)
			assert.Equal(t, user.Email, retrieved.Email)

			user.Name = "Updated Aggregate User"
			err = repo.Update(ctx, user)
			require.NoError(t, err)

			err = repo.Delete(ctx, user.ID)
			require.NoError(t, err)
		})

		// Test 7: Domain events (placeholder for future implementation)
		t.Run("domain events placeholder", func(t *testing.T) {
			// This test serves as a placeholder for when domain events are implemented
			// Domain events should be raised when aggregate state changes
			t.Skip("Domain events not yet implemented")
		})
	})
}

func setupIntegrationTestDB(t *testing.T) *gorm.DB {
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