package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cctw-zed/wonder/internal/container"
	"github.com/cctw-zed/wonder/internal/domain/user"
	"github.com/cctw-zed/wonder/internal/server"
)

// E2ETestSuite represents an end-to-end test suite that starts a real server
type E2ETestSuite struct {
	container    *container.Container
	server       *server.Server
	httpServer   *http.Server
	baseURL      string
	httpClient   *http.Client
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewE2ETestSuite creates a new E2E test suite
func NewE2ETestSuite(t *testing.T) *E2ETestSuite {
	ctx, cancel := context.WithCancel(context.Background())

	// Create container with testing configuration
	c, err := container.NewContainerForEnvironment(ctx, "testing")
	require.NoError(t, err, "Failed to create container for testing")

	// Get server configuration
	port := c.Config.Server.Port
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	suite := &E2ETestSuite{
		container: c,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	return suite
}

// StartServer starts the HTTP server for testing
func (s *E2ETestSuite) StartServer(t *testing.T) {
	// Create server instance
	s.server = server.New(s.container)

	// Start server in a goroutine
	go func() {
		if err := s.server.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	// Wait for server to be ready
	s.waitForServerReady(t)
}

// CleanupDatabase cleans up test data from database
func (s *E2ETestSuite) CleanupDatabase(t *testing.T) {
	// Clean users table using GORM
	err := s.container.Database.DB().Where("email LIKE ?", "%test.com").Delete(&user.User{}).Error
	if err != nil {
		t.Logf("Warning: Failed to cleanup test data: %v", err)
	}
}

// Cleanup cleans up resources after tests
func (s *E2ETestSuite) Cleanup() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
	if s.cancel != nil {
		s.cancel()
	}
	if s.container != nil {
		s.container.Close()
	}
}

// waitForServerReady waits for the server to be ready to accept requests
func (s *E2ETestSuite) waitForServerReady(t *testing.T) {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for server to be ready")
		case <-ticker.C:
			resp, err := s.httpClient.Get(s.baseURL + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				t.Log("Server is ready")
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

// TestServerE2E demonstrates a complete end-to-end test
func TestServerE2E(t *testing.T) {
	// Skip E2E tests in short mode
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	suite := NewE2ETestSuite(t)
	defer suite.Cleanup()

	// Clean database before tests
	suite.CleanupDatabase(t)

	suite.StartServer(t)

	t.Run("Health Check E2E", func(t *testing.T) {
		resp, err := suite.httpClient.Get(suite.baseURL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)

		assert.Equal(t, "wonder", health["app"])
		assert.Equal(t, "testing", health["environment"])
		assert.Equal(t, "healthy", health["status"])
	})

	t.Run("User Registration E2E Flow", func(t *testing.T) {
		// Test complete user registration flow
		requestBody := map[string]string{
			"email": "e2e_new@test.com",
			"name":  "E2E Test User",
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		resp, err := suite.httpClient.Post(
			suite.baseURL+"/api/v1/users/register",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)

		assert.NotEmpty(t, user["id"])
		assert.Equal(t, "e2e_new@test.com", user["email"])
		assert.Equal(t, "E2E Test User", user["name"])
		assert.NotEmpty(t, user["created_at"])
		assert.NotEmpty(t, user["updated_at"])
	})

	t.Run("Duplicate Email E2E", func(t *testing.T) {
		// First create a user
		firstRequestBody := map[string]string{
			"email": "e2e_duplicate@test.com",
			"name":  "First User",
		}

		firstJsonBody, err := json.Marshal(firstRequestBody)
		require.NoError(t, err)

		firstResp, err := suite.httpClient.Post(
			suite.baseURL+"/api/v1/users/register",
			"application/json",
			bytes.NewBuffer(firstJsonBody),
		)
		require.NoError(t, err)
		firstResp.Body.Close()
		assert.Equal(t, http.StatusCreated, firstResp.StatusCode)

		// Now test duplicate email handling
		requestBody := map[string]string{
			"email": "e2e_duplicate@test.com", // Same email as above
			"name":  "Another User",
		}

		jsonBody, err := json.Marshal(requestBody)
		require.NoError(t, err)

		resp, err := suite.httpClient.Post(
			suite.baseURL+"/api/v1/users/register",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // Current behavior

		var errorResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.Contains(t, errorResp["error"], "already exists")
	})
}

// BenchmarkServerE2E benchmarks the server performance
func BenchmarkServerE2E(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping E2E benchmarks in short mode")
	}

	suite := NewE2ETestSuite(&testing.T{})
	defer suite.Cleanup()

	// suite.StartServer(&testing.T{})

	b.ResetTimer()

	b.Run("Health Check", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := suite.httpClient.Get(suite.baseURL + "/health")
			if err == nil {
				resp.Body.Close()
			}
		}
	})
}