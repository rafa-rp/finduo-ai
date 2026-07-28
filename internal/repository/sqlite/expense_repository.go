package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"finduo-ai/internal/domain"
)

type ExpenseRepository struct {
	db *DB
}

// NewExpenseRepository creates a new SQLite implementation of ExpenseRepository.
func NewExpenseRepository(db *DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create inserts a new expense into the database.
func (r *ExpenseRepository) Create(ctx context.Context, exp *domain.Expense) error {
	if exp.ID == "" {
		exp.ID = uuid.New().String()
	}

	query := `
		INSERT INTO expenses (id, description, amount, date, category, payer_id, is_shared)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		exp.ID, exp.Description, exp.Amount, exp.Date, exp.Category, exp.PayerID, exp.IsShared,
	)

	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	return nil
}

// Update modifies an existing expense.
func (r *ExpenseRepository) Update(ctx context.Context, exp *domain.Expense) error {
	query := `
		UPDATE expenses
		SET description = ?, amount = ?, date = ?, category = ?, payer_id = ?, is_shared = ?
		WHERE id = ?
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
	query := "DELETE FROM expenses WHERE id = ?"
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
func (r *ExpenseRepository) ListByMonth(ctx context.Context, year int, month int) ([]domain.Expense, error) {
	query := `
		SELECT id, description, amount, date, category, payer_id, is_shared
		FROM expenses
		WHERE strftime('%Y', date) = ? AND strftime('%m', date) = ?
		ORDER BY date ASC, description ASC
	`
	yearStr := fmt.Sprintf("%04d", year)
	monthStr := fmt.Sprintf("%02d", month)

	rows, err := r.db.QueryContext(ctx, query, yearStr, monthStr)
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
		SELECT id, description, amount, date, category, payer_id, is_shared
		FROM expenses
		WHERE id = ?
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
