package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	inframetrics "github.com/cctw-zed/wonder/internal/infrastructure/metrics"
)

// MetricsMiddleware records Prometheus metrics for incoming HTTP requests.
func MetricsMiddleware() gin.HandlerFunc {
	inframetrics.EnsureHTTPMetrics()

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start).Seconds()
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		inframetrics.ObserveHTTPRequest(
			c.Request.Method,
			route,
			strconv.Itoa(c.Writer.Status()),
			duration,
		)
	}
}
