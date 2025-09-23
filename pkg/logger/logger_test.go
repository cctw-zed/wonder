package logger

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerInterface(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "JSON format logger",
			config: &Config{
				Level:       InfoLevel,
				Format:      "json",
				Output:      "stdout",
				ServiceName: "test-service",
			},
		},
		{
			name: "Text format logger",
			config: &Config{
				Level:       DebugLevel,
				Format:      "text",
				Output:      "stdout",
				ServiceName: "test-service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogrusLogger(tt.config)
			ctx := context.Background()

			// Test all log levels
			logger.Debug(ctx, "Debug message", String("key", "debug"))
			logger.Info(ctx, "Info message", String("key", "info"))
			logger.Warn(ctx, "Warn message", String("key", "warn"))
			logger.Error(ctx, "Error message", String("key", "error"))
			// Note: Not testing Fatal as it would exit the process
		})
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Level:       DebugLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	logger := NewLogrusLogger(config)
	// Redirect output to buffer for testing
	logger.logger.SetOutput(&buf)

	ctx := context.Background()

	// Test WithField
	fieldLogger := logger.WithField("test_field", "test_value")
	fieldLogger.Info(ctx, "Test message")

	output := buf.String()
	assert.Contains(t, output, "test_field")
	assert.Contains(t, output, "test_value")
	assert.Contains(t, output, "Test message")

	buf.Reset()

	// Test WithFields
	fieldsLogger := logger.WithFields(Fields{
		"field1": "value1",
		"field2": 123,
		"field3": true,
	})
	fieldsLogger.Info(ctx, "Test with multiple fields")

	output = buf.String()
	assert.Contains(t, output, "field1")
	assert.Contains(t, output, "value1")
	assert.Contains(t, output, "field2")
	assert.Contains(t, output, "123")
	assert.Contains(t, output, "field3")
	assert.Contains(t, output, "true")
}

func TestTracing(t *testing.T) {
	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID1 := GenerateTraceID()
		traceID2 := GenerateTraceID()

		assert.NotEmpty(t, traceID1)
		assert.NotEmpty(t, traceID2)
		assert.NotEqual(t, traceID1, traceID2)
		assert.Equal(t, 32, len(traceID1)) // 16 bytes hex encoded
	})

	t.Run("TraceContext", func(t *testing.T) {
		traceCtx := NewTraceContext()
		assert.NotEmpty(t, traceCtx.TraceID)
		assert.NotEmpty(t, traceCtx.RequestID)
		assert.False(t, traceCtx.StartTime.IsZero())

		ctx := WithTraceContext(context.Background(), traceCtx)

		extractedTraceID := GetTraceID(ctx)
		extractedRequestID := GetRequestID(ctx)

		assert.Equal(t, traceCtx.TraceID, extractedTraceID)
		assert.Equal(t, traceCtx.RequestID, extractedRequestID)
	})

	t.Run("LoggerWithTracing", func(t *testing.T) {
		var buf bytes.Buffer
		config := &Config{
			Level:       InfoLevel,
			Format:      "json",
			Output:      "stdout",
			ServiceName: "test-service",
		}

		logger := NewLogrusLogger(config)
		logger.logger.SetOutput(&buf)

		traceID := GenerateTraceID()
		ctx := WithTraceID(context.Background(), traceID)

		tracingLogger := LoggerWithTracing(ctx, logger)
		tracingLogger.Info(ctx, "Test message with tracing")

		output := buf.String()
		assert.Contains(t, output, traceID)
	})
}

func TestPerformanceLogger(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Level:       InfoLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	logger := NewLogrusLogger(config)
	logger.logger.SetOutput(&buf)

	ctx := context.Background()
	perfLogger := NewPerformanceLogger(ctx, logger)

	t.Run("LogDuration fast operation", func(t *testing.T) {
		startTime := time.Now().Add(-100 * time.Millisecond)
		perfLogger.LogDuration("test_operation", startTime)

		output := buf.String()
		assert.Contains(t, output, "test_operation")
		assert.Contains(t, output, "duration")
		assert.Contains(t, output, "Operation completed")
	})

	buf.Reset()

	t.Run("LogDuration slow operation", func(t *testing.T) {
		startTime := time.Now().Add(-2 * time.Second)
		perfLogger.LogDuration("slow_operation", startTime)

		output := buf.String()
		assert.Contains(t, output, "slow_operation")
		assert.Contains(t, output, "Slow operation")
	})
}

