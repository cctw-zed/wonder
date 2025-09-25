package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/cctw-zed/wonder/internal/domain/user"
)

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new database migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// MigrateAll runs all database migrations
func (m *Migrator) MigrateAll() error {
	if err := m.migrateUserTable(); err != nil {
		return fmt.Errorf("failed to migrate user table: %w", err)
	}

	return nil
}

// migrateUserTable creates or updates the users table
func (m *Migrator) migrateUserTable() error {
	// Check if the table exists
	hasTable := m.db.Migrator().HasTable(&user.User{})

	if hasTable {
		// Check if password_hash column exists
		hasPasswordHash := m.db.Migrator().HasColumn(&user.User{}, "password_hash")

		if !hasPasswordHash {
			// First, add the column as nullable
			if err := m.db.Exec("ALTER TABLE users ADD COLUMN password_hash VARCHAR(255) NULL").Error; err != nil {
				return fmt.Errorf("failed to add password_hash column: %w", err)
			}

			// Set a default password hash for existing users (they'll need to reset their password)
			defaultPasswordHash := "$2a$10$CwTycUXWue0Thq9StjUM0uNzZCbF8FFdQLJkN8XWNq3xKWbkjT2Ba" // "changeme"
			if err := m.db.Exec("UPDATE users SET password_hash = ? WHERE password_hash IS NULL", defaultPasswordHash).Error; err != nil {
				return fmt.Errorf("failed to set default password hash: %w", err)
			}

			// Now make it NOT NULL
			if err := m.db.Exec("ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL").Error; err != nil {
				return fmt.Errorf("failed to make password_hash NOT NULL: %w", err)
			}
		}
	}

	// Run the normal auto-migration
	if err := m.db.AutoMigrate(&user.User{}); err != nil {
		return fmt.Errorf("failed to auto-migrate User model: %w", err)
	}

	return nil
}

// DropAll drops all tables (use with caution!)
func (m *Migrator) DropAll() error {
	if err := m.db.Migrator().DropTable(&user.User{}); err != nil {
		return fmt.Errorf("failed to drop user table: %w", err)
	}

	return nil
}

// CreateIndexes creates additional database indexes
func (m *Migrator) CreateIndexes() error {
	// Email index is handled by uniqueIndex tag in the model, skip manual creation

	// Create created_at index for time-based queries
	if !m.db.Migrator().HasIndex(&user.User{}, "created_at") {
		if err := m.db.Migrator().CreateIndex(&user.User{}, "created_at"); err != nil {
			return fmt.Errorf("failed to create created_at index: %w", err)
		}
	}

	return nil
}

// CheckTables verifies that all required tables exist
func (m *Migrator) CheckTables() error {
	if !m.db.Migrator().HasTable(&user.User{}) {
		return fmt.Errorf("users table does not exist")
	}

	return nil
}

// GetTableStatus returns information about table structure
func (m *Migrator) GetTableStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// Check if users table exists
	status["users_table_exists"] = m.db.Migrator().HasTable(&user.User{})

	// Check columns
	userColumns := []string{"id", "email", "name", "created_at", "updated_at"}
	existingColumns := make(map[string]bool)

	for _, column := range userColumns {
		existingColumns[column] = m.db.Migrator().HasColumn(&user.User{}, column)
	}
	status["users_columns"] = existingColumns

	// Check indexes
	userIndexes := []string{"email", "created_at"}
	existingIndexes := make(map[string]bool)

	for _, index := range userIndexes {
		existingIndexes[index] = m.db.Migrator().HasIndex(&user.User{}, index)
	}
	status["users_indexes"] = existingIndexes

	return status, nil
}
