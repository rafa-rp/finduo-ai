package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"finduo-ai/internal/domain"
)

type SettlementRepository struct {
	db *DB
}

// NewSettlementRepository creates a new PostgreSQL implementation of MonthlySettlementRepository.
func NewSettlementRepository(db *DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

// Get retrieves the settlement status for a specific month.
// If it does not exist in the database, it returns a default unsettled record.
func (r *SettlementRepository) Get(ctx context.Context, year int, month int) (*domain.MonthlySettlement, error) {
	query := "SELECT year, month, is_settled, settled_by_id FROM monthly_settlements WHERE year = $1 AND month = $2"

	var set domain.MonthlySettlement
	err := r.db.QueryRowContext(ctx, query, year, month).Scan(&set.Year, &set.Month, &set.IsSettled, &set.SettledByID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Not yet created in DB, return default unsettled state
			return &domain.MonthlySettlement{
				Year:        year,
				Month:       month,
				IsSettled:   false,
				SettledByID: nil,
			}, nil
		}
		return nil, fmt.Errorf("failed to get monthly settlement: %w", err)
	}

	return &set, nil
}

// Save inserts or updates the settlement status for a month.
func (r *SettlementRepository) Save(ctx context.Context, set *domain.MonthlySettlement) error {
	query := `
		INSERT INTO monthly_settlements (year, month, is_settled, settled_by_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (year, month) DO UPDATE SET
			is_settled = EXCLUDED.is_settled,
			settled_by_id = EXCLUDED.settled_by_id
	`
	_, err := r.db.ExecContext(ctx, query, set.Year, set.Month, set.IsSettled, set.SettledByID)
	if err != nil {
		return fmt.Errorf("failed to save monthly settlement: %w", err)
	}

	return nil
}
