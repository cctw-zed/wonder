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
	container  *container.Container
	server     *server.Server
	httpServer *http.Server
	baseURL    string
	httpClient *http.Client
	ctx        context.Context
	cancel     context.CancelFunc
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

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		user := response["user"].(map[string]interface{})
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

		assert.Equal(t, http.StatusConflict, resp.StatusCode) // Correct behavior - 409 for conflicts

		var errorResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		// The error could be in different fields depending on the error type
		var errorMessage string
		if errorResp["error"] != nil {
			errorMessage = errorResp["error"].(string)
		} else if errorResp["message"] != nil {
			errorMessage = errorResp["message"].(string)
		} else if errorResp["details"] != nil {
			if details, ok := errorResp["details"].(map[string]interface{}); ok {
				if msg, ok := details["message"].(string); ok {
					errorMessage = msg
				}
			}
		}
		require.NotEmpty(t, errorMessage, "Expected error message in response")
		assert.Contains(t, errorMessage, "conflict")
	})

	t.Run("Complete User Lifecycle E2E Flow", func(t *testing.T) {
		// Step 1: Create a new user
		createRequestBody := map[string]string{
			"email": "e2e_lifecycle@test.com",
			"name":  "E2E Lifecycle User",
		}

		createJsonBody, err := json.Marshal(createRequestBody)
		require.NoError(t, err)

		createResp, err := suite.httpClient.Post(
			suite.baseURL+"/api/v1/users/register",
			"application/json",
			bytes.NewBuffer(createJsonBody),
		)
		require.NoError(t, err)
		defer createResp.Body.Close()

		assert.Equal(t, http.StatusCreated, createResp.StatusCode)

		var createResponse map[string]interface{}
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		require.NoError(t, err)

		createdUser := createResponse["user"].(map[string]interface{})
		userID := createdUser["id"].(string)
		assert.NotEmpty(t, userID)

		// Step 2: Get user profile
		getResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users/" + userID)
		require.NoError(t, err)
		defer getResp.Body.Close()

		assert.Equal(t, http.StatusOK, getResp.StatusCode)

		var getResponse map[string]interface{}
		err = json.NewDecoder(getResp.Body).Decode(&getResponse)
		require.NoError(t, err)

		retrievedUser := getResponse["user"].(map[string]interface{})
		assert.Equal(t, userID, retrievedUser["id"])
		assert.Equal(t, "e2e_lifecycle@test.com", retrievedUser["email"])
		assert.Equal(t, "E2E Lifecycle User", retrievedUser["name"])

		// Step 3: Update user profile - name only
		updateNameBody := map[string]string{
			"name": "Updated Lifecycle User",
		}

		updateNameJsonBody, err := json.Marshal(updateNameBody)
		require.NoError(t, err)

		updateNameReq, err := http.NewRequest(
			http.MethodPut,
			suite.baseURL+"/api/v1/users/"+userID,
			bytes.NewBuffer(updateNameJsonBody),
		)
		require.NoError(t, err)
		updateNameReq.Header.Set("Content-Type", "application/json")

		updateNameResp, err := suite.httpClient.Do(updateNameReq)
		require.NoError(t, err)
		defer updateNameResp.Body.Close()

		assert.Equal(t, http.StatusOK, updateNameResp.StatusCode)

		var updateNameResponse map[string]interface{}
		err = json.NewDecoder(updateNameResp.Body).Decode(&updateNameResponse)
		require.NoError(t, err)

		updatedUser := updateNameResponse["user"].(map[string]interface{})
		assert.Equal(t, userID, updatedUser["id"])
		assert.Equal(t, "e2e_lifecycle@test.com", updatedUser["email"])
		assert.Equal(t, "Updated Lifecycle User", updatedUser["name"])

		// Step 4: Update user profile - email only
		updateEmailBody := map[string]string{
			"email": "e2e_updated@test.com",
		}

		updateEmailJsonBody, err := json.Marshal(updateEmailBody)
		require.NoError(t, err)

		updateEmailReq, err := http.NewRequest(
			http.MethodPut,
			suite.baseURL+"/api/v1/users/"+userID,
			bytes.NewBuffer(updateEmailJsonBody),
		)
		require.NoError(t, err)
		updateEmailReq.Header.Set("Content-Type", "application/json")

		updateEmailResp, err := suite.httpClient.Do(updateEmailReq)
		require.NoError(t, err)
		defer updateEmailResp.Body.Close()

		assert.Equal(t, http.StatusOK, updateEmailResp.StatusCode)

		var updateEmailResponse map[string]interface{}
		err = json.NewDecoder(updateEmailResp.Body).Decode(&updateEmailResponse)
		require.NoError(t, err)

		updatedUserWithEmail := updateEmailResponse["user"].(map[string]interface{})
		assert.Equal(t, userID, updatedUserWithEmail["id"])
		assert.Equal(t, "e2e_updated@test.com", updatedUserWithEmail["email"])
		assert.Equal(t, "Updated Lifecycle User", updatedUserWithEmail["name"])

		// Step 5: List users to verify the updated user exists
		listResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users?page=1&page_size=10")
		require.NoError(t, err)
		defer listResp.Body.Close()

		assert.Equal(t, http.StatusOK, listResp.StatusCode)

		var listResponse map[string]interface{}
		err = json.NewDecoder(listResp.Body).Decode(&listResponse)
		require.NoError(t, err)

		users := listResponse["users"].([]interface{})
		assert.GreaterOrEqual(t, len(users), 1)

		// Find our updated user in the list
		var foundUser map[string]interface{}
		for _, u := range users {
			userData := u.(map[string]interface{})
			if userData["id"] == userID {
				foundUser = userData
				break
			}
		}
		require.NotNil(t, foundUser, "Updated user should be found in users list")
		assert.Equal(t, "e2e_updated@test.com", foundUser["email"])
		assert.Equal(t, "Updated Lifecycle User", foundUser["name"])

		// Step 6: Delete the user
		deleteReq, err := http.NewRequest(
			http.MethodDelete,
			suite.baseURL+"/api/v1/users/"+userID,
			nil,
		)
		require.NoError(t, err)

		deleteResp, err := suite.httpClient.Do(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()

		assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

		// Step 7: Verify user is deleted - should return 404
		verifyDeleteResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users/" + userID)
		require.NoError(t, err)
		defer verifyDeleteResp.Body.Close()

		assert.Equal(t, http.StatusNotFound, verifyDeleteResp.StatusCode)
	})

	t.Run("User List with Pagination and Filters E2E", func(t *testing.T) {
		// Create multiple test users for pagination testing
		testUsers := []map[string]string{
			{"email": "e2e_pagination1@test.com", "name": "Pagination User 1"},
			{"email": "e2e_pagination2@test.com", "name": "Pagination User 2"},
			{"email": "e2e_pagination3@test.com", "name": "Pagination User 3"},
		}

		userIDs := make([]string, 0, len(testUsers))
		defer func() {
			// Clean up created users
			for _, userID := range userIDs {
				deleteReq, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/users/"+userID, nil)
				if deleteReq != nil {
					suite.httpClient.Do(deleteReq)
				}
			}
		}()

		// Create test users
		for _, testUser := range testUsers {
			jsonBody, err := json.Marshal(testUser)
			require.NoError(t, err)

			resp, err := suite.httpClient.Post(
				suite.baseURL+"/api/v1/users/register",
				"application/json",
				bytes.NewBuffer(jsonBody),
			)
			require.NoError(t, err)

			var createResponse map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&createResponse)
			resp.Body.Close()
			require.NoError(t, err)

			createdUser := createResponse["user"].(map[string]interface{})
			userIDs = append(userIDs, createdUser["id"].(string))
		}

		// Test pagination - page 1, size 2
		listResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users?page=1&page_size=2")
		require.NoError(t, err)
		defer listResp.Body.Close()

		assert.Equal(t, http.StatusOK, listResp.StatusCode)

		var listResponse map[string]interface{}
		err = json.NewDecoder(listResp.Body).Decode(&listResponse)
		require.NoError(t, err)

		users := listResponse["users"].([]interface{})
		assert.LessOrEqual(t, len(users), 2) // Should return at most 2 users

		meta := listResponse["meta"].(map[string]interface{})
		assert.Equal(t, float64(1), meta["page"])
		assert.Equal(t, float64(2), meta["page_size"])
		assert.GreaterOrEqual(t, int(meta["total"].(float64)), 3) // Should have at least our 3 test users

		// Test with name filter
		filterResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users?name=Pagination&page=1&page_size=10")
		require.NoError(t, err)
		defer filterResp.Body.Close()

		assert.Equal(t, http.StatusOK, filterResp.StatusCode)

		var filterResponse map[string]interface{}
		err = json.NewDecoder(filterResp.Body).Decode(&filterResponse)
		require.NoError(t, err)

		filteredUsers := filterResponse["users"].([]interface{})
		assert.GreaterOrEqual(t, len(filteredUsers), 3) // Should find our pagination test users

		// Verify all returned users have "Pagination" in their name
		for _, u := range filteredUsers {
			userData := u.(map[string]interface{})
			userName := userData["name"].(string)
			assert.Contains(t, userName, "Pagination")
		}
	})

	t.Run("Error Handling E2E", func(t *testing.T) {
		// Test 1: Get non-existent user
		getResp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users/nonexistent-id")
		require.NoError(t, err)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)

		// Test 2: Update non-existent user
		updateBody := map[string]string{
			"name": "Updated Name",
		}
		updateJsonBody, err := json.Marshal(updateBody)
		require.NoError(t, err)

		updateReq, err := http.NewRequest(
			http.MethodPut,
			suite.baseURL+"/api/v1/users/nonexistent-id",
			bytes.NewBuffer(updateJsonBody),
		)
		require.NoError(t, err)
		updateReq.Header.Set("Content-Type", "application/json")

		updateResp, err := suite.httpClient.Do(updateReq)
		require.NoError(t, err)
		defer updateResp.Body.Close()
		assert.Equal(t, http.StatusNotFound, updateResp.StatusCode)

		// Test 3: Delete non-existent user
		deleteReq, err := http.NewRequest(
			http.MethodDelete,
			suite.baseURL+"/api/v1/users/nonexistent-id",
			nil,
		)
		require.NoError(t, err)

		deleteResp, err := suite.httpClient.Do(deleteReq)
		require.NoError(t, err)
		defer deleteResp.Body.Close()
		assert.Equal(t, http.StatusNotFound, deleteResp.StatusCode)

		// Test 4: Invalid update data
		invalidUpdateBody := map[string]string{
			"email": "invalid-email-format",
		}
		invalidJsonBody, err := json.Marshal(invalidUpdateBody)
		require.NoError(t, err)

		// First create a user to test invalid update
		createBody := map[string]string{
			"email": "e2e_invalid_update@test.com",
			"name":  "Test User for Invalid Update",
		}
		createJsonBody, err := json.Marshal(createBody)
		require.NoError(t, err)

		createResp, err := suite.httpClient.Post(
			suite.baseURL+"/api/v1/users/register",
			"application/json",
			bytes.NewBuffer(createJsonBody),
		)
		require.NoError(t, err)

		var createResponse map[string]interface{}
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		createResp.Body.Close()
		require.NoError(t, err)

		createdUser := createResponse["user"].(map[string]interface{})
		userID := createdUser["id"].(string)

		defer func() {
			// Clean up
			deleteReq, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/users/"+userID, nil)
			if deleteReq != nil {
				suite.httpClient.Do(deleteReq)
			}
		}()

		// Try to update with invalid email
		invalidUpdateReq, err := http.NewRequest(
			http.MethodPut,
			suite.baseURL+"/api/v1/users/"+userID,
			bytes.NewBuffer(invalidJsonBody),
		)
		require.NoError(t, err)
		invalidUpdateReq.Header.Set("Content-Type", "application/json")

		invalidUpdateResp, err := suite.httpClient.Do(invalidUpdateReq)
		require.NoError(t, err)
		defer invalidUpdateResp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, invalidUpdateResp.StatusCode)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(invalidUpdateResp.Body).Decode(&errorResponse)
		require.NoError(t, err)

		// Should contain validation error message
		assert.NotEmpty(t, errorResponse)
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

