package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"test-project-grpc/internal/handler"
	pb "test-project-grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	
	// Register service handlers
	userHandler := handler.NewUserHandler()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	
	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down server...")
	
	// Close database connection
	if err := userHandler.Close(); err != nil {
		log.Printf("Error closing DB connection: %v", err)
	}
	
	grpcServer.GracefulStop()
	log.Println("Server stopped")
} 