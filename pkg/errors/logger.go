package errors

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ErrorLogger provides structured error logging functionality
type ErrorLogger interface {
	LogError(ctx context.Context, err error, traceID string, metadata map[string]interface{})
	LogErrorWithLevel(ctx context.Context, err error, level LogLevel, traceID string, metadata map[string]interface{})
}

// LogLevel represents the severity level of the error
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// DefaultErrorLogger is a basic implementation of ErrorLogger
type DefaultErrorLogger struct {
	serviceName string
}

// NewDefaultErrorLogger creates a new default error logger
func NewDefaultErrorLogger(serviceName string) *DefaultErrorLogger {
	return &DefaultErrorLogger{
		serviceName: serviceName,
	}
}

// LogError logs an error with default ERROR level
func (l *DefaultErrorLogger) LogError(ctx context.Context, err error, traceID string, metadata map[string]interface{}) {
	l.LogErrorWithLevel(ctx, err, LogLevelError, traceID, metadata)
}

// LogErrorWithLevel logs an error with specified level
func (l *DefaultErrorLogger) LogErrorWithLevel(ctx context.Context, err error, level LogLevel, traceID string, metadata map[string]interface{}) {
	// Create log entry with structured fields
	logEntry := map[string]interface{}{
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"service":      l.serviceName,
		"level":        string(level),
		"trace_id":     traceID,
		"error":        err.Error(),
		"error_type":   l.getErrorType(err),
		"retryable":    l.isRetryable(err),
	}

	// Add error-specific details using the new classification system
	errorType := Classifier.ClassifyError(err)
	logEntry["layer"] = string(errorType)

	if baseErr, ok := err.(BaseError); ok {
		logEntry["error_code"] = baseErr.Code()
		logEntry["error_details"] = baseErr.Details()
	}

	// Add custom metadata
	for k, v := range metadata {
		logEntry[k] = v
	}

	// Log the structured entry (in production, this would use a proper logger like logrus or zap)
	l.logStructured(logEntry)
}

// getErrorType returns the type of error for classification
func (l *DefaultErrorLogger) getErrorType(err error) string {
	errorType := Classifier.ClassifyError(err)
	switch errorType {
	case ErrorTypeDomain:
		return "domain_error"
	case ErrorTypeApplication:
		return "application_error"
	case ErrorTypeInfrastructure:
		return "infrastructure_error"
	case ErrorTypeInterface:
		return "interface_error"
	case ErrorTypeSystem:
		return "system_error"
	default:
		return "unknown_error"
	}
}

// isRetryable determines if the error is retryable
func (l *DefaultErrorLogger) isRetryable(err error) bool {
	return IsRetryable(err)
}

// logStructured outputs the structured log entry
func (l *DefaultErrorLogger) logStructured(entry map[string]interface{}) {
	// In production, this would use a proper structured logger
	// For now, we'll format it as JSON-like output
	logLine := fmt.Sprintf("[%s] %s - %s",
		entry["level"],
		entry["timestamp"],
		entry["error"])

	if traceID, ok := entry["trace_id"].(string); ok && traceID != "" {
		logLine += fmt.Sprintf(" (trace_id: %s)", traceID)
	}

	if layer, ok := entry["layer"].(string); ok {
		logLine += fmt.Sprintf(" [layer: %s]", layer)
	}

	if code, ok := entry["error_code"].(string); ok {
		logLine += fmt.Sprintf(" [code: %s]", code)
	}

	if retryable, ok := entry["retryable"].(bool); ok && retryable {
		logLine += " [retryable]"
	}

	log.Println(logLine)
}

// LogErrorsMiddleware creates Gin middleware for error logging
func LogErrorsMiddleware(logger ErrorLogger) func(ctx context.Context, err error, traceID string) {
	return func(ctx context.Context, err error, traceID string) {
		// Determine log level based on error type
		level := LogLevelError
		errorType := Classifier.ClassifyError(err)

		switch errorType {
		case ErrorTypeDomain:
			level = LogLevelWarn // Domain errors are usually user errors
		case ErrorTypeApplication:
			level = LogLevelWarn // Application errors are business logic issues
		case ErrorTypeInfrastructure:
			if IsRetryable(err) {
				level = LogLevelWarn // Retryable errors are warnings
			} else {
				level = LogLevelError // Non-retryable infrastructure errors are serious
			}
		default:
			level = LogLevelError
		}

		logger.LogErrorWithLevel(ctx, err, level, traceID, map[string]interface{}{
			"component": "http_handler",
		})
	}
}