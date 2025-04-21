package app

import (
	"context"
	"time"

	"test-project-grpc/internal/domain"
	"github.com/google/uuid"
)

type userService struct {
	repo domain.UserRepository
}

// NewUserService creates a new instance of the user service
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser implements the domain.UserService interface
func (s *userService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	user := &domain.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser implements the domain.UserService interface
func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
} 