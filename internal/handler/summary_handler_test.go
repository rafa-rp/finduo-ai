package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"finduo-ai/internal/domain"
)

type mockSettlementRepo struct {
	settlement *domain.MonthlySettlement
	err        error
}

func (m *mockSettlementRepo) Get(ctx context.Context, year int, month int) (*domain.MonthlySettlement, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.settlement == nil {
		return &domain.MonthlySettlement{Year: year, Month: month, IsSettled: false}, nil
	}
	return m.settlement, nil
}

func (m *mockSettlementRepo) Save(ctx context.Context, s *domain.MonthlySettlement) error {
	if m.err != nil {
		return m.err
	}
	m.settlement = s
	return nil
}

func TestSummaryHandler_Get(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]domain.User{
			"u1": {ID: "u1", Name: "Rafael", Salary: 6000},
			"u2": {ID: "u2", Name: "Partner", Salary: 4000},
		},
	}
	expRepo := &mockExpenseRepo{
		expenses: map[string]domain.Expense{
			"e1": {ID: "e1", Description: "Market", Amount: 500, Category: "Mercado", PayerID: "u1", IsShared: true, Date: "2026-06-10"},
		},
	}
	settleRepo := &mockSettlementRepo{}

	t.Run("Success", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, expRepo, settleRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026-06", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Missing or invalid month format length", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, expRepo, settleRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid month number", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, expRepo, settleRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026-13", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("User repo error", func(t *testing.T) {
		h := NewSummaryHandler(&mockUserRepo{err: errors.New("user error")}, expRepo, settleRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026-06", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})

	t.Run("Expense repo error", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, &mockExpenseRepo{err: errors.New("exp error")}, settleRepo)

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026-06", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})

	t.Run("Settlement repo error", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, expRepo, &mockSettlementRepo{err: errors.New("settle error")})

		req := httptest.NewRequest(http.MethodGet, "/api/summary?month=2026-06", nil)
		rec := httptest.NewRecorder()

		h.Get(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}

func TestSummaryHandler_Settle(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]domain.User{
			"u1": {ID: "u1", Name: "Rafael", Salary: 6000},
		},
	}
	settleRepo := &mockSettlementRepo{}

	t.Run("Success", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, settleRepo)

		userID := "u1"
		body := `{"year":2026,"month":6,"is_settled":true,"settled_by_id":"u1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString(body))
		_ = userID
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Invalid body", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, settleRepo)

		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString("invalid json"))
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid year", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, settleRepo)

		body := `{"year":1800,"month":6}`
		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid month", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, settleRepo)

		body := `{"year":2026,"month":15}`
		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Settling user not found", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, settleRepo)

		body := `{"year":2026,"month":6,"is_settled":true,"settled_by_id":"unknown"}`
		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Settlement repo save error", func(t *testing.T) {
		h := NewSummaryHandler(userRepo, nil, &mockSettlementRepo{err: errors.New("save error")})

		body := `{"year":2026,"month":6,"is_settled":true}`
		req := httptest.NewRequest(http.MethodPost, "/api/summary/settle", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Settle(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}
