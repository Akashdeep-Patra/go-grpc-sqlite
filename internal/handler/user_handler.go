package handler

import (
	"context"

	"github.com/Akashdeep-Patra/go-grpc-sqlite/internal/app"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/internal/domain"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/internal/repo/sqlite"
	pb "github.com/Akashdeep-Patra/go-grpc-sqlite/user"
	"github.com/Akashdeep-Patra/go-grpc-sqlite/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements the UserService gRPC service
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service domain.UserService
	repo    *sqlite.SQLiteUserRepository
}

// NewUserHandler creates a new instance of the user gRPC handler
func NewUserHandler() *UserHandler {
	dbPath := db.GetSQLiteDBPath()
	
	// Initialize SQLite repository
	repo, err := sqlite.NewSQLiteUserRepository(dbPath)
	if err != nil {
		panic(err)
	}
	
	// Initialize user service with the repository
	service := app.NewUserService(repo)
	
	return &UserHandler{
		service: service,
		repo:    repo,
	}
}

// Close closes the repository connection
func (h *UserHandler) Close() error {
	return h.repo.Close()
}

// CreateUser handles the CreateUser RPC call
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := h.service.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// GetUser handles the GetUser RPC call
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
} 