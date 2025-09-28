package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/cctw-zed/wonder/internal/domain/user/mocks"
	"github.com/cctw-zed/wonder/pkg/logger"
	idMocks "github.com/cctw-zed/wonder/pkg/snowflake/id/mocks"
)

func TestNewUserServiceWithLogger_Validation(t *testing.T) {
	// Initialize logger for tests
	logger.Initialize()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockIDGen := idMocks.NewMockGenerator(ctrl)
	log := logger.Get().WithLayer("application").WithComponent("user_service")

	t.Run("successful creation", func(t *testing.T) {
		service := NewUserServiceWithLogger(mockRepo, mockIDGen, log)
		assert.NotNil(t, service)
	})

	t.Run("nil repository panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserServiceWithLogger(nil, mockIDGen, log)
		})
	})

	t.Run("nil ID generator panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserServiceWithLogger(mockRepo, nil, log)
		})
	})

	t.Run("nil logger panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewUserServiceWithLogger(mockRepo, mockIDGen, nil)
		})
	})

	t.Run("standard constructor works", func(t *testing.T) {
		service := NewUserService(mockRepo, mockIDGen)
		assert.NotNil(t, service)
	})
}
