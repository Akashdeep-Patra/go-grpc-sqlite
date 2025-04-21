package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Parse flags
	service := flag.String("service", "", "The service name to check (default is empty, which checks the server's overall health)")
	addr := flag.String("addr", "localhost:50051", "The server address to check")
	timeout := flag.Duration("timeout", time.Second*3, "The timeout for the health check")
	flag.Parse()

	// Set up a connection to the server
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a health client
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Perform health check
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: *service,
	})
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		os.Exit(1)
	}

	// Check status
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		fmt.Printf("Service is not serving: %s\n", resp.Status.String())
		os.Exit(1)
	}

	fmt.Println("Service is healthy")
	os.Exit(0)
} 