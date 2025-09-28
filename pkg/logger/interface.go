package logger

import (
	"context"
	"time"
)

// Logger defines a minimal interface for structured logging
// Inspired by go-kit/log and Kubernetes klog design principles
type Logger interface {
	// Core logging methods - simplified to key-value pairs only
	Debug(ctx context.Context, msg string, keyvals ...interface{})
	Info(ctx context.Context, msg string, keyvals ...interface{})
	Warn(ctx context.Context, msg string, keyvals ...interface{})
	Error(ctx context.Context, msg string, keyvals ...interface{})

	// Level checking for performance optimization
	DebugEnabled() bool
	InfoEnabled() bool

	// Context enrichment - returns new logger with additional context
	With(keyvals ...interface{}) Logger
	WithLayer(layer string) Logger
	WithComponent(component string) Logger
	WithError(err error) Logger
}

// Log represents a single log method interface (go-kit style)
// For maximum flexibility, some use cases may prefer this minimal approach
type Log interface {
	Log(keyvals ...interface{}) error
}

// Level constants for simple usage
const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

// Performance-optimized field creators that avoid allocations
func S(key, value string) interface{}  { return key }
func Sv(key, value string) interface{} { return value }

// KV creates a key-value pair
func KV(key string, value interface{}) []interface{} {
	return []interface{}{key, value}
}

// Merge multiple key-value pairs
func Merge(kvs ...[]interface{}) []interface{} {
	var result []interface{}
	for _, kv := range kvs {
		result = append(result, kv...)
	}
	return result
}

// Common field helpers for frequent use cases
func Err(err error) []interface{} {
	if err == nil {
		return nil
	}
	return []interface{}{"error", err.Error()}
}

func Duration(d time.Duration) []interface{} {
	return []interface{}{"duration", d.String()}
}

func UserID(id string) []interface{} {
	return []interface{}{"user_id", id}
}

func Email(email string) []interface{} {
	return []interface{}{"email", email}
}

func Operation(op string) []interface{} {
	return []interface{}{"operation", op}
}
