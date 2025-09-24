package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// LogConfig represents logger configuration
type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, text
	Output     string // stdout, file, both
	FilePath   string // path to log file (when Output is file or both)
	EnableFile bool   // enable file logging
}

// simpleLogger implements Logger interface with minimal overhead
type simpleLogger struct {
	logger    *logrus.Logger
	baseKV    map[string]interface{} // Pre-stored key-values for performance
	component string
	layer     string
}

// NewLogger creates a new simplified logger instance with default configuration
func NewLogger() Logger {
	return NewLoggerWithConfig(LogConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	})
}

// NewLoggerWithConfig creates a new logger with specified configuration
func NewLoggerWithConfig(config LogConfig) Logger {
	l := logrus.New()

	// Set formatter
	if config.Format == "text" {
		l.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		})
	} else {
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	}

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.DebugLevel // fallback to debug
	}
	l.SetLevel(level)

	// Set output destination
	var output io.Writer = os.Stdout

	switch config.Output {
	case "file":
		if config.FilePath != "" {
			fileOutput, err := createLogFile(config.FilePath)
			if err != nil {
				// Fallback to stdout on file creation error
				fmt.Fprintf(os.Stderr, "Failed to create log file %s: %v. Using stdout.\n", config.FilePath, err)
			} else {
				output = fileOutput
			}
		}
	case "both":
		if config.FilePath != "" {
			fileOutput, err := createLogFile(config.FilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create log file %s: %v. Using stdout only.\n", config.FilePath, err)
			} else {
				output = io.MultiWriter(os.Stdout, fileOutput)
			}
		}
	default: // "stdout" or any other value
		output = os.Stdout
	}

	l.SetOutput(output)

	return &simpleLogger{
		logger: l,
		baseKV: make(map[string]interface{}),
	}
}

// createLogFile creates and opens a log file, ensuring the directory exists
func createLogFile(filePath string) (*os.File, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}

	return file, nil
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

	// Try to extract from standard context keys used by our middleware
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

	// Support for other common trace ID context keys
	if traceID := ctx.Value("traceID"); traceID != nil {
		if str, ok := traceID.(string); ok {
			return str
		}
	}

	return ""
}

// Global logger instance for convenience
var defaultLogger Logger

// Initialize sets up the global logger with default configuration
func Initialize() {
	defaultLogger = NewLogger()
}

// InitializeWithConfig sets up the global logger with custom configuration
func InitializeWithConfig(config LogConfig) {
	defaultLogger = NewLoggerWithConfig(config)
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