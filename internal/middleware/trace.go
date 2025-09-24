package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDKey is the context key for storing trace ID
	TraceIDKey = "trace_id"
	// TraceIDHeader is the HTTP header name for trace ID
	TraceIDHeader = "X-Trace-ID"
)

// TraceIDMiddleware creates a middleware that automatically generates and injects
// a TraceID into the request context for distributed tracing and logging
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceID string

		// First, check if trace ID is provided in the request header
		if headerTraceID := c.GetHeader(TraceIDHeader); headerTraceID != "" {
			traceID = headerTraceID
		} else {
			// Generate a new UUID for trace ID if not provided
			traceID = uuid.New().String()
		}

		// Set trace ID in response header for client visibility
		c.Header(TraceIDHeader, traceID)

		// Inject trace ID into the request context
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// Continue with the next handler
		c.Next()
	}
}

// GetTraceIDFromContext extracts trace ID from context
// This is a convenience function for manual trace ID extraction if needed
func GetTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if str, ok := traceID.(string); ok {
			return str
		}
	}

	return ""
}