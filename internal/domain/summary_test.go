package domain

import (
	"testing"
)

func TestCalculateSummary(t *testing.T) {
	// 1. Setup mock data
	users := []User{
		{ID: "user-1", Name: "Rafael", Salary: 6000.00},
		{ID: "user-2", Name: "Partner", Salary: 4000.00},
	}

	expenses := []Expense{
		{ID: "exp-1", Description: "Supermarket", Amount: 500.00, Category: CategoryMercado, PayerID: "user-1", IsShared: true, Date: "2026-06-10"},
		{ID: "exp-2", Description: "Dog Food", Amount: 200.00, Category: CategoryDog, PayerID: "user-1", IsShared: true, Date: "2026-06-12"},
		{ID: "exp-3", Description: "Rent", Amount: 300.00, Category: CategoryCasa, PayerID: "user-2", IsShared: true, Date: "2026-06-15"},
		{ID: "exp-4", Description: "Gym Clothes", Amount: 150.00, Category: CategoryLazer, PayerID: "user-1", IsShared: false, Date: "2026-06-15"},
	}

	settlement := &MonthlySettlement{
		Year:        2026,
		Month:       6,
		IsSettled:   false,
		SettledByID: nil,
	}

	// 2. Run calculation
	summary := CalculateSummary(2026, 6, settlement, users, expenses)

	// 3. Assertions
	if summary.Month != "2026-06" {
		t.Errorf("Expected month 2026-06, got %s", summary.Month)
	}

	if summary.TotalSharedExpenses != 1000.00 {
		t.Errorf("Expected total shared expenses 1000.00, got %.2f", summary.TotalSharedExpenses)
	}

	// Find users in summary
	var u1, u2 UserSummary
	for _, u := range summary.Users {
		if u.ID == "user-1" {
			u1 = u
		} else if u.ID == "user-2" {
			u2 = u
		}
	}

	// Rafael (user-1) checks
	if u1.Proportion != 0.60 {
		t.Errorf("Expected Rafael proportion 0.60, got %.2f", u1.Proportion)
	}
	if u1.TotalPaidShared != 700.00 {
		t.Errorf("Expected Rafael paid shared 700.00, got %.2f", u1.TotalPaidShared)
	}
	if u1.TotalPaidIndividual != 150.00 {
		t.Errorf("Expected Rafael paid individual 150.00, got %.2f", u1.TotalPaidIndividual)
	}
	if u1.FairShare != 600.00 {
		t.Errorf("Expected Rafael fair share 600.00, got %.2f", u1.FairShare)
	}
	if u1.Balance != -100.00 {
		t.Errorf("Expected Rafael balance -100.00, got %.2f", u1.Balance)
	}

	// Partner (user-2) checks
	if u2.Proportion != 0.40 {
		t.Errorf("Expected Partner proportion 0.40, got %.2f", u2.Proportion)
	}
	if u2.TotalPaidShared != 300.00 {
		t.Errorf("Expected Partner paid shared 300.00, got %.2f", u2.TotalPaidShared)
	}
	if u2.TotalPaidIndividual != 0.00 {
		t.Errorf("Expected Partner paid individual 0.00, got %.2f", u2.TotalPaidIndividual)
	}
	if u2.FairShare != 400.00 {
		t.Errorf("Expected Partner fair share 400.00, got %.2f", u2.FairShare)
	}
	if u2.Balance != 100.00 {
		t.Errorf("Expected Partner balance 100.00, got %.2f", u2.Balance)
	}

	// Settlement message check
	expectedMsg := "Partner deve R$ 100.00 para Rafael"
	if summary.SettlementMessage != expectedMsg {
		t.Errorf("Expected settlement message '%s', got '%s'", expectedMsg, summary.SettlementMessage)
	}

	// Category breakdown check
	expectedBreakdown := map[string]float64{
		CategoryMercado: 500.00,
		CategoryDog:     200.00,
		CategoryCasa:    300.00,
		CategoryLazer:   150.00,
	}

	for cat, amt := range expectedBreakdown {
		if summary.CategoryBreakdown[cat] != amt {
			t.Errorf("Expected category %s to be %.2f, got %.2f", cat, amt, summary.CategoryBreakdown[cat])
		}
	}
}
