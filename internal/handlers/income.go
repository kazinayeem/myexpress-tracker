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

// IncomeHandler handles income requests
type IncomeHandler struct {
	incomeRepo   *repository.IncomeRepository
	categoryRepo *repository.CategoryRepository
}

// NewIncomeHandler creates a new income handler
func NewIncomeHandler(incomeRepo *repository.IncomeRepository, categoryRepo *repository.CategoryRepository) *IncomeHandler {
	return &IncomeHandler{
		incomeRepo:   incomeRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateIncome creates a new income record
func (h *IncomeHandler) CreateIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if income.CategoryID == 0 || income.Amount <= 0 || income.IncomeDate == "" {
		http.Error(w, `{"error":"category_id, amount (>0), and income_date are required"}`, http.StatusBadRequest)
		return
	}

	// Verify category exists and is income type
	category, err := h.categoryRepo.GetByID(income.CategoryID)
	if err != nil || category == nil || category.Type != "income" {
		http.Error(w, `{"error":"invalid income category"}`, http.StatusBadRequest)
		return
	}

	income.UserID = userID

	if err := h.incomeRepo.Create(&income); err != nil {
		http.Error(w, `{"error":"failed to create income"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(income)
}

// GetIncomes retrieves income records with optional filters
func (h *IncomeHandler) GetIncomes(w http.ResponseWriter, r *http.Request) {
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

	incomes, err := h.incomeRepo.GetByUser(userID, filters)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch incomes"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incomes)
}

// UpdateIncome updates an existing income record
func (h *IncomeHandler) UpdateIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get income ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error":"income id required"}`, http.StatusBadRequest)
		return
	}

	incomeID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid income id"}`, http.StatusBadRequest)
		return
	}

	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if income.CategoryID == 0 || income.Amount <= 0 || income.IncomeDate == "" {
		http.Error(w, `{"error":"category_id, amount (>0), and income_date are required"}`, http.StatusBadRequest)
		return
	}

	// Verify category exists and is income type
	category, err := h.categoryRepo.GetByID(income.CategoryID)
	if err != nil || category == nil || category.Type != "income" {
		http.Error(w, `{"error":"invalid income category"}`, http.StatusBadRequest)
		return
	}

	income.ID = incomeID
	income.UserID = userID

	if err := h.incomeRepo.Update(&income); err != nil {
		http.Error(w, `{"error":"failed to update income"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(income)
}

// DeleteIncome deletes an income record
func (h *IncomeHandler) DeleteIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get income ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error":"income id required"}`, http.StatusBadRequest)
		return
	}

	incomeID, err := strconv.ParseInt(pathParts[len(pathParts)-1], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid income id"}`, http.StatusBadRequest)
		return
	}

	if err := h.incomeRepo.Delete(incomeID, userID); err != nil {
		http.Error(w, `{"error":"failed to delete income"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "income deleted successfully"})
}
