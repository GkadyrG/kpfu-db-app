package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	DBURL string
}

// Load returns a new Config struct with values from environment variables.
func Load() *Config {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/shipments_db?sslmode=disable"
	}
	return &Config{
		DBURL: dbURL,
	}
}

