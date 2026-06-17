package handler

import (
	"encoding/json"
	"net/http"

	"finduo-ai/internal/domain"
)

type ExpenseHandler struct {
	repo     domain.ExpenseRepository
	userRepo domain.UserRepository
}

// NewExpenseHandler initializes a new ExpenseHandler.
func NewExpenseHandler(repo domain.ExpenseRepository, userRepo domain.UserRepository) *ExpenseHandler {
	return &ExpenseHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

// Create handles POST /api/expenses to register a new expense.
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Category    string  `json:"category"`
		PayerID     string  `json:"payer_id"`
		IsShared    *bool   `json:"is_shared"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Description == "" {
		writeError(w, http.StatusBadRequest, "description is required")
		return
	}
	if input.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be greater than zero")
		return
	}
	if input.Date == "" {
		writeError(w, http.StatusBadRequest, "date is required")
		return
	}
	if !domain.IsValidCategory(input.Category) {
		writeError(w, http.StatusBadRequest, "invalid category")
		return
	}
	if input.PayerID == "" {
		writeError(w, http.StatusBadRequest, "payer_id is required")
		return
	}

	// Verify the payer user exists
	_, err := h.userRepo.Get(r.Context(), input.PayerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "payer not found: "+err.Error())
		return
	}

	isShared := true
	if input.IsShared != nil {
		isShared = *input.IsShared
	}

	exp := domain.Expense{
		Description: input.Description,
		Amount:      input.Amount,
		Date:        input.Date,
		Category:    input.Category,
		PayerID:     input.PayerID,
		IsShared:    isShared,
	}

	if err := h.repo.Create(r.Context(), &exp); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, exp)
}

// Update handles PUT /api/expenses/{id} to modify an expense.
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	existing, err := h.repo.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "expense not found: "+err.Error())
		return
	}

	var input struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Category    string  `json:"category"`
		PayerID     string  `json:"payer_id"`
		IsShared    *bool   `json:"is_shared"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.Amount > 0 {
		existing.Amount = input.Amount
	}
	if input.Date != "" {
		existing.Date = input.Date
	}
	if input.Category != "" {
		if !domain.IsValidCategory(input.Category) {
			writeError(w, http.StatusBadRequest, "invalid category")
			return
		}
		existing.Category = input.Category
	}
	if input.PayerID != "" {
		_, err := h.userRepo.Get(r.Context(), input.PayerID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "payer not found: "+err.Error())
			return
		}
		existing.PayerID = input.PayerID
	}
	if input.IsShared != nil {
		existing.IsShared = *input.IsShared
	}

	if err := h.repo.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// Delete handles DELETE /api/expenses/{id} to remove an expense.
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "expense deleted successfully"})
}
