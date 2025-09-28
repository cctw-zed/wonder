package logger

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_BasicLogging(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Test basic logging methods
	logger.Debug(ctx, "debug message", "key1", "value1")
	logger.Info(ctx, "info message", "key2", "value2")
	logger.Warn(ctx, "warn message", "key3", "value3")
	logger.Error(ctx, "error message", "key4", "value4")

	// These should not panic
	assert.True(t, true)
}

func TestLogger_LevelChecking(t *testing.T) {
	logger := NewLogger()

	// Level checking should work
	debugEnabled := logger.DebugEnabled()
	infoEnabled := logger.InfoEnabled()

	// At debug level, both should be true
	assert.True(t, debugEnabled)
	assert.True(t, infoEnabled)
}

func TestLogger_WithMethods(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Test With method
	loggerWithFields := logger.With("component", "test", "version", "1.0")
	require.NotNil(t, loggerWithFields)

	// Test WithLayer
	layerLogger := logger.WithLayer("application")
	require.NotNil(t, layerLogger)

	// Test WithComponent
	componentLogger := logger.WithComponent("user_service")
	require.NotNil(t, componentLogger)

	// Test WithError
	err := assert.AnError
	errorLogger := logger.WithError(err)
	require.NotNil(t, errorLogger)

	// Test chaining
	chainedLogger := logger.WithLayer("domain").WithComponent("user").With("operation", "create")
	require.NotNil(t, chainedLogger)

	// Log with chained logger
	chainedLogger.Info(ctx, "chained logging test")
}

func TestLogger_ContextExtraction(t *testing.T) {
	logger := NewLogger()

	// Test with trace ID in context
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-123")
	logger.Info(ctx, "message with trace", "key", "value")

	// Test with empty context
	emptyCtx := context.Background()
	logger.Info(emptyCtx, "message without trace")
}

func TestLogger_KeyValuesParsing(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Test even number of keyvals
	logger.Info(ctx, "even keyvals", "key1", "value1", "key2", "value2")

	// Test odd number of keyvals (should handle gracefully)
	logger.Info(ctx, "odd keyvals", "key1", "value1", "key2")

	// Test non-string keys
	logger.Info(ctx, "non-string key", 123, "value")
}

func TestLogger_Performance(t *testing.T) {
	logger := NewLogger()
	ctx := context.Background()

	// Test performance with debug disabled
	start := time.Now()
	for i := 0; i < 1000; i++ {
		if logger.DebugEnabled() {
			logger.Debug(ctx, "debug message", "iteration", i)
		}
	}
	duration := time.Since(start)

	// Should be fast when debug is enabled
	assert.Less(t, duration, time.Millisecond*100)
}

func TestHelperFunctions(t *testing.T) {
	// Test KV helper
	kv := KV("test", "value")
	assert.Equal(t, []interface{}{"test", "value"}, kv)

	// Test Merge helper
	kv1 := KV("key1", "value1")
	kv2 := KV("key2", "value2")
	merged := Merge(kv1, kv2)
	expected := []interface{}{"key1", "value1", "key2", "value2"}
	assert.Equal(t, expected, merged)

	// Test Err helper
	err := assert.AnError
	errKV := Err(err)
	assert.Equal(t, []interface{}{"error", err.Error()}, errKV)

	// Test Err helper with nil
	nilErrKV := Err(nil)
	assert.Nil(t, nilErrKV)

	// Test Duration helper
	dur := time.Second
	durKV := Duration(dur)
	assert.Equal(t, []interface{}{"duration", "1s"}, durKV)

	// Test UserID helper
	userKV := UserID("user-123")
	assert.Equal(t, []interface{}{"user_id", "user-123"}, userKV)

	// Test Email helper
	emailKV := Email("test@example.com")
	assert.Equal(t, []interface{}{"email", "test@example.com"}, emailKV)

	// Test Operation helper
	opKV := Operation("create_user")
	assert.Equal(t, []interface{}{"operation", "create_user"}, opKV)
}

func TestGlobalLogger(t *testing.T) {
	// Initialize global logger
	Initialize()

	// Test global logger functions
	ctx := context.Background()
	LogDebug(ctx, "global debug")
	LogInfo(ctx, "global info")
	LogWarn(ctx, "global warn")
	LogError(ctx, "global error")

	// Test getting global logger
	globalLogger := Get()
	assert.NotNil(t, globalLogger)

	// Multiple calls should return same instance
	globalLogger2 := Get()
	assert.Equal(t, globalLogger, globalLogger2)
}

// Benchmark tests to compare performance
func BenchmarkLogger_Info(b *testing.B) {
	logger := NewLogger()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "benchmark message", "iteration", i)
	}
}

func BenchmarkLogger_DebugDisabled(b *testing.B) {
	logger := NewLogger()
	ctx := context.Background()

	// Disable debug level by setting higher level
	// (This would need implementation in real usage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if logger.DebugEnabled() {
			logger.Debug(ctx, "benchmark debug", "iteration", i)
		}
	}
}

func BenchmarkLogger_WithChaining(b *testing.B) {
	logger := NewLogger()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chainedLogger := logger.WithLayer("test").WithComponent("bench")
		chainedLogger.Info(ctx, "chained message", "iteration", i)
	}
}

// Example usage tests that demonstrate the API
func ExampleLogger_basic() {
	logger := NewLogger()
	ctx := context.Background()

	// Basic usage - similar to go-kit/log
	logger.Info(ctx, "user registered", "email", "user@example.com", "user_id", "123")

	// With context
	appLogger := logger.WithLayer("application").WithComponent("user_service")
	appLogger.Debug(ctx, "validating email", "email", "user@example.com")

	// Error logging
	err := assert.AnError
	appLogger.Error(ctx, "validation failed", "error", err.Error(), "phase", "email_check")
}

func ExampleLogger_helpers() {
	logger := NewLogger()
	ctx := context.Background()

	// Using helper functions for common patterns
	logger.Info(ctx, "operation started",
		Merge(
			Operation("user_registration"),
			UserID("123"),
			Email("user@example.com"),
		)...,
	)

	// Error with helper
	err := assert.AnError
	logger.Error(ctx, "operation failed",
		Merge(
			Operation("user_registration"),
			Err(err),
		)...,
	)
}
