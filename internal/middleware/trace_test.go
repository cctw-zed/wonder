package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceIDMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("generates trace ID when none provided", func(t *testing.T) {
		// Create a test router with the middleware
		router := gin.New()
		router.Use(TraceIDMiddleware())

		// Add a test endpoint that extracts trace ID from context
		var capturedTraceID string
		router.GET("/test", func(c *gin.Context) {
			capturedTraceID = GetTraceIDFromContext(c.Request.Context())
			c.JSON(http.StatusOK, gin.H{"trace_id": capturedTraceID})
		})

		// Make a request without trace ID header
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, capturedTraceID)
		assert.Equal(t, capturedTraceID, w.Header().Get(TraceIDHeader))
	})

	t.Run("uses provided trace ID from header", func(t *testing.T) {
		// Create a test router with the middleware
		router := gin.New()
		router.Use(TraceIDMiddleware())

		// Add a test endpoint that extracts trace ID from context
		var capturedTraceID string
		router.GET("/test", func(c *gin.Context) {
			capturedTraceID = GetTraceIDFromContext(c.Request.Context())
			c.JSON(http.StatusOK, gin.H{"trace_id": capturedTraceID})
		})

		// Make a request with existing trace ID header
		providedTraceID := "test-trace-id-123"
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set(TraceIDHeader, providedTraceID)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, providedTraceID, capturedTraceID)
		assert.Equal(t, providedTraceID, w.Header().Get(TraceIDHeader))
	})

	t.Run("trace ID header is always set in response", func(t *testing.T) {
		// Create a test router with the middleware
		router := gin.New()
		router.Use(TraceIDMiddleware())

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Make a request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get(TraceIDHeader))
	})
}

func TestGetTraceIDFromContext(t *testing.T) {
	t.Run("returns empty string for nil context", func(t *testing.T) {
		traceID := GetTraceIDFromContext(nil)
		assert.Empty(t, traceID)
	})

	t.Run("returns empty string when no trace ID in context", func(t *testing.T) {
		ctx := httptest.NewRequest("GET", "/", nil).Context()
		traceID := GetTraceIDFromContext(ctx)
		assert.Empty(t, traceID)
	})
}
