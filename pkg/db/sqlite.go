package db

import (
	"log"
	"os"
	"path/filepath"
)

// GetSQLiteDBPath returns the path to the SQLite database file
func GetSQLiteDBPath() string {
	// Check if DB_PATH env var is set
	dbPath := os.Getenv("DB_PATH")
	if dbPath != "" {
		return dbPath
	}

	// Use default path
	dataDir := "./data"
	
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Failed to create data directory: %v", err)
		// Fall back to current directory
		dataDir = "."
	}

	return filepath.Join(dataDir, "users.db")
} 