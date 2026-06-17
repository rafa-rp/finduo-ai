package domain

import "context"

// User represents a participant in the system with their monthly salary.
type User struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
}

// UserRepository defines the database operations for User entities.
type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	Save(ctx context.Context, user *User) error
	Get(ctx context.Context, id string) (*User, error)
}
