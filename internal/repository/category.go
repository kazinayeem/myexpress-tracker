package repository

import (
	"database/sql"
	"fmt"
	"myexpress-tracker/internal/models"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `
		SELECT id, name, type, created_at
		FROM categories
		ORDER BY type, name
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Type, &cat.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// GetByType retrieves categories by type (income or expense)
func (r *CategoryRepository) GetByType(categoryType string) ([]models.Category, error) {
	query := `
		SELECT id, name, type, created_at
		FROM categories
		WHERE type = ?
		ORDER BY name
	`
	
	rows, err := r.db.Query(query, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories by type: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Type, &cat.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int64) (*models.Category, error) {
	query := `
		SELECT id, name, type, created_at
		FROM categories
		WHERE id = ?
	`
	
	cat := &models.Category{}
	err := r.db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Type, &cat.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return cat, nil
}
