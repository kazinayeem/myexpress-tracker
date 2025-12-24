package handlers

import (
	"encoding/json"
	"myexpress-tracker/internal/auth"
	"myexpress-tracker/internal/models"
	"myexpress-tracker/internal/repository"
	"net/http"
	"strings"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo    *repository.UserRepository
	authService *auth.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repository.UserRepository, authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		authService: authService,
	}
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token    string       `json:"token"`
	User     *models.User `json:"user"`
	Message  string       `json:"message"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"email, username, and password are required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"error":"password must be at least 6 characters"}`, http.StatusBadRequest)
		return
	}

	// Check if email exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, `{"error":"email already exists"}`, http.StatusConflict)
		return
	}

	// Check if username exists
	existingUser, err = h.userRepo.GetByUsername(req.Username)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, `{"error":"username already exists"}`, http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Email:        strings.ToLower(req.Email),
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Currency:     "USD",
		Theme:        "light",
	}

	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Send response
	response := AuthResponse{
		Token:   token,
		User:    user,
		Message: "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.EmailOrUsername == "" || req.Password == "" {
		http.Error(w, `{"error":"email/username and password are required"}`, http.StatusBadRequest)
		return
	}

	// Try to find user by email or username
	var user *models.User
	var err error

	// Check if it's an email (contains @)
	if strings.Contains(req.EmailOrUsername, "@") {
		user, err = h.userRepo.GetByEmail(strings.ToLower(req.EmailOrUsername))
	} else {
		user, err = h.userRepo.GetByUsername(req.EmailOrUsername)
	}

	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Send response
	response := AuthResponse{
		Token:   token,
		User:    user,
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
