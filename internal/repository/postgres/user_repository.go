package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"finduo-ai/internal/domain"
)

type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new PostgreSQL implementation of UserRepository.
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// List retrieves all users from the database.
func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, salary FROM users ORDER BY name ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Salary); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

// Save inserts or updates a user in the database.
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		// Insert and let Postgres generate the UUID
		query := "INSERT INTO users (name, salary) VALUES ($1, $2) RETURNING id"
		err := r.db.QueryRowContext(ctx, query, user.Name, user.Salary).Scan(&user.ID)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
		return nil
	}

	// Upsert if ID is provided
	query := `
		INSERT INTO users (id, name, salary)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			salary = EXCLUDED.salary
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Salary)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

// Get retrieves a specific user by ID.
func (r *UserRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	query := "SELECT id, name, salary FROM users WHERE id = $1"
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Salary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}
