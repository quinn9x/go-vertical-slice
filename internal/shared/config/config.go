// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultAppPort = 9080
	defaultDBPort  = 5432
)

// Config contains the application and database configuration.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

// AppConfig contains the application server configuration.
type AppConfig struct {
	Host string
	Port int
}

// DatabaseConfig contains the database connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// Load loads the application configuration from environment variables.
func Load() (Config, error) {
	// Load .env if it exists.
	// Environment variables already set take precedence.
	_ = godotenv.Load()

	port, err := getInt("APP_PORT", defaultAppPort)
	if err != nil {
		return Config{}, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	dbPort, err := getInt("DB_PORT", defaultDBPort)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return Config{
		App: AppConfig{
			Host: getString("APP_HOST", "0.0.0.0"),
			Port: port,
		},
		Database: DatabaseConfig{
			Host:     getString("DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     getString("DB_NAME", "vertical_slice"),
			User:     getString("DB_USER", "postgres"),
			Password: getString("DB_PASSWORD", "postgres"),
			SSLMode:  getString("DB_SSL_MODE", "disable"),
		},
	}, nil
}

// DSN returns the database connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

func getInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	return strconv.Atoi(value)
}

func getString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
