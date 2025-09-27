package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/pkg/jwt"
)

func TestNewAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a simple mock that implements the AuthService interface
	mockAuthService := &mockAuthService{}
	handler := NewAuthHandler(mockAuthService)

	require.NotNil(t, handler)
}

// Simple mock implementation for testing constructor only
type mockAuthService struct{}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*service.LoginResponse, error) {
	return nil, nil
}

func (m *mockAuthService) Logout(ctx context.Context, token string) error {
	return nil
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	return nil, nil
}
