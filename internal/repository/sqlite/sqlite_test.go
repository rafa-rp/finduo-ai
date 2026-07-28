package sqlite

import (
	"context"
	"testing"
	"time"

	"finduo-ai/internal/domain"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()
	// Use in-memory SQLite database for fast unit testing
	db, err := Connect(":memory:")
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.InitSchema(ctx); err != nil {
		db.Close()
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := NewUserRepository(db)

	// Save new user
	u1 := &domain.User{Name: "Rafael", Salary: 5000.0}
	if err := repo.Save(ctx, u1); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}
	if u1.ID == "" {
		t.Errorf("expected generated UUID for user, got empty string")
	}

	// Get user
	fetched, err := repo.Get(ctx, u1.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if fetched.Name != "Rafael" || fetched.Salary != 5000.0 {
		t.Errorf("unexpected user data: %+v", fetched)
	}

	// List users
	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestExpenseRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	userRepo := NewUserRepository(db)
	expenseRepo := NewExpenseRepository(db)

	// Create user
	u1 := &domain.User{Name: "Partner A", Salary: 3000.0}
	if err := userRepo.Save(ctx, u1); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	// Create expense
	exp := &domain.Expense{
		Description: "Supermarket",
		Amount:      150.50,
		Date:        "2026-07-28",
		Category:    "Groceries",
		PayerID:     u1.ID,
		IsShared:    true,
	}
	if err := expenseRepo.Create(ctx, exp); err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}
	if exp.ID == "" {
		t.Errorf("expected generated UUID for expense")
	}

	// List by month
	expenses, err := expenseRepo.ListByMonth(ctx, 2026, 7)
	if err != nil {
		t.Fatalf("failed to list expenses by month: %v", err)
	}
	if len(expenses) != 1 {
		t.Errorf("expected 1 expense for 2026-07, got %d", len(expenses))
	}

	// Update expense
	exp.Amount = 180.00
	if err := expenseRepo.Update(ctx, exp); err != nil {
		t.Fatalf("failed to update expense: %v", err)
	}

	updated, err := expenseRepo.Get(ctx, exp.ID)
	if err != nil {
		t.Fatalf("failed to get expense: %v", err)
	}
	if updated.Amount != 180.00 {
		t.Errorf("expected updated amount 180.00, got %f", updated.Amount)
	}

	// Delete expense
	if err := expenseRepo.Delete(ctx, exp.ID); err != nil {
		t.Fatalf("failed to delete expense: %v", err)
	}
}

func TestSettlementRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	userRepo := NewUserRepository(db)
	settleRepo := NewSettlementRepository(db)

	u := &domain.User{Name: "Rafael", Salary: 4000.0}
	if err := userRepo.Save(ctx, u); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	// Get initial default state
	st, err := settleRepo.Get(ctx, 2026, 7)
	if err != nil {
		t.Fatalf("failed to get settlement: %v", err)
	}
	if st.IsSettled {
		t.Errorf("expected unsettled by default")
	}

	// Save settlement
	st.IsSettled = true
	st.SettledByID = &u.ID
	if err := settleRepo.Save(ctx, st); err != nil {
		t.Fatalf("failed to save settlement: %v", err)
	}

	// Fetch again
	updated, err := settleRepo.Get(ctx, 2026, 7)
	if err != nil {
		t.Fatalf("failed to get updated settlement: %v", err)
	}
	if !updated.IsSettled || updated.SettledByID == nil || *updated.SettledByID != u.ID {
		t.Errorf("unexpected settlement status: %+v", updated)
	}
}
