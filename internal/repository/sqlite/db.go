package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the standard *sql.DB and manages database connections.
type DB struct {
	*sql.DB
}

// Connect establishes a connection to the SQLite database.
func Connect(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// SQLite connection settings (limit open conns for thread safety with SQLite file access)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection is active
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	// Enable foreign key constraints in SQLite
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
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
