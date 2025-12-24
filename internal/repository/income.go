package repository

import (
	"database/sql"
	"fmt"
	"myexpress-tracker/internal/models"
	"strings"
)

// IncomeRepository handles database operations for income
type IncomeRepository struct {
	db *sql.DB
}

// NewIncomeRepository creates a new income repository
func NewIncomeRepository(db *sql.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

// Create creates a new income record
func (r *IncomeRepository) Create(income *models.Income) error {
	query := `
		INSERT INTO income (user_id, category_id, amount, description, income_date)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, income.UserID, income.CategoryID, income.Amount, income.Description, income.IncomeDate)
	if err != nil {
		return fmt.Errorf("failed to create income: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	income.ID = id
	return nil
}

// Update updates an existing income record
func (r *IncomeRepository) Update(income *models.Income) error {
	query := `
		UPDATE income
		SET category_id = ?, amount = ?, description = ?, income_date = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.Exec(query, income.CategoryID, income.Amount, income.Description, income.IncomeDate, income.ID, income.UserID)
	if err != nil {
		return fmt.Errorf("failed to update income: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("income not found or unauthorized")
	}

	return nil
}

// Delete deletes an income record
func (r *IncomeRepository) Delete(id, userID int64) error {
	query := `DELETE FROM income WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete income: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("income not found or unauthorized")
	}

	return nil
}

// GetByID retrieves an income record by ID
func (r *IncomeRepository) GetByID(id, userID int64) (*models.Income, error) {
	query := `
		SELECT i.id, i.user_id, i.category_id, i.amount, i.description, i.income_date, i.created_at, i.updated_at, c.name
		FROM income i
		JOIN categories c ON i.category_id = c.id
		WHERE i.id = ? AND i.user_id = ?
	`
	
	income := &models.Income{}
	err := r.db.QueryRow(query, id, userID).Scan(
		&income.ID, &income.UserID, &income.CategoryID, &income.Amount, &income.Description,
		&income.IncomeDate, &income.CreatedAt, &income.UpdatedAt, &income.CategoryName,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get income by id: %w", err)
	}

	return income, nil
}

// GetByUser retrieves all income records for a user with optional filters
func (r *IncomeRepository) GetByUser(userID int64, filters map[string]interface{}) ([]models.Income, error) {
	query := `
		SELECT i.id, i.user_id, i.category_id, i.amount, i.description, i.income_date, i.created_at, i.updated_at, c.name
		FROM income i
		JOIN categories c ON i.category_id = c.id
		WHERE i.user_id = ?
	`
	
	args := []interface{}{userID}
	
	// Add filters
	if categoryID, ok := filters["category_id"].(int64); ok && categoryID > 0 {
		query += " AND i.category_id = ?"
		args = append(args, categoryID)
	}
	
	if date, ok := filters["date"].(string); ok && date != "" {
		query += " AND i.income_date = ?"
		args = append(args, date)
	}
	
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query += " AND i.income_date >= ?"
		args = append(args, startDate)
	}
	
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query += " AND i.income_date <= ?"
		args = append(args, endDate)
	}
	
	query += " ORDER BY i.income_date DESC, i.created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query income: %w", err)
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		if err := rows.Scan(
			&income.ID, &income.UserID, &income.CategoryID, &income.Amount, &income.Description,
			&income.IncomeDate, &income.CreatedAt, &income.UpdatedAt, &income.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan income: %w", err)
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

// GetTotalByUser calculates total income for a user with optional filters
func (r *IncomeRepository) GetTotalByUser(userID int64, filters map[string]interface{}) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM income WHERE user_id = ?`
	args := []interface{}{userID}
	
	// Add filters
	conditions := []string{}
	if categoryID, ok := filters["category_id"].(int64); ok && categoryID > 0 {
		conditions = append(conditions, "category_id = ?")
		args = append(args, categoryID)
	}
	
	if date, ok := filters["date"].(string); ok && date != "" {
		conditions = append(conditions, "income_date = ?")
		args = append(args, date)
	}
	
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		conditions = append(conditions, "income_date >= ?")
		args = append(args, startDate)
	}
	
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		conditions = append(conditions, "income_date <= ?")
		args = append(args, endDate)
	}
	
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	
	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total income: %w", err)
	}

	return total, nil
}
