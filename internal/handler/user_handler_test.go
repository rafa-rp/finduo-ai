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

type mockUserRepo struct {
	users  map[string]domain.User
	err    error
	getErr error
}

func (m *mockUserRepo) List(ctx context.Context) ([]domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []domain.User
	for _, u := range m.users {
		res = append(res, u)
	}
	return res, nil
}

func (m *mockUserRepo) Save(ctx context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	if user.ID == "" {
		user.ID = "generated-id"
	}
	if m.users == nil {
		m.users = make(map[string]domain.User)
	}
	m.users[user.ID] = *user
	return nil
}

func (m *mockUserRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func TestUserHandler_List(t *testing.T) {
	t.Run("Success with users", func(t *testing.T) {
		repo := &mockUserRepo{
			users: map[string]domain.User{
				"1": {ID: "1", Name: "Rafael", Salary: 5000},
			},
		}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Success empty users", func(t *testing.T) {
		repo := &mockUserRepo{}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		repo := &mockUserRepo{err: errors.New("db error")}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()

		h.List(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}

func TestUserHandler_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &mockUserRepo{}
		h := NewUserHandler(repo)

		body := bytes.NewBufferString(`{"name":"Rafael","salary":5000}`)
		req := httptest.NewRequest(http.MethodPost, "/api/users", body)
		rec := httptest.NewRecorder()

		h.Save(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Invalid body", func(t *testing.T) {
		repo := &mockUserRepo{}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString("invalid json"))
		rec := httptest.NewRecorder()

		h.Save(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Missing name", func(t *testing.T) {
		repo := &mockUserRepo{}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{"salary":5000}`))
		rec := httptest.NewRecorder()

		h.Save(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Negative salary", func(t *testing.T) {
		repo := &mockUserRepo{}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{"name":"Rafael","salary":-10}`))
		rec := httptest.NewRecorder()

		h.Save(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Repo error", func(t *testing.T) {
		repo := &mockUserRepo{err: errors.New("db error")}
		h := NewUserHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString(`{"name":"Rafael","salary":5000}`))
		rec := httptest.NewRecorder()

		h.Save(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
	})
}
