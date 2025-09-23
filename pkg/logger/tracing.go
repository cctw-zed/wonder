package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// TraceIDKey is the context key for trace ID
	TraceIDKey ContextKey = "trace_id"
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// CorrelationIDKey is the context key for correlation ID
	CorrelationIDKey ContextKey = "correlation_id"
)

// TraceContext holds tracing information
type TraceContext struct {
	TraceID       string    `json:"trace_id"`
	RequestID     string    `json:"request_id"`
	CorrelationID string    `json:"correlation_id"`
	StartTime     time.Time `json:"start_time"`
	UserID        string    `json:"user_id,omitempty"`
	SessionID     string    `json:"session_id,omitempty"`
}

// GenerateTraceID generates a new trace ID
func GenerateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("trace_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("req_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// NewTraceContext creates a new trace context
func NewTraceContext() *TraceContext {
	return &TraceContext{
		TraceID:   GenerateTraceID(),
		RequestID: GenerateRequestID(),
		StartTime: time.Now(),
	}
}

// WithTraceContext adds trace context to the context
func WithTraceContext(ctx context.Context, traceCtx *TraceContext) context.Context {
	ctx = context.WithValue(ctx, TraceIDKey, traceCtx.TraceID)
	ctx = context.WithValue(ctx, RequestIDKey, traceCtx.RequestID)
	if traceCtx.CorrelationID != "" {
		ctx = context.WithValue(ctx, CorrelationIDKey, traceCtx.CorrelationID)
	}
	return ctx
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithCorrelationID adds a correlation ID to the context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}

	// Fallback to checking other common keys
	if traceID := extractTraceID(ctx); traceID != "" {
		return traceID
	}

	return ""
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if correlationID, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return correlationID
	}

	return ""
}

// GetTraceContext extracts complete trace context from context
func GetTraceContext(ctx context.Context) *TraceContext {
	if ctx == nil {
		return nil
	}

	traceCtx := &TraceContext{
		TraceID:       GetTraceID(ctx),
		RequestID:     GetRequestID(ctx),
		CorrelationID: GetCorrelationID(ctx),
	}

	// Only return if we have at least a trace ID
	if traceCtx.TraceID == "" {
		return nil
	}

	return traceCtx
}

// LoggerWithTracing creates a logger with tracing information from context
func LoggerWithTracing(ctx context.Context, baseLogger Logger) Logger {
	logger := baseLogger

	if traceID := GetTraceID(ctx); traceID != "" {
		logger = logger.WithTraceID(traceID)
	}

	if requestID := GetRequestID(ctx); requestID != "" {
		logger = logger.WithField("request_id", requestID)
	}

	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		logger = logger.WithField("correlation_id", correlationID)
	}

	return logger
}

// PerformanceLogger helps with performance logging
type PerformanceLogger struct {
	logger Logger
	ctx    context.Context
}

// NewPerformanceLogger creates a new performance logger
func NewPerformanceLogger(ctx context.Context, logger Logger) *PerformanceLogger {
	return &PerformanceLogger{
		logger: LoggerWithTracing(ctx, logger),
		ctx:    ctx,
	}
}

// LogDuration logs the duration of an operation
func (p *PerformanceLogger) LogDuration(operation string, startTime time.Time, fields ...Field) {
	duration := time.Since(startTime)

	allFields := append(fields,
		String("operation", operation),
		Duration("duration", duration),
		Duration("duration_ms", duration/time.Millisecond),
	)

	// Log as warning if operation is slow (>1s), info otherwise
	if duration > time.Second {
		p.logger.Warn(p.ctx, fmt.Sprintf("Slow operation: %s", operation), allFields...)
	} else {
		p.logger.Info(p.ctx, fmt.Sprintf("Operation completed: %s", operation), allFields...)
	}
}

// LogError logs an operation error with duration
func (p *PerformanceLogger) LogError(operation string, startTime time.Time, err error, fields ...Field) {
	duration := time.Since(startTime)

	allFields := append(fields,
		String("operation", operation),
		Duration("duration", duration),
		Duration("duration_ms", duration/time.Millisecond),
		String("error", err.Error()),
	)

	p.logger.Error(p.ctx, fmt.Sprintf("Operation failed: %s", operation), allFields...)
}