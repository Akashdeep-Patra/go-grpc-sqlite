package handler

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"test-project-grpc/pkg/logger"
)

// HealthHandler implements the gRPC Health Checking Protocol
type HealthHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	statusMap map[string]grpc_health_v1.HealthCheckResponse_ServingStatus
	mu        sync.RWMutex
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		statusMap: make(map[string]grpc_health_v1.HealthCheckResponse_ServingStatus),
	}
}

// Check returns the server's health status
func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	service := req.GetService()
	logger.Debug("Health check received", zap.String("service", service))

	// If no service is specified, return the overall server status
	if service == "" {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}

	status, ok := h.statusMap[service]
	if !ok {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}

// Watch returns a stream of health check statuses
func (h *HealthHandler) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	service := req.GetService()
	logger.Debug("Health watch received", zap.String("service", service))

	// This is a simplified implementation that just returns the current status
	// A complete implementation would watch for status changes
	h.mu.RLock()
	defer h.mu.RUnlock()

	// If no service is specified, return the overall server status
	if service == "" {
		return stream.Send(&grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		})
	}

	status, ok := h.statusMap[service]
	if !ok {
		status = grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN
	}

	return stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: status,
	})
}

// SetServingStatus updates the serving status of a service
func (h *HealthHandler) SetServingStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.statusMap[service] = status
	logger.Info("Health status updated",
		zap.String("service", service),
		zap.String("status", status.String()),
	)
} 