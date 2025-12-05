package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new database connection pool.
func NewConnection(dbURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	
	// Test connection
	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	
	return dbpool, nil
}

