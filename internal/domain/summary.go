package domain

import (
	"context"
	"fmt"
	"math"
)

// MonthlySettlement tracks if a specific month has been settled.
type MonthlySettlement struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	IsSettled   bool    `json:"is_settled"`
	SettledByID *string `json:"settled_by_id"` // user ID who settled it
}

// MonthlySettlementRepository defines database operations for settlements.
type MonthlySettlementRepository interface {
	Get(ctx context.Context, year int, month int) (*MonthlySettlement, error)
	Save(ctx context.Context, set *MonthlySettlement) error
}

// UserSummary represents a single user's calculated totals for the month.
type UserSummary struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Salary              float64 `json:"salary"`
	Proportion          float64 `json:"proportion"`
	TotalPaidShared     float64 `json:"total_paid_shared"`
	TotalPaidIndividual float64 `json:"total_paid_individual"`
	FairShare           float64 `json:"fair_share"`
	Balance             float64 `json:"balance"` // Positive means they owe money, negative means they are owed
}

// MonthlySummary is the aggregated view of a month's finances.
type MonthlySummary struct {
	Month               string             `json:"month"` // YYYY-MM
	IsSettled           bool               `json:"is_settled"`
	SettledByID         *string            `json:"settled_by_id"`
	TotalSharedExpenses float64            `json:"total_shared_expenses"`
	SettlementMessage   string             `json:"settlement_message"`
	Users               []UserSummary      `json:"users"`
	CategoryBreakdown   map[string]float64 `json:"category_breakdown"`
	Expenses            []Expense          `json:"expenses"`
}

// Helper to round float64 to 2 decimal places
func round2(val float64) float64 {
	return math.Round(val*100) / 100
}

// CalculateSummary computes the proportions, totals, and balances for a given month.
func CalculateSummary(year int, month int, settlement *MonthlySettlement, users []User, expenses []Expense) *MonthlySummary {
	monthStr := fmt.Sprintf("%04d-%02d", year, month)

	summary := &MonthlySummary{
		Month:             monthStr,
		Users:             make([]UserSummary, 0, len(users)),
		CategoryBreakdown: make(map[string]float64),
		Expenses:          expenses,
	}

	if settlement != nil {
		summary.IsSettled = settlement.IsSettled
		summary.SettledByID = settlement.SettledByID
	}

	// 1. Calculate Total Salary
	var totalSalary float64
	for _, u := range users {
		totalSalary += u.Salary
	}

	// 2. Sum Expenses
	var totalShared float64
	paidSharedByUser := make(map[string]float64)
	paidIndivByUser := make(map[string]float64)

	for _, exp := range expenses {
		// Category breakdown (for both shared and individual)
		summary.CategoryBreakdown[exp.Category] = round2(summary.CategoryBreakdown[exp.Category] + exp.Amount)

		if exp.IsShared {
			totalShared += exp.Amount
			paidSharedByUser[exp.PayerID] = paidSharedByUser[exp.PayerID] + exp.Amount
		} else {
			paidIndivByUser[exp.PayerID] = paidIndivByUser[exp.PayerID] + exp.Amount
		}
	}
	summary.TotalSharedExpenses = round2(totalShared)

	// 3. Calculate each user's breakdown
	for _, u := range users {
		var proportion float64
		if totalSalary > 0 {
			proportion = u.Salary / totalSalary
		}

		paidShared := round2(paidSharedByUser[u.ID])
		paidIndiv := round2(paidIndivByUser[u.ID])
		fairShare := round2(proportion * totalShared)
		balance := round2(fairShare - paidShared)

		summary.Users = append(summary.Users, UserSummary{
			ID:                  u.ID,
			Name:                u.Name,
			Salary:              u.Salary,
			Proportion:          round2(proportion),
			TotalPaidShared:     paidShared,
			TotalPaidIndividual: paidIndiv,
			FairShare:           fairShare,
			Balance:             balance,
		})
	}

	// 4. Generate Settlement Message
	if len(summary.Users) == 2 {
		u1 := summary.Users[0]
		u2 := summary.Users[1]

		if u1.Balance > 0.005 {
			summary.SettlementMessage = fmt.Sprintf("%s deve R$ %.2f para %s", u1.Name, u1.Balance, u2.Name)
		} else if u2.Balance > 0.005 {
			summary.SettlementMessage = fmt.Sprintf("%s deve R$ %.2f para %s", u2.Name, u2.Balance, u1.Name)
		} else {
			summary.SettlementMessage = "Tudo acertado!"
		}
	} else if len(summary.Users) > 2 {
		// Simple fallback if there are more than 2 users (though the app is family-focused for 2)
		var debtors []UserSummary
		var creditors []UserSummary
		for _, us := range summary.Users {
			if us.Balance > 0.01 {
				debtors = append(debtors, us)
			} else if us.Balance < -0.01 {
				creditors = append(creditors, us)
			}
		}

		if len(debtors) == 0 && len(creditors) == 0 {
			summary.SettlementMessage = "Tudo acertado!"
		} else {
			summary.SettlementMessage = "Pendências financeiras em aberto."
		}
	} else {
		summary.SettlementMessage = "Adicione participantes para calcular divisões."
	}

	return summary
}
