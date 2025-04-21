package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	
	"test-project-grpc/pkg/logger"
)

var (
	// RequestCounter counts the number of requests received by method
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "The total number of gRPC requests received",
		},
		[]string{"method", "status"},
	)

	// RequestDuration tracks request duration by method
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "The gRPC request latencies in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// ActiveRequests tracks the number of requests in flight
	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_active_requests",
			Help: "The number of active gRPC requests",
		},
		[]string{"method"},
	)

	// ErrorCounter tracks the number of errors by method and code
	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_errors_total",
			Help: "The total number of errors in gRPC requests",
		},
		[]string{"method", "code"},
	)
)

// StartMetricsServer starts an HTTP server for Prometheus metrics
func StartMetricsServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	
	addr := fmt.Sprintf(":%d", port)
	go func() {
		logger.Info("Starting metrics server", zap.String("address", addr))
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Error("Metrics server error", zap.Error(err))
		}
	}()
} 