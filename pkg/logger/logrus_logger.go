package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogrusLogger implements the Logger interface using logrus
type LogrusLogger struct {
	logger      *logrus.Logger
	serviceName string
	fields      logrus.Fields
}

// NewLogrusLogger creates a new LogrusLogger instance
func NewLogrusLogger(config *Config) *LogrusLogger {
	logger := logrus.New()

	// Set log level
	level := parseLogLevel(config.Level)
	logger.SetLevel(level)

	// Set formatter based on format
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default: // console
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
		})
	}

	// Set output destination
	output := setupOutput(config)
	logger.SetOutput(output)

	serviceName := config.ServiceName
	if serviceName == "" {
		serviceName = "wonder"
	}

	return &LogrusLogger{
		logger:      logger,
		serviceName: serviceName,
		fields:      make(logrus.Fields),
	}
}

// setupOutput configures the logger output based on configuration
func setupOutput(config *Config) io.Writer {
	var outputs []io.Writer

	// Add console output
	switch config.Output {
	case "stderr":
		outputs = append(outputs, os.Stderr)
	case "file":
		if config.EnableFile && config.FilePath != "" {
			fileWriter := setupFileOutput(config)
			outputs = append(outputs, fileWriter)
		} else {
			outputs = append(outputs, os.Stdout)
		}
	default: // stdout
		outputs = append(outputs, os.Stdout)
	}

	// Add file output if enabled
	if config.EnableFile && config.FilePath != "" && config.Output != "file" {
		fileWriter := setupFileOutput(config)
		outputs = append(outputs, fileWriter)
	}

	if len(outputs) == 1 {
		return outputs[0]
	}

	return io.MultiWriter(outputs...)
}

// setupFileOutput creates a file writer with rotation
func setupFileOutput(config *Config) io.Writer {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
		// Fallback to stdout if directory creation fails
		return os.Stdout
	}

	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxFileSize, // MB
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		Compress:   config.Compress,
	}
}

// parseLogLevel converts string level to logrus level
func parseLogLevel(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

// Debug logs a debug level message
func (l *LogrusLogger) Debug(ctx context.Context, message string, fields ...Field) {
	l.LogWithFields(ctx, DebugLevel, message, fields...)
}

// Info logs an info level message
func (l *LogrusLogger) Info(ctx context.Context, message string, fields ...Field) {
	l.LogWithFields(ctx, InfoLevel, message, fields...)
}

// Warn logs a warning level message
func (l *LogrusLogger) Warn(ctx context.Context, message string, fields ...Field) {
	l.LogWithFields(ctx, WarnLevel, message, fields...)
}

// Error logs an error level message
func (l *LogrusLogger) Error(ctx context.Context, message string, fields ...Field) {
	l.LogWithFields(ctx, ErrorLevel, message, fields...)
}

// Fatal logs a fatal level message and exits
func (l *LogrusLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	l.LogWithFields(ctx, FatalLevel, message, fields...)
}

// LogWithFields logs a message with explicit fields and level
func (l *LogrusLogger) LogWithFields(ctx context.Context, level Level, message string, fields ...Field) {
	entry := l.logger.WithFields(l.buildLogFields(ctx, fields...))

	logrusLevel := parseLogLevel(level)
	entry.Log(logrusLevel, message)
}

// buildLogFields constructs the logrus fields from context and provided fields
func (l *LogrusLogger) buildLogFields(ctx context.Context, fields ...Field) logrus.Fields {
	logFields := logrus.Fields{
		"service": l.serviceName,
	}

	// Add logger instance fields
	for k, v := range l.fields {
		logFields[k] = v
	}

	// Extract trace ID from context if available
	if traceID := extractTraceID(ctx); traceID != "" {
		logFields["trace_id"] = traceID
	}

	// Add provided fields
	for _, field := range fields {
		logFields[field.Key] = field.Value
	}

	return logFields
}

// WithField returns a new logger with an additional field
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &LogrusLogger{
		logger:      l.logger,
		serviceName: l.serviceName,
		fields:      newFields,
	}
}

// WithFields returns a new logger with additional fields
func (l *LogrusLogger) WithFields(fields Fields) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &LogrusLogger{
		logger:      l.logger,
		serviceName: l.serviceName,
		fields:      newFields,
	}
}

// WithError returns a new logger with error information
func (l *LogrusLogger) WithError(err error) Logger {
	return l.WithField("error", err.Error())
}

// WithTraceID returns a new logger with trace ID
func (l *LogrusLogger) WithTraceID(traceID string) Logger {
	return l.WithField("trace_id", traceID)
}

// WithComponent returns a new logger with component information
func (l *LogrusLogger) WithComponent(component string) Logger {
	return l.WithField("component", component)
}

// WithLayer returns a new logger with DDD layer information
func (l *LogrusLogger) WithLayer(layer DDDLayer) Logger {
	return l.WithField("layer", string(layer))
}

// extractTraceID extracts trace ID from context
func extractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Check for common trace ID keys
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	if traceID, ok := ctx.Value("traceID").(string); ok {
		return traceID
	}
	if traceID, ok := ctx.Value("requestID").(string); ok {
		return traceID
	}
	if traceID, ok := ctx.Value("request_id").(string); ok {
		return traceID
	}

	return ""
}

// Ensure LogrusLogger implements Logger interface
var _ Logger = (*LogrusLogger)(nil)