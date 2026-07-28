package domain

import "context"

// Valid categories as specified by the user
const (
	CategoryCasa    = "Casa"
	CategoryCarro   = "Carro"
	CategoryDog     = "Dog"
	CategoryMercado = "Mercado"
	CategoryViagem  = "Viagem"
	CategoryLazer   = "Lazer"
	CategorySaude   = "Saúde"
	CategoryOutros  = "Outros"
)

// IsValidCategory checks if a given category name is valid.
func IsValidCategory(cat string) bool {
	switch cat {
	case CategoryCasa, CategoryCarro, CategoryDog, CategoryMercado, CategoryViagem, CategoryLazer, CategorySaude, CategoryOutros:
		return true
	}
	return false
}

// Expense represents a single expense record, which can be individual or shared.
type Expense struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"` // YYYY-MM-DD format
	Category    string  `json:"category"`
	PayerID     string  `json:"payer_id"`
	IsShared    bool    `json:"is_shared"`
}

// ExpenseRepository defines the database operations for Expense entities.
type ExpenseRepository interface {
	Create(ctx context.Context, exp *Expense) error
	Update(ctx context.Context, exp *Expense) error
	Delete(ctx context.Context, id string) error
	ListByMonth(ctx context.Context, year int, month int) ([]Expense, error)
	Get(ctx context.Context, id string) (*Expense, error)
}