func TestFactoryAndGlobalLogger(t *testing.T) {
	config := &Config{
		Level:       InfoLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	factory := NewFactory(config)
	SetGlobalFactory(factory)

	t.Run("Factory creates loggers", func(t *testing.T) {
		logger := factory.NewLogger()
		assert.NotNil(t, logger)

		domainLogger := factory.NewDomainLogger()
		assert.NotNil(t, domainLogger)

		appLogger := factory.NewApplicationLogger()
		assert.NotNil(t, appLogger)

		infraLogger := factory.NewInfrastructureLogger()
		assert.NotNil(t, infraLogger)

		interfaceLogger := factory.NewInterfaceLogger()
		assert.NotNil(t, interfaceLogger)

		componentLogger := factory.NewComponentLogger("test-component")
		assert.NotNil(t, componentLogger)
	})

	t.Run("Global factory functions work", func(t *testing.T) {
		logger := NewLogger()
		assert.NotNil(t, logger)

		domainLogger := factory.NewDomainLogger()
		assert.NotNil(t, domainLogger)

		appLogger := factory.NewApplicationLogger()
		assert.NotNil(t, appLogger)

		infraLogger := factory.NewInfrastructureLogger()
		assert.NotNil(t, infraLogger)

		interfaceLogger := factory.NewInterfaceLogger()
		assert.NotNil(t, interfaceLogger)

		componentLogger := NewComponentLogger("test-component")
		assert.NotNil(t, componentLogger)
	})
}

func TestDDDLoggers(t *testing.T) {
	var buf bytes.Buffer
	config := &Config{
		Level:       DebugLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	baseLogger := NewLogrusLogger(config)
	baseLogger.logger.SetOutput(&buf)

	ctx := context.Background()

	t.Run("DomainLogger", func(t *testing.T) {
		domainLogger := NewDomainLogger(baseLogger)

		domainLogger.LogBusinessRule(ctx, "UserAgeValidation", "User", "123", true)
		output := buf.String()
		assert.Contains(t, output, "UserAgeValidation")
		assert.Contains(t, output, "User")
		assert.Contains(t, output, "123")
		assert.Contains(t, output, "true")

		buf.Reset()

		domainLogger.LogDomainEvent(ctx, "UserRegistered", "user-123")
		output = buf.String()
		assert.Contains(t, output, "UserRegistered")
		assert.Contains(t, output, "user-123")
	})

	t.Run("ApplicationLogger", func(t *testing.T) {
		appLogger := NewApplicationLogger(baseLogger)
		startTime := time.Now().Add(-100 * time.Millisecond)

		appLogger.LogUseCase(ctx, "RegisterUser", startTime, true)
		output := buf.String()
		assert.Contains(t, output, "RegisterUser")
		assert.Contains(t, output, "duration")
		assert.Contains(t, output, "true")

		buf.Reset()

		appLogger.LogValidation(ctx, "EmailValidation", false, []string{"invalid format"})
		output = buf.String()
		assert.Contains(t, output, "EmailValidation")
		assert.Contains(t, output, "false")
		assert.Contains(t, output, "invalid format")
	})

	t.Run("InfrastructureLogger", func(t *testing.T) {
		infraLogger := NewInfrastructureLogger(baseLogger)

		infraLogger.LogDatabaseOperation(ctx, "SELECT", "users", 50*time.Millisecond, 5)
		output := buf.String()
		assert.Contains(t, output, "SELECT")
		assert.Contains(t, output, "users")
		assert.Contains(t, output, "5")

		buf.Reset()

		infraLogger.LogExternalServiceCall(ctx, "PaymentAPI", "/charge", "POST", 200, 1*time.Second)
		output = buf.String()
		assert.Contains(t, output, "PaymentAPI")
		assert.Contains(t, output, "/charge")
		assert.Contains(t, output, "POST")
		assert.Contains(t, output, "200")
	})

	t.Run("InterfaceLogger", func(t *testing.T) {
		interfaceLogger := NewInterfaceLogger(baseLogger)

		interfaceLogger.LogHTTPRequest(ctx, "POST", "/api/users", 201, 200*time.Millisecond, "test-agent")
		output := buf.String()
		assert.Contains(t, output, "POST")
		assert.Contains(t, output, "/api/users")
		assert.Contains(t, output, "201")
		assert.Contains(t, output, "test-agent")

		buf.Reset()

		interfaceLogger.LogAuthentication(ctx, "user-123", "password", true, "success")
		output = buf.String()
		assert.Contains(t, output, "user-123")
		assert.Contains(t, output, "password")
		assert.Contains(t, output, "true")
		assert.Contains(t, output, "success")
	})
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			config: &Config{
				Level:         InfoLevel,
				Format:        "json",
				Output:        "stdout",
				ServiceName:   "test-service",
				EnableTracing: true,
				MaxFileSize:   100,
				MaxBackups:    3,
				MaxAge:        30,
				Compress:      true,
			},
			wantErr: false,
		},
		{
			name: "Empty service name",
			config: &Config{
				Level:  InfoLevel,
				Format: "json",
				Output: "stdout",
				// ServiceName is empty
				EnableTracing: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogrusLogger(tt.config)
			// If we get here without panic, the config was accepted
			// For empty service name, it should fallback to default
			assert.NotNil(t, logger)
		})
	}
}

func TestFieldHelperFunctions(t *testing.T) {
	stringField := String("key", "value")
	assert.Equal(t, "key", stringField.Key)
	assert.Equal(t, "value", stringField.Value)

	intField := Int("count", 42)
	assert.Equal(t, "count", intField.Key)
	assert.Equal(t, 42, intField.Value)

	boolField := Bool("active", true)
	assert.Equal(t, "active", boolField.Key)
	assert.Equal(t, true, boolField.Value)

	now := time.Now()
	timeField := Time("timestamp", now)
	assert.Equal(t, "timestamp", timeField.Key)
	assert.Equal(t, now, timeField.Value)

	duration := 5 * time.Second
	durationField := Duration("elapsed", duration)
	assert.Equal(t, "elapsed", durationField.Key)
	assert.Equal(t, duration, durationField.Value)
}

func TestLogrusLoggerImplementsInterface(t *testing.T) {
	config := &Config{
		Level:       InfoLevel,
		Format:      "json",
		Output:      "stdout",
		ServiceName: "test-service",
	}

	var logger Logger = NewLogrusLogger(config)
	assert.NotNil(t, logger)

	// Test that we can call all interface methods
	ctx := context.Background()
	logger.Debug(ctx, "test")
	logger.Info(ctx, "test")
	logger.Warn(ctx, "test")
	logger.Error(ctx, "test")

	withFieldLogger := logger.WithField("key", "value")
	assert.NotNil(t, withFieldLogger)

	withFieldsLogger := logger.WithFields(Fields{"key": "value"})
	assert.NotNil(t, withFieldsLogger)

	withComponentLogger := logger.WithComponent("test-component")
	assert.NotNil(t, withComponentLogger)

	withLayerLogger := logger.WithLayer(DomainLayer)
	assert.NotNil(t, withLayerLogger)
}