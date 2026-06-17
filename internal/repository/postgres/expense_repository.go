package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"finduo-ai/internal/domain"
)

type ExpenseRepository struct {
	db *DB
}

// NewExpenseRepository creates a new PostgreSQL implementation of ExpenseRepository.
func NewExpenseRepository(db *DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create inserts a new expense into the database and returns the generated UUID.
func (r *ExpenseRepository) Create(ctx context.Context, exp *domain.Expense) error {
	query := `
		INSERT INTO expenses (description, amount, date, category, payer_id, is_shared)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		exp.Description, exp.Amount, exp.Date, exp.Category, exp.PayerID, exp.IsShared,
	).Scan(&exp.ID)

	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	return nil
}

// Update modifies an existing expense.
func (r *ExpenseRepository) Update(ctx context.Context, exp *domain.Expense) error {
	query := `
		UPDATE expenses
		SET description = $1, amount = $2, date = $3, category = $4, payer_id = $5, is_shared = $6
		WHERE id = $7
	`
	res, err := r.db.ExecContext(
		ctx, query,
		exp.Description, exp.Amount, exp.Date, exp.Category, exp.PayerID, exp.IsShared, exp.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("expense not found: %s", exp.ID)
	}

	return nil
}

// Delete removes an expense from the database.
func (r *ExpenseRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM expenses WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("expense not found: %s", id)
	}

	return nil
}

// ListByMonth retrieves all expenses for a specific month and year.
// Uses date::text to get YYYY-MM-DD directly.
func (r *ExpenseRepository) ListByMonth(ctx context.Context, year int, month int) ([]domain.Expense, error) {
	query := `
		SELECT id, description, amount, date::text, category, payer_id, is_shared
		FROM expenses
		WHERE EXTRACT(YEAR FROM date) = $1 AND EXTRACT(MONTH FROM date) = $2
		ORDER BY date ASC, description ASC
	`
	rows, err := r.db.QueryContext(ctx, query, year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var exp domain.Expense
		err := rows.Scan(
			&exp.ID, &exp.Description, &exp.Amount, &exp.Date,
			&exp.Category, &exp.PayerID, &exp.IsShared,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		expenses = append(expenses, exp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return expenses, nil
}

// Get retrieves a specific expense by ID.
func (r *ExpenseRepository) Get(ctx context.Context, id string) (*domain.Expense, error) {
	query := `
		SELECT id, description, amount, date::text, category, payer_id, is_shared
		FROM expenses
		WHERE id = $1
	`
	var exp domain.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exp.ID, &exp.Description, &exp.Amount, &exp.Date,
		&exp.Category, &exp.PayerID, &exp.IsShared,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("expense not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}
	return &exp, nil
}
