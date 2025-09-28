package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/cctw-zed/wonder/pkg/logger"
)

func TestNewUserRepositoryWithLogger_Validation(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

	// Mock database (nil for validation tests)
	var db *gorm.DB
	log := logger.Get().WithLayer("infrastructure").WithComponent("user_repository")

	t.Run("nil database panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserRepositoryWithLogger(nil, log)
		})
	})

	t.Run("nil logger panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserRepositoryWithLogger(db, nil)
		})
	})

	// Note: We can't test successful creation without a real database connection
	// That would be covered by integration tests
}
