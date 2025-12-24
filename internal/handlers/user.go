package handlers

import (
	"encoding/json"
	"myexpress-tracker/internal/middleware"
	"myexpress-tracker/internal/repository"
	"net/http"
)

// UserHandler handles user settings requests
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// UpdateSettingsRequest represents settings update request
type UpdateSettingsRequest struct {
	Currency string `json:"currency"`
	Theme    string `json:"theme"`
}

// GetProfile retrieves user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil || user == nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateSettings updates user settings (currency, theme)
func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Update currency if provided
	if req.Currency != "" {
		if err := h.userRepo.UpdateCurrency(userID, req.Currency); err != nil {
			http.Error(w, `{"error":"failed to update currency"}`, http.StatusInternalServerError)
			return
		}
	}

	// Update theme if provided
	if req.Theme != "" {
		if err := h.userRepo.UpdateTheme(userID, req.Theme); err != nil {
			http.Error(w, `{"error":"failed to update theme"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
}
