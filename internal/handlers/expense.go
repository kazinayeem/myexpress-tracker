package handlers

import (
	"encoding/json"
	"myexpress-tracker/internal/middleware"
	"myexpress-tracker/internal/models"
	"myexpress-tracker/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

// ExpenseHandler handles expense requests
type ExpenseHandler struct {
	expenseRepo  *repository.ExpenseRepository
	categoryRepo *repository.CategoryRepository
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseRepo *repository.ExpenseRepository, categoryRepo *repository.CategoryRepository) *ExpenseHandler {
	return &ExpenseHandler{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateExpense creates a new expense record
func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if expense.CategoryID == 0 || expense.Amount <= 0 || expense.ExpenseDate == "" {
		http.Error(w, `{"error":"category_id, amount (>0), and expense_date are required"}`, http.StatusBadRequest)
		return
	}

	// Verify category exists and is expense type
	category, err := h.categoryRepo.GetByID(expense.CategoryID)
	if err != nil || category == nil || category.Type != "expense" {
		http.Error(w, `{"error":"invalid expense category"}`, http.StatusBadRequest)
		return
	}

	expense.UserID = userID

	if err := h.expenseRepo.Create(&expense); err != nil {
		http.Error(w, `{"error":"failed to create expense"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

// GetExpenses retrieves expense records with optional filters
func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	filters := make(map[string]interface{})
	
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if id, err := strconv.ParseInt(categoryID, 10, 64); err == nil {
			filters["category_id"] = id
		}
	}
	
	if date := r.URL.Query().Get("date"); date != "" {
		filters["date"] = date
	}
	
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	expenses, err := h.expenseRepo.GetByUser(userID, filters)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch expenses"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// UpdateExpense updates an existing expense record
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get expense ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error":"expense id required"}`, http.StatusBadRequest)
		return
	}

	expenseID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid expense id"}`, http.StatusBadRequest)
		return
	}

	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if expense.CategoryID == 0 || expense.Amount <= 0 || expense.ExpenseDate == "" {
		http.Error(w, `{"error":"category_id, amount (>0), and expense_date are required"}`, http.StatusBadRequest)
		return
	}

	// Verify category exists and is expense type
	category, err := h.categoryRepo.GetByID(expense.CategoryID)
	if err != nil || category == nil || category.Type != "expense" {
		http.Error(w, `{"error":"invalid expense category"}`, http.StatusBadRequest)
		return
	}

	expense.ID = expenseID
	expense.UserID = userID

	if err := h.expenseRepo.Update(&expense); err != nil {
		http.Error(w, `{"error":"failed to update expense"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// DeleteExpense deletes an expense record
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get expense ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error":"expense id required"}`, http.StatusBadRequest)
		return
	}

	expenseID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid expense id"}`, http.StatusBadRequest)
		return
	}

	if err := h.expenseRepo.Delete(expenseID, userID); err != nil {
		http.Error(w, `{"error":"failed to delete expense"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "expense deleted successfully"})
}
