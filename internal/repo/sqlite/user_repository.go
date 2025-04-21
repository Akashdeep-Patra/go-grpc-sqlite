package sqlite

import (
	"context"
	"database/sql"
	"time"

	"test-project-grpc/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteUserRepository is a SQLite implementation of the UserRepository interface
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new instance of the SQLite user repository
func NewSQLiteUserRepository(dbPath string) (*SQLiteUserRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ensure connection works
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{
		db: db,
	}, nil
}

// Close closes the database connection
func (r *SQLiteUserRepository) Close() error {
	return r.db.Close()
}

// Create adds a new user to the SQLite database
func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetByID retrieves a user by ID from the SQLite database
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id)

	var user domain.User
	var createdAt, updatedAt string

	err := row.Scan(&user.ID, &user.Name, &user.Email, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse the time strings
	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// Update modifies an existing user in the SQLite database
func (r *SQLiteUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ?
	`, user.Name, user.Email, user.UpdatedAt, user.ID)
	return err
}

// Delete removes a user from the SQLite database
func (r *SQLiteUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE id = ?
	`, id)
	return err
} 