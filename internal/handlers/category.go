package handlers

import (
	"encoding/json"
	"myexpress-tracker/internal/repository"
	"net/http"
)

// CategoryHandler handles category requests
type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryRepo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// GetCategories retrieves all categories or filtered by type
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	categoryType := r.URL.Query().Get("type")
	
	var categories interface{}
	var err error
	
	if categoryType != "" && (categoryType == "income" || categoryType == "expense") {
		categories, err = h.categoryRepo.GetByType(categoryType)
	} else {
		categories, err = h.categoryRepo.GetAll()
	}

	if err != nil {
		http.Error(w, `{"error":"failed to fetch categories"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
