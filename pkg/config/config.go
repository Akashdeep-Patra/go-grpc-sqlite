package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	Host         string `mapstructure:"host"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	SQLiteDBPath string `mapstructure:"sqlite_db_path"`
}

// Load loads the configuration from files and environment variables
func Load(configPaths ...string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/app/")

	// Add all provided config paths
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %s", err)
		}
	}

	// Set default values
	setDefaults(v)

	// Override config from environment variables
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Create a new config struct
	cfg := &Config{}

	// Unmarshal the config into the struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %s", err)
	}

	// Ensure database path exists
	if cfg.Database.SQLiteDBPath != "" {
		dir := filepath.Dir(cfg.Database.SQLiteDBPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("unable to create database directory: %s", err)
		}
	}

	return cfg, nil
}

// setDefaults sets the default values for configuration
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "go-grpc-user-service")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.log_level", "info")

	// Server defaults
	v.SetDefault("server.port", 50051)
	v.SetDefault("server.read_timeout", 5)
	v.SetDefault("server.write_timeout", 10)
	v.SetDefault("server.idle_timeout", 15)
	v.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	v.SetDefault("database.sqlite_db_path", "./data/users.db")
} 