package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// simpleLogger implements Logger interface with minimal overhead
type simpleLogger struct {
	logger    *logrus.Logger
	baseKV    map[string]interface{} // Pre-stored key-values for performance
	component string
	layer     string
}

// NewLogger creates a new simplified logger instance
func NewLogger() Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)

	return &simpleLogger{
		logger: l,
		baseKV: make(map[string]interface{}),
	}
}

// Debug logs a debug level message with key-value pairs
func (s *simpleLogger) Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	if !s.DebugEnabled() {
		return
	}
	s.logWithLevel(ctx, logrus.DebugLevel, msg, keyvals...)
}

// Info logs an info level message with key-value pairs
func (s *simpleLogger) Info(ctx context.Context, msg string, keyvals ...interface{}) {
	if !s.InfoEnabled() {
		return
	}
	s.logWithLevel(ctx, logrus.InfoLevel, msg, keyvals...)
}

// Warn logs a warning level message with key-value pairs
func (s *simpleLogger) Warn(ctx context.Context, msg string, keyvals ...interface{}) {
	s.logWithLevel(ctx, logrus.WarnLevel, msg, keyvals...)
}

// Error logs an error level message with key-value pairs
func (s *simpleLogger) Error(ctx context.Context, msg string, keyvals ...interface{}) {
	s.logWithLevel(ctx, logrus.ErrorLevel, msg, keyvals...)
}

// DebugEnabled returns true if debug logging is enabled
func (s *simpleLogger) DebugEnabled() bool {
	return s.logger.IsLevelEnabled(logrus.DebugLevel)
}

// InfoEnabled returns true if info logging is enabled
func (s *simpleLogger) InfoEnabled() bool {
	return s.logger.IsLevelEnabled(logrus.InfoLevel)
}

// With returns a new logger with additional key-value pairs
func (s *simpleLogger) With(keyvals ...interface{}) Logger {
	newKV := make(map[string]interface{})
	// Copy existing key-values
	for k, v := range s.baseKV {
		newKV[k] = v
	}

	// Add new key-values
	fields := s.parseKeyvals(keyvals...)
	for k, v := range fields {
		newKV[k] = v
	}

	return &simpleLogger{
		logger:    s.logger,
		baseKV:    newKV,
		component: s.component,
		layer:     s.layer,
	}
}

// WithLayer returns a new logger with DDD layer context
func (s *simpleLogger) WithLayer(layer string) Logger {
	newLogger := s.With("layer", layer)
	if sl, ok := newLogger.(*simpleLogger); ok {
		sl.layer = layer
	}
	return newLogger
}

// WithComponent returns a new logger with component context
func (s *simpleLogger) WithComponent(component string) Logger {
	newLogger := s.With("component", component)
	if sl, ok := newLogger.(*simpleLogger); ok {
		sl.component = component
	}
	return newLogger
}

// WithError returns a new logger with error context
func (s *simpleLogger) WithError(err error) Logger {
	if err == nil {
		return s
	}
	return s.With("error", err.Error())
}

// logWithLevel performs the actual logging with automatic context extraction
func (s *simpleLogger) logWithLevel(ctx context.Context, level logrus.Level, msg string, keyvals ...interface{}) {
	fields := logrus.Fields{}

	// Add base key-values
	for k, v := range s.baseKV {
		fields[k] = v
	}

	// Auto-extract from context
	if traceID := extractTraceID(ctx); traceID != "" {
		fields["trace_id"] = traceID
	}

	// Add provided key-values
	kvFields := s.parseKeyvals(keyvals...)
	for k, v := range kvFields {
		fields[k] = v
	}

	// Add layer and component if set
	if s.layer != "" {
		fields["layer"] = s.layer
	}
	if s.component != "" {
		fields["component"] = s.component
	}

	s.logger.WithFields(fields).Log(level, msg)
}

// parseKeyvals converts variadic interface{} to logrus.Fields
func (s *simpleLogger) parseKeyvals(keyvals ...interface{}) logrus.Fields {
	fields := logrus.Fields{}

	if len(keyvals)%2 != 0 {
		// Add unpaired value as "extra"
		keyvals = append(keyvals, "MISSING_VALUE")
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("key_%d", i)
		}
		fields[key] = keyvals[i+1]
	}

	return fields
}

// extractTraceID extracts trace ID from context
func extractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Try to extract from standard context keys
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if str, ok := traceID.(string); ok {
			return str
		}
	}

	if traceID := ctx.Value("x-trace-id"); traceID != nil {
		if str, ok := traceID.(string); ok {
			return str
		}
	}

	return ""
}

// Global logger instance for convenience
var defaultLogger Logger

// Initialize sets up the global logger
func Initialize() {
	defaultLogger = NewLogger()
}

// Get returns the global logger instance
func Get() Logger {
	if defaultLogger == nil {
		Initialize()
	}
	return defaultLogger
}

// Convenience functions for global logger
func LogDebug(ctx context.Context, msg string, keyvals ...interface{}) {
	Get().Debug(ctx, msg, keyvals...)
}

func LogInfo(ctx context.Context, msg string, keyvals ...interface{}) {
	Get().Info(ctx, msg, keyvals...)
}

func LogWarn(ctx context.Context, msg string, keyvals ...interface{}) {
	Get().Warn(ctx, msg, keyvals...)
}

func LogError(ctx context.Context, msg string, keyvals ...interface{}) {
	Get().Error(ctx, msg, keyvals...)
}