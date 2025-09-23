package logger

import (
	"context"
	"time"
)

// Logger defines the interface for all logging operations
type Logger interface {
	// Core logging methods with different levels
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, fields ...Field)
	Warn(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, fields ...Field)
	Fatal(ctx context.Context, message string, fields ...Field)

	// Structured logging with explicit fields
	LogWithFields(ctx context.Context, level Level, message string, fields ...Field)

	// With methods for adding context
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	WithTraceID(traceID string) Logger

	// Component and layer-specific loggers
	WithComponent(component string) Logger
	WithLayer(layer DDDLayer) Logger
}

// Level represents the logging level
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// DDDLayer represents the DDD architecture layers
type DDDLayer string

const (
	DomainLayer       DDDLayer = "domain"
	ApplicationLayer  DDDLayer = "application"
	InfrastructureLayer DDDLayer = "infrastructure"
	InterfaceLayer    DDDLayer = "interface"
)

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Fields represents a collection of fields
type Fields map[string]interface{}

// LogEntry represents a complete log entry
type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Level       Level     `json:"level"`
	Message     string    `json:"message"`
	TraceID     string    `json:"trace_id,omitempty"`
	Component   string    `json:"component,omitempty"`
	Layer       DDDLayer  `json:"layer,omitempty"`
	Error       string    `json:"error,omitempty"`
	ErrorType   string    `json:"error_type,omitempty"`
	ServiceName string    `json:"service"`
	Fields      Fields    `json:"fields,omitempty"`
}

// Factory defines the interface for creating loggers
type Factory interface {
	// Create a new logger instance
	NewLogger() Logger

	// Create layer-specific loggers
	NewDomainLogger() Logger
	NewApplicationLogger() Logger
	NewInfrastructureLogger() Logger
	NewInterfaceLogger() Logger

	// Create component-specific loggers
	NewComponentLogger(component string) Logger
}

// Config represents the logging configuration
type Config struct {
	Level         Level  `yaml:"level" mapstructure:"level"`
	Format        string `yaml:"format" mapstructure:"format"`
	Output        string `yaml:"output" mapstructure:"output"`
	EnableFile    bool   `yaml:"enable_file" mapstructure:"enable_file"`
	FilePath      string `yaml:"file_path" mapstructure:"file_path"`
	ServiceName   string `yaml:"service_name" mapstructure:"service_name"`
	EnableTracing bool   `yaml:"enable_tracing" mapstructure:"enable_tracing"`
	MaxFileSize   int    `yaml:"max_file_size" mapstructure:"max_file_size"`
	MaxBackups    int    `yaml:"max_backups" mapstructure:"max_backups"`
	MaxAge        int    `yaml:"max_age" mapstructure:"max_age"`
	Compress      bool   `yaml:"compress" mapstructure:"compress"`
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}