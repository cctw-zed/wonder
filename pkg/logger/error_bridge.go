package logger

import (
	"context"

	"github.com/cctw-zed/wonder/pkg/errors"
)

// ErrorLogger implements the errors.ErrorLogger interface using our new logger
type ErrorLogger struct {
	logger Logger
}

// NewErrorLogger creates a new error logger bridge
func NewErrorLogger(logger Logger) *ErrorLogger {
	return &ErrorLogger{
		logger: logger.WithComponent("error_handler"),
	}
}

// LogError logs an error with default ERROR level
func (e *ErrorLogger) LogError(ctx context.Context, err error, traceID string, metadata map[string]interface{}) {
	e.LogErrorWithLevel(ctx, err, errors.LogLevelError, traceID, metadata)
}

// LogErrorWithLevel logs an error with specified level
func (e *ErrorLogger) LogErrorWithLevel(ctx context.Context, err error, level errors.LogLevel, traceID string, metadata map[string]interface{}) {
	// Convert errors.LogLevel to our Level type
	logLevel := convertErrorLogLevel(level)

	// Create a logger with tracing information
	logger := LoggerWithTracing(ctx, e.logger)

	// Add trace ID if provided and not already in context
	if traceID != "" && GetTraceID(ctx) == "" {
		logger = logger.WithTraceID(traceID)
	}

	// Add error-specific fields
	fields := []Field{
		String("error", err.Error()),
		String("error_type", getErrorType(err)),
		Bool("retryable", errors.IsRetryable(err)),
	}

	// Add error classification
	errorType := errors.Classifier.ClassifyError(err)
	fields = append(fields, String("layer", string(errorType)))

	// Add error details if available
	if baseErr, ok := err.(errors.BaseError); ok {
		fields = append(fields,
			String("error_code", string(baseErr.Code())),
			Any("error_details", baseErr.Details()),
		)
	}

	// Add custom metadata
	for k, v := range metadata {
		fields = append(fields, Any(k, v))
	}

	// Log with appropriate level
	message := "Error occurred"
	if errorType == errors.ErrorTypeDomain {
		message = "Domain error"
	} else if errorType == errors.ErrorTypeApplication {
		message = "Application error"
	} else if errorType == errors.ErrorTypeInfrastructure {
		message = "Infrastructure error"
	}

	logger.LogWithFields(ctx, logLevel, message, fields...)
}

// convertErrorLogLevel converts errors.LogLevel to our Level type
func convertErrorLogLevel(level errors.LogLevel) Level {
	switch level {
	case errors.LogLevelDebug:
		return DebugLevel
	case errors.LogLevelInfo:
		return InfoLevel
	case errors.LogLevelWarn:
		return WarnLevel
	case errors.LogLevelError:
		return ErrorLevel
	case errors.LogLevelFatal:
		return FatalLevel
	default:
		return ErrorLevel
	}
}

// getErrorType returns the type of error for classification
func getErrorType(err error) string {
	errorType := errors.Classifier.ClassifyError(err)
	switch errorType {
	case errors.ErrorTypeDomain:
		return "domain_error"
	case errors.ErrorTypeApplication:
		return "application_error"
	case errors.ErrorTypeInfrastructure:
		return "infrastructure_error"
	case errors.ErrorTypeInterface:
		return "interface_error"
	case errors.ErrorTypeSystem:
		return "system_error"
	default:
		return "unknown_error"
	}
}

// CreateErrorLoggerMiddleware creates middleware for error logging
func CreateErrorLoggerMiddleware(logger Logger) func(ctx context.Context, err error, traceID string) {
	errorLogger := NewErrorLogger(logger)

	return func(ctx context.Context, err error, traceID string) {
		// Determine log level based on error type
		level := errors.LogLevelError
		errorType := errors.Classifier.ClassifyError(err)

		switch errorType {
		case errors.ErrorTypeDomain:
			level = errors.LogLevelWarn // Domain errors are usually user errors
		case errors.ErrorTypeApplication:
			level = errors.LogLevelWarn // Application errors are business logic issues
		case errors.ErrorTypeInfrastructure:
			if errors.IsRetryable(err) {
				level = errors.LogLevelWarn // Retryable errors are warnings
			} else {
				level = errors.LogLevelError // Non-retryable infrastructure errors are serious
			}
		default:
			level = errors.LogLevelError
		}

		errorLogger.LogErrorWithLevel(ctx, err, level, traceID, map[string]interface{}{
			"component": "http_handler",
		})
	}
}

// Ensure ErrorLogger implements errors.ErrorLogger interface
var _ errors.ErrorLogger = (*ErrorLogger)(nil)