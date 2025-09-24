# Wonder Testing Guide

This document provides comprehensive guidance on testing strategies, test execution, and the testing framework for the Wonder project.

## Testing Architecture Overview

The Wonder project implements a multi-layered testing strategy covering:

1. **Unit Tests**: Fast, isolated tests for individual components
2. **Integration Tests**: Tests for component interactions and database operations
3. **End-to-End (E2E) Tests**: Full application stack testing with real HTTP server

## Test Runner Script

The project includes a comprehensive test runner script at `scripts/test.sh` that supports different testing scenarios.

### Basic Usage

```bash
# Make script executable (first time only)
chmod +x scripts/test.sh

# Run different types of tests
./scripts/test.sh unit           # Unit tests only
./scripts/test.sh integration    # Integration tests only
./scripts/test.sh e2e           # End-to-end tests only
./scripts/test.sh all           # All tests (unit + integration + e2e)
./scripts/test.sh coverage     # Generate coverage report
```

### Advanced Options

```bash
# Verbose output
./scripts/test.sh unit -v
./scripts/test.sh integration --verbose

# Short mode (skips long-running tests)
./scripts/test.sh all -s
./scripts/test.sh e2e --short

# Race condition detection
./scripts/test.sh unit --race
./scripts/test.sh all --race

# Clean test cache before running
./scripts/test.sh e2e --clean
./scripts/test.sh coverage --clean
```

### Help and Usage

```bash
# Show help message
./scripts/test.sh --help
./scripts/test.sh -h
```

## Test Categories

### 1. Unit Tests

**Location**: `internal/*/` directories
**Purpose**: Test individual components in isolation
**Dependencies**: Mock objects and test utilities

```bash
# Run unit tests
./scripts/test.sh unit

# Run with verbose output
./scripts/test.sh unit -v

# Run specific package
go test ./internal/domain/user -v
```

**Coverage**: Tests cover domain entities, application services, HTTP handlers, and utility functions.

### 2. Integration Tests

**Location**: `test/integration/` directory
**Purpose**: Test component interactions and database operations
**Dependencies**: Test database (uses same credentials as E2E tests)

```bash
# Run integration tests
./scripts/test.sh integration

# Set up test database environment
export DB_USERNAME="test"
export DB_PASSWORD="test"
export DB_DATABASE="wonder_test"
./scripts/test.sh integration
```

### 3. End-to-End (E2E) Tests

**Location**: `test/e2e/` directory
**Purpose**: Test complete application stack with real HTTP server
**Dependencies**: Test database, actual server startup

#### E2E Test Environment Setup

E2E tests require a PostgreSQL test database. Set up environment variables:

```bash
export DB_USERNAME="test"
export DB_PASSWORD="test"
export DB_DATABASE="wonder_test"
```

#### Running E2E Tests

```bash
# Run E2E tests with environment setup
DB_USERNAME="test" DB_PASSWORD="test" DB_DATABASE="wonder_test" ./scripts/test.sh e2e

# Run with verbose output
DB_USERNAME="test" DB_PASSWORD="test" DB_DATABASE="wonder_test" ./scripts/test.sh e2e -v

# Run directly with Go
DB_USERNAME="test" DB_PASSWORD="test" DB_DATABASE="wonder_test" go test ./test/e2e/... -v
```

#### E2E Test Features

**Real Server Startup**: Tests start an actual HTTP server using the production server code.

**Complete Request Flow**: HTTP requests go through the entire application stack:
- Gin router and middleware
- HTTP handlers
- Application services
- Domain logic
- Database repositories
- Response serialization

**Database Integration**: Uses real PostgreSQL connections with automatic test data cleanup.

**Test Isolation**: Each test run automatically cleans up test data using GORM operations.

**Environment Configuration**: Uses `testing` environment configuration for database isolation.

#### E2E Test Examples

```go
// Health check test
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
})

// User registration flow test
t.Run("User Registration E2E Flow", func(t *testing.T) {
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
})
```

## Coverage Reports

### Generate Coverage Reports

```bash
# Generate HTML coverage report
./scripts/test.sh coverage

# Coverage files generated:
# - coverage.out: Raw coverage data
# - coverage.html: HTML coverage report
```

### Coverage Analysis

```bash
# View coverage summary in terminal
go tool cover -func=coverage.out

# Open HTML coverage report
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Current Coverage Stats

- **Domain Layer**: 100% (User entity and business logic)
- **Application Layer**: 94.1% (User service with mock dependencies)
- **Interface Layer**: 100% (HTTP handlers with request/response validation)
- **Repository Layer**: 84.6% (Database operations and GORM integration)
- **Test Utilities**: 84.8% (Builder pattern implementation)
- **Overall Coverage**: 93.4% ✅

## Testing Best Practices

### 1. Test Isolation

- Each test should be independent and not rely on other tests
- Use test data builders for consistent test setup
- Clean up test data after each test run

### 2. Test Data Management

- Use unique identifiers for test data (e.g., emails ending with `test.com`)
- Implement cleanup functions for database tests
- Use mock objects for unit tests to avoid external dependencies

### 3. Error Testing

- Test both success and failure scenarios
- Verify proper error handling and error messages
- Test edge cases and boundary conditions

### 4. Performance Considerations

- Use `testing.Short()` to skip long-running tests in short mode
- Implement benchmarks for performance-critical code
- Monitor test execution time and optimize slow tests

## Continuous Integration

### Local Pre-commit Checks

```bash
# Run all tests before committing
./scripts/test.sh all

# Run with coverage verification
./scripts/test.sh coverage

# Format and lint code
go fmt ./...
go vet ./...
go mod tidy
```

### CI/CD Pipeline Integration

The test runner script is designed for CI/CD integration:

```yaml
# Example CI configuration
test:
  script:
    - source .envrc
    - ./scripts/test.sh all --race
    - ./scripts/test.sh coverage
```

## Troubleshooting

### Common Issues

**1. Database Connection Errors**
```bash
# Ensure database credentials are set
export DB_USERNAME="test"
export DB_PASSWORD="test"
export DB_DATABASE="wonder_test"

# Verify database is running and accessible
psql -h localhost -U test -d wonder_test -c "SELECT 1;"
```

**2. Port Conflicts in E2E Tests**
- E2E tests use port 8080 by default
- Ensure no other services are running on the same port
- Check for existing server processes

**3. Test Data Conflicts**
- E2E tests automatically clean up test data
- If tests fail unexpectedly, manually clean test data:
```sql
DELETE FROM users WHERE email LIKE '%test.com';
```

**4. Environment Variable Issues**
```bash
# Source environment configuration
source .envrc

# Verify Go environment
go env GOPROXY
go env GOSUMDB
```

### Debug Mode

```bash
# Run tests with verbose output for debugging
./scripts/test.sh e2e -v

# Run specific test function
go test ./test/e2e -run TestServerE2E/Health_Check_E2E -v

# Enable race detection for concurrency issues
./scripts/test.sh all --race
```

This comprehensive testing framework ensures code quality, reliability, and maintainability across the entire Wonder application stack.