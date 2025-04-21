package domain

import (
	"context"
	"time"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data storage
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
} 