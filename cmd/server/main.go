package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/Akashdeep-Patra/go-grpc-sqlite/user"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/internal/handler"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/config"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/logger"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/metrics"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.App.Environment)
	defer logger.Sync()

	// Log basic information
	logger.Info("Starting service",
		zap.String("app", cfg.App.Name),
		zap.String("environment", cfg.App.Environment),
	)

	// Start Prometheus metrics server
	metrics.StartMetricsServer(9100)

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(cfg.Server.IdleTimeout) * time.Second,
			MaxConnectionAge:      time.Hour,
			MaxConnectionAgeGrace: time.Minute * 5,
			Time:                  time.Minute,
			Timeout:               time.Second * 20,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Second * 5,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(),
			middleware.LoggingInterceptor(),
			// Uncomment when you have authentication set up
			// middleware.AuthInterceptor(),
			middleware.RateLimitInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			middleware.RecoveryStreamInterceptor(),
			middleware.LoggingStreamInterceptor(),
			// Uncomment when you have authentication set up
			// middleware.AuthStreamInterceptor(),
			middleware.RateLimitStreamInterceptor(),
		),
	)

	// Create TCP listener
	port := cfg.Server.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// Create health check service
	healthHandler := handler.NewHealthHandler()
	healthHandler.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)

	// Register service handlers
	userHandler := handler.NewUserHandler()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	
	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Set service as serving
	healthHandler.SetServingStatus(cfg.App.Name, grpc_health_v1.HealthCheckResponse_SERVING)

	// Start server in a goroutine
	go func() {
		logger.Info("Server listening", zap.Int("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down server...")
	
	// Set service to NOT_SERVING during shutdown
	healthHandler.SetServingStatus(cfg.App.Name, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	
	// Stop accepting new requests
	grpcServer.GracefulStop()
	
	// Close database connection
	if err := userHandler.Close(); err != nil {
		logger.Error("Error closing DB connection", zap.Error(err))
	}
	
	// Allow some time for existing requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		logger.Info("Timeout waiting for connections to close, forcing shutdown")
	case <-time.After(10 * time.Millisecond):
		logger.Info("All connections closed successfully")
	}
	
	logger.Info("Server stopped")
} 