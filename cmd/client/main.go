package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Akashdeep-Patra/go-grpc-sqlite/gen/go/github.com/Akashdeep-Patra/go-grpc-sqlite/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Command-line flags
	serverAddr := flag.String("server", "localhost:50051", "The server address in the format of host:port")
	createUser := flag.Bool("create", false, "Create a new user")
	getUserID := flag.String("get", "", "Get user by ID")
	userName := flag.String("name", "", "User name for create operation")
	userEmail := flag.String("email", "", "User email for create operation")
	flag.Parse()

	// Set up connection to server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new user
	if *createUser {
		if *userName == "" || *userEmail == "" {
			log.Fatal("Name and email are required for user creation")
		}

		resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
			Name:  *userName,
			Email: *userEmail,
		})
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		log.Printf("User created successfully with ID: %s", resp.Id)
		log.Printf("User details: Name=%s, Email=%s", resp.Name, resp.Email)
	}

	// Get user by ID
	if *getUserID != "" {
		resp, err := client.GetUser(ctx, &pb.GetUserRequest{
			Id: *getUserID,
		})
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}

		log.Printf("User found: ID=%s, Name=%s, Email=%s", resp.Id, resp.Name, resp.Email)
	}

	// If no operation was specified
	if !*createUser && *getUserID == "" {
		log.Println("No operation specified. Use --create to create a user or --get=<id> to get a user.")
		log.Println("Example usage:")
		log.Println("  ./client --create --name=\"John Doe\" --email=\"john@example.com\"")
		log.Println("  ./client --get=<user_id>")
	}
} 