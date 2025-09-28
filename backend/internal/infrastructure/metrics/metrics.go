package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	registerOnce        sync.Once
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
)

func initDefault() {
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wonder",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests processed, labeled by method, route, and status code.",
	}, []string{"method", "route", "status"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "wonder",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of latencies for HTTP requests in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "route"})

	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// EnsureHTTPMetrics registers the default HTTP metrics once per process.
func EnsureHTTPMetrics() {
	registerOnce.Do(initDefault)
}

// ObserveHTTPRequest records metrics for a single HTTP request.
func ObserveHTTPRequest(method, route, status string, durationSeconds float64) {
	EnsureHTTPMetrics()
	httpRequestsTotal.WithLabelValues(method, route, status).Inc()
	httpRequestDuration.WithLabelValues(method, route).Observe(durationSeconds)
}
