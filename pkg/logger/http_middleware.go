package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPLoggerMiddleware creates Gin middleware for HTTP request logging
func HTTPLoggerMiddleware(logger Logger) gin.HandlerFunc {
	interfaceLogger := NewInterfaceLogger(logger.WithComponent("http"))

	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Generate trace ID if not present
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = GenerateTraceID()
		}

		// Add trace context to request context
		ctx := WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// Add trace ID to response headers
		c.Header("X-Trace-ID", traceID)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log the request
		interfaceLogger.LogHTTPRequest(
			ctx,
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			duration,
			c.GetHeader("User-Agent"),
			String("remote_addr", c.ClientIP()),
			String("content_length", c.GetHeader("Content-Length")),
			Int("response_size", c.Writer.Size()),
		)
	}
}

// HTTPErrorMiddleware creates Gin middleware for HTTP error logging
func HTTPErrorMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			ctx := c.Request.Context()
			traceID := GetTraceID(ctx)

			errorLogger := NewErrorLogger(logger)

			for _, ginErr := range c.Errors {
				errorLogger.LogError(ctx, ginErr.Err, traceID, map[string]interface{}{
					"method":     c.Request.Method,
					"path":       c.Request.RequestURI,
					"status":     c.Writer.Status(),
					"remote_ip":  c.ClientIP(),
					"user_agent": c.GetHeader("User-Agent"),
				})
			}
		}
	}
}

// RecoveryMiddleware creates Gin middleware for panic recovery with logging
func RecoveryMiddleware(logger Logger) gin.HandlerFunc {
	interfaceLogger := NewInterfaceLogger(logger.WithComponent("recovery"))

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ctx := c.Request.Context()
		traceID := GetTraceID(ctx)

		interfaceLogger.Error(ctx, "Panic recovered",
			String("method", c.Request.Method),
			String("path", c.Request.RequestURI),
			String("remote_ip", c.ClientIP()),
			String("user_agent", c.GetHeader("User-Agent")),
			String("trace_id", traceID),
			Any("panic", recovered),
		)

		c.AbortWithStatus(500)
	})
}