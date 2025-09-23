package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	wonderErrors "github.com/cctw-zed/wonder/pkg/errors"
)

func TestErrorLogger(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Level:       DebugLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	baseLogger := NewLogrusLogger(config)
	baseLogger.logger.SetOutput(&buf)

	errorLogger := NewErrorLogger(baseLogger)
	ctx := context.Background()

	t.Run("LogError with BaseError", func(t *testing.T) {
		domainErr := wonderErrors.NewValidationError(wonderErrors.CodeValidationError, "email", "invalid@", "test validation error")

		errorLogger.LogError(ctx, domainErr, "trace-123", map[string]interface{}{
			"user_id": "user-456",
		})

		output := buf.String()
		assert.Contains(t, output, "test validation error")
		assert.Contains(t, output, "trace-123")
		assert.Contains(t, output, "user-456")
		assert.Contains(t, output, string(wonderErrors.CodeValidationError))
		assert.Contains(t, output, "domain_error")
	})

	buf.Reset()

	t.Run("LogErrorWithLevel", func(t *testing.T) {
		appErr := wonderErrors.NewEntityNotFoundError("User", "123")

		errorLogger.LogErrorWithLevel(ctx, appErr, wonderErrors.LogLevelWarn, "trace-456", map[string]interface{}{
			"operation": "GetUser",
		})

		output := buf.String()
		assert.Contains(t, output, "User")
		assert.Contains(t, output, "123")
		assert.Contains(t, output, "trace-456")
		assert.Contains(t, output, "GetUser")
		assert.Contains(t, output, "application_error")
	})

	buf.Reset()

	t.Run("LogError with context tracing", func(t *testing.T) {
		traceID := GenerateTraceID()
		ctxWithTrace := WithTraceID(context.Background(), traceID)

		infraErr := wonderErrors.NewDatabaseError("SELECT", "users", errors.New("connection failed"), false)

		errorLogger.LogError(ctxWithTrace, infraErr, "", map[string]interface{}{
			"database": "postgres",
		})

		output := buf.String()
		assert.Contains(t, output, "connection failed")
		assert.Contains(t, output, traceID)
		assert.Contains(t, output, "postgres")
		assert.Contains(t, output, "infrastructure_error")
	})
}

func TestCreateErrorLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Level:       DebugLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	baseLogger := NewLogrusLogger(config)
	baseLogger.logger.SetOutput(&buf)

	middleware := CreateErrorLoggerMiddleware(baseLogger)
	ctx := context.Background()

	t.Run("Domain error logs as warning", func(t *testing.T) {
		domainErr := wonderErrors.NewBusinessRuleError("age_validation", "age must be positive")
		middleware(ctx, domainErr, "trace-123")

		output := buf.String()
		assert.Contains(t, output, "age must be positive")
		assert.Contains(t, output, "trace-123")
		assert.Contains(t, output, "http_handler")
		// Should be warn level for domain errors
	})

	buf.Reset()

	t.Run("Infrastructure retryable error logs as warning", func(t *testing.T) {
		infraErr := wonderErrors.NewNetworkError("payment", "/charge", "POST", errors.New("timeout"), true)
		middleware(ctx, infraErr, "trace-456")

		output := buf.String()
		assert.Contains(t, output, "timeout")
		assert.Contains(t, output, "trace-456")
		assert.Contains(t, output, "retryable")
	})

	buf.Reset()

	t.Run("Infrastructure non-retryable error logs as error", func(t *testing.T) {
		infraErr := wonderErrors.NewDatabaseError("SELECT", "users", errors.New("connection failed"), false)
		middleware(ctx, infraErr, "trace-789")

		output := buf.String()
		assert.Contains(t, output, "connection failed")
		assert.Contains(t, output, "trace-789")
		assert.Contains(t, output, "false") // retryable = false
	})
}

func TestConvertErrorLogLevel(t *testing.T) {
	tests := []struct {
		input    wonderErrors.LogLevel
		expected Level
	}{
		{wonderErrors.LogLevelDebug, DebugLevel},
		{wonderErrors.LogLevelInfo, InfoLevel},
		{wonderErrors.LogLevelWarn, WarnLevel},
		{wonderErrors.LogLevelError, ErrorLevel},
		{wonderErrors.LogLevelFatal, FatalLevel},
		{"unknown", ErrorLevel}, // default case
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := convertErrorLogLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "Domain error",
			err:      wonderErrors.NewValidationError(wonderErrors.CodeValidationError, "test", "value", "test message"),
			expected: "domain_error",
		},
		{
			name:     "Application error",
			err:      wonderErrors.NewEntityNotFoundError("User", "123"),
			expected: "application_error",
		},
		{
			name:     "Infrastructure error",
			err:      wonderErrors.NewDatabaseError("SELECT", "users", errors.New("connection failed"), false),
			expected: "infrastructure_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getErrorType(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrorLoggerImplementsInterface(t *testing.T) {
	config := &Config{
		Level:       InfoLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	baseLogger := NewLogrusLogger(config)

	// Test that ErrorLogger implements wonderErrors.ErrorLogger interface
	var errorLogger wonderErrors.ErrorLogger = NewErrorLogger(baseLogger)
	assert.NotNil(t, errorLogger)

	// Test interface methods can be called
	ctx := context.Background()
	testErr := wonderErrors.NewValidationError(wonderErrors.CodeValidationError, "test", "value", "test message")

	errorLogger.LogError(ctx, testErr, "trace-123", map[string]interface{}{
		"test": "metadata",
	})

	errorLogger.LogErrorWithLevel(ctx, testErr, wonderErrors.LogLevelWarn, "trace-456", map[string]interface{}{
		"test": "metadata",
	})
}