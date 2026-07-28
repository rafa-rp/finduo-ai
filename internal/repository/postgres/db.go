package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the standard *sql.DB and manages database connections.
type DB struct {
	*sql.DB
}

// Connect establishes a connection pool to the PostgreSQL database.
func Connect(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool configurations
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection is active
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// InitSchema runs the DDL schema statement to ensure all tables exist.
func (db *DB) InitSchema(ctx context.Context) error {
	_, err := db.ExecContext(ctx, SchemaSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return nil
}
