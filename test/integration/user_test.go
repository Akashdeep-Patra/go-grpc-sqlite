// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Akashdeep-Patra/go-grpc-sqlite/gen/go/github.com/Akashdeep-Patra/go-grpc-sqlite/user"
)

const (
	serverAddress = "localhost:50051"
	timeout       = time.Second * 5
)

func TestCreateAndGetUser(t *testing.T) {
	// Set up a connection to the server
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to connect to server")
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Create a user
	createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Integration Test User",
		Email: "integration-test@example.com",
	})
	require.NoError(t, err, "Failed to create user")
	assert.NotEmpty(t, createResp.Id, "User ID should not be empty")
	assert.Equal(t, "Integration Test User", createResp.Name, "User name should match")
	assert.Equal(t, "integration-test@example.com", createResp.Email, "User email should match")

	// Get the user
	getResp, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: createResp.Id,
	})
	require.NoError(t, err, "Failed to get user")
	assert.Equal(t, createResp.Id, getResp.Id, "User ID should match")
	assert.Equal(t, createResp.Name, getResp.Name, "User name should match")
	assert.Equal(t, createResp.Email, getResp.Email, "User email should match")
} 