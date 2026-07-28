package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"finduo-ai/internal/domain"
)

type SummaryHandler struct {
	userRepo       domain.UserRepository
	expenseRepo    domain.ExpenseRepository
	settlementRepo domain.MonthlySettlementRepository
}

// NewSummaryHandler initializes a new SummaryHandler.
func NewSummaryHandler(
	userRepo domain.UserRepository,
	expenseRepo domain.ExpenseRepository,
	settlementRepo domain.MonthlySettlementRepository,
) *SummaryHandler {
	return &SummaryHandler{
		userRepo:       userRepo,
		expenseRepo:    expenseRepo,
		settlementRepo: settlementRepo,
	}
}

// Get handles GET /api/summary?month=YYYY-MM.
// It returns the combined monthly status (proportions, totals, balances, category breakdown, and expenses).
func (h *SummaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	monthParam := r.URL.Query().Get("month")
	if len(monthParam) != 7 || monthParam[4] != '-' {
		writeError(w, http.StatusBadRequest, "month parameter must be in YYYY-MM format")
		return
	}

	var year, month int
	_, err := fmt.Sscanf(monthParam, "%d-%d", &year, &month)
	if err != nil || month < 1 || month > 12 {
		writeError(w, http.StatusBadRequest, "invalid month format")
		return
	}

	ctx := r.Context()

	// 1. Fetch Users
	users, err := h.userRepo.List(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch users: "+err.Error())
		return
	}

	// 2. Fetch Expenses for the month
	expenses, err := h.expenseRepo.ListByMonth(ctx, year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch expenses: "+err.Error())
		return
	}
	if expenses == nil {
		expenses = []domain.Expense{}
	}

	// 3. Fetch Settlement Status
	settlement, err := h.settlementRepo.Get(ctx, year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch settlement status: "+err.Error())
		return
	}

	// 4. Calculate Summary
	summary := domain.CalculateSummary(year, month, settlement, users, expenses)

	writeJSON(w, http.StatusOK, summary)
}

// Settle handles POST /api/summary/settle to toggle the settled status of a month.
func (h *SummaryHandler) Settle(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Year        int     `json:"year"`
		Month       int     `json:"month"`
		IsSettled   bool    `json:"is_settled"`
		SettledByID *string `json:"settled_by_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Year < 1970 || input.Year > 2100 {
		writeError(w, http.StatusBadRequest, "invalid year")
		return
	}
	if input.Month < 1 || input.Month > 12 {
		writeError(w, http.StatusBadRequest, "invalid month")
		return
	}

	ctx := r.Context()

	// If settling, optionally check if the settling user exists
	if input.IsSettled && input.SettledByID != nil && *input.SettledByID != "" {
		_, err := h.userRepo.Get(ctx, *input.SettledByID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "settling user not found: "+err.Error())
			return
		}
	}

	set := domain.MonthlySettlement{
		Year:        input.Year,
		Month:       input.Month,
		IsSettled:   input.IsSettled,
		SettledByID: input.SettledByID,
	}

	if err := h.settlementRepo.Save(ctx, &set); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save settlement status: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, set)
}