// BenchmarkLifecycleAPIsE2E benchmarks the new lifecycle APIs performance
func BenchmarkLifecycleAPIsE2E(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping E2E benchmarks in short mode")
	}

	suite := NewE2ETestSuite(&testing.T{})
	defer suite.Cleanup()

	suite.CleanupDatabase(&testing.T{})
	// suite.StartServer(&testing.T{})

	// Create a test user for benchmarking get/update/delete operations
	createBody := map[string]string{
		"email": "benchmark@test.com",
		"name":  "Benchmark User",
	}
	createJsonBody, _ := json.Marshal(createBody)

	createResp, err := suite.httpClient.Post(
		suite.baseURL+"/api/v1/users/register",
		"application/json",
		bytes.NewBuffer(createJsonBody),
	)
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	var createResponse map[string]interface{}
	json.NewDecoder(createResp.Body).Decode(&createResponse)
	createResp.Body.Close()

	createdUser := createResponse["user"].(map[string]interface{})
	userID := createdUser["id"].(string)

	defer func() {
		// Clean up test user
		deleteReq, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/users/"+userID, nil)
		if deleteReq != nil {
			suite.httpClient.Do(deleteReq)
		}
	}()

	b.ResetTimer()

	b.Run("Get User Profile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users/" + userID)
			if err == nil {
				resp.Body.Close()
			}
		}
	})

	b.Run("Update User Profile", func(b *testing.B) {
		updateBody := map[string]string{
			"name": "Updated Benchmark User",
		}
		updateJsonBody, _ := json.Marshal(updateBody)

		for i := 0; i < b.N; i++ {
			updateReq, _ := http.NewRequest(
				http.MethodPut,
				suite.baseURL+"/api/v1/users/"+userID,
				bytes.NewBuffer(updateJsonBody),
			)
			updateReq.Header.Set("Content-Type", "application/json")

			resp, err := suite.httpClient.Do(updateReq)
			if err == nil {
				resp.Body.Close()
			}
		}
	})

	b.Run("List Users", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := suite.httpClient.Get(suite.baseURL + "/api/v1/users?page=1&page_size=10")
			if err == nil {
				resp.Body.Close()
			}
		}
	})
}
