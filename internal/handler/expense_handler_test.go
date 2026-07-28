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

type mockExpenseRepo struct {
	expenses map[string]domain.Expense
	err      error
	getErr   error
}

func (m *mockExpenseRepo) Create(ctx context.Context, exp *domain.Expense) error {
	if m.err != nil {
		return m.err
	}
	if exp.ID == "" {
		exp.ID = "generated-exp-id"
	}
	if m.expenses == nil {
		m.expenses = make(map[string]domain.Expense)
	}
	m.expenses[exp.ID] = *exp
	return nil
}

func (m *mockExpenseRepo) Update(ctx context.Context, exp *domain.Expense) error {
	if m.err != nil {
		return m.err
	}
	if m.expenses == nil {
		return errors.New("expense not found")
	}
	m.expenses[exp.ID] = *exp
	return nil
}

func (m *mockExpenseRepo) Delete(ctx context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.expenses, id)
	return nil
}

func (m *mockExpenseRepo) ListByMonth(ctx context.Context, year int, month int) ([]domain.Expense, error) {
	if m.err != nil {
		return nil, m.err
	}
	var list []domain.Expense
	for _, e := range m.expenses {
		list = append(list, e)
	}
	return list, nil
}

func (m *mockExpenseRepo) Get(ctx context.Context, id string) (*domain.Expense, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	e, ok := m.expenses[id]
	if !ok {
		return nil, errors.New("expense not found")
	}
	return &e, nil
}

func TestExpenseHandler_Create(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]domain.User{
			"user-1": {ID: "user-1", Name: "Rafael", Salary: 5000},
		},
	}

	t.Run("Success", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Supermarket","amount":150.50,"date":"2026-07-27","category":"Mercado","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rec.Code)
		}
	})

	t.Run("Invalid body", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString("invalid json"))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Missing description", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"amount":150.50,"date":"2026-07-27","category":"Mercado","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Amount zero or negative", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":0,"date":"2026-07-27","category":"Mercado","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Missing date", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":10,"category":"Mercado","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid category", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":10,"date":"2026-07-27","category":"InvalidCat","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Missing payer_id", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":10,"date":"2026-07-27","category":"Mercado"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Payer not found", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":10,"date":"2026-07-27","category":"Mercado","payer_id":"non-existent"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		expRepo := &mockExpenseRepo{err: errors.New("db insert error")}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Test","amount":10,"date":"2026-07-27","category":"Mercado","payer_id":"user-1"}`
		req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}

func TestExpenseHandler_Update(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]domain.User{
			"user-1": {ID: "user-1", Name: "Rafael", Salary: 5000},
		},
	}

	t.Run("Success", func(t *testing.T) {
		expRepo := &mockExpenseRepo{
			expenses: map[string]domain.Expense{
				"exp-1": {ID: "exp-1", Description: "Old", Amount: 50, Date: "2026-07-01", Category: "Casa", PayerID: "user-1"},
			},
		}
		h := NewExpenseHandler(expRepo, userRepo)

		body := `{"description":"Updated","amount":100}`
		req := httptest.NewRequest(http.MethodPut, "/api/expenses/exp-1", bytes.NewBufferString(body))
		req.SetPathValue("id", "exp-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Missing path ID", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		req := httptest.NewRequest(http.MethodPut, "/api/expenses/", bytes.NewBufferString(`{}`))
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Expense not found", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, userRepo)

		req := httptest.NewRequest(http.MethodPut, "/api/expenses/exp-99", bytes.NewBufferString(`{}`))
		req.SetPathValue("id", "exp-99")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rec.Code)
		}
	})

	t.Run("Invalid category update", func(t *testing.T) {
		expRepo := &mockExpenseRepo{
			expenses: map[string]domain.Expense{
				"exp-1": {ID: "exp-1", Category: "Casa"},
			},
		}
		h := NewExpenseHandler(expRepo, userRepo)

		req := httptest.NewRequest(http.MethodPut, "/api/expenses/exp-1", bytes.NewBufferString(`{"category":"Invalid"}`))
		req.SetPathValue("id", "exp-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid payer update", func(t *testing.T) {
		expRepo := &mockExpenseRepo{
			expenses: map[string]domain.Expense{
				"exp-1": {ID: "exp-1", Category: "Casa"},
			},
		}
		h := NewExpenseHandler(expRepo, userRepo)

		req := httptest.NewRequest(http.MethodPut, "/api/expenses/exp-1", bytes.NewBufferString(`{"payer_id":"unknown"}`))
		req.SetPathValue("id", "exp-1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})
}

func TestExpenseHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expRepo := &mockExpenseRepo{
			expenses: map[string]domain.Expense{
				"exp-1": {ID: "exp-1"},
			},
		}
		h := NewExpenseHandler(expRepo, nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/expenses/exp-1", nil)
		req.SetPathValue("id", "exp-1")
		rec := httptest.NewRecorder()

		h.Delete(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Missing path ID", func(t *testing.T) {
		expRepo := &mockExpenseRepo{}
		h := NewExpenseHandler(expRepo, nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/expenses/", nil)
		rec := httptest.NewRecorder()

		h.Delete(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		expRepo := &mockExpenseRepo{err: errors.New("delete error")}
		h := NewExpenseHandler(expRepo, nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/expenses/exp-1", nil)
		req.SetPathValue("id", "exp-1")
		rec := httptest.NewRecorder()

		h.Delete(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}
