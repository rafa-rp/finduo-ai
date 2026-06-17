package handler

import (
	"encoding/json"
	"net/http"

	"finduo-ai/internal/domain"
)

type UserHandler struct {
	repo domain.UserRepository
}

// NewUserHandler initializes a new UserHandler.
func NewUserHandler(repo domain.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// List handles GET /api/users to list all participants.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []domain.User{}
	}
	writeJSON(w, http.StatusOK, users)
}

// Save handles POST /api/users to create or update a participant.
func (h *UserHandler) Save(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Salary float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if input.Salary < 0 {
		writeError(w, http.StatusBadRequest, "salary cannot be negative")
		return
	}

	user := domain.User{
		ID:     input.ID,
		Name:   input.Name,
		Salary: input.Salary,
	}

	if err := h.repo.Save(r.Context(), &user); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}
