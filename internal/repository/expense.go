package repository

import (
	"database/sql"
	"fmt"
	"myexpress-tracker/internal/models"
	"strings"
)

// ExpenseRepository handles database operations for expenses
type ExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create creates a new expense record
func (r *ExpenseRepository) Create(expense *models.Expense) error {
	query := `
		INSERT INTO expense (user_id, category_id, amount, description, expense_date)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, expense.UserID, expense.CategoryID, expense.Amount, expense.Description, expense.ExpenseDate)
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	expense.ID = id
	return nil
}

// Update updates an existing expense record
func (r *ExpenseRepository) Update(expense *models.Expense) error {
	query := `
		UPDATE expense
		SET category_id = ?, amount = ?, description = ?, expense_date = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.Exec(query, expense.CategoryID, expense.Amount, expense.Description, expense.ExpenseDate, expense.ID, expense.UserID)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("expense not found or unauthorized")
	}

	return nil
}

// Delete deletes an expense record
func (r *ExpenseRepository) Delete(id, userID int64) error {
	query := `DELETE FROM expense WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("expense not found or unauthorized")
	}

	return nil
}

// GetByID retrieves an expense record by ID
func (r *ExpenseRepository) GetByID(id, userID int64) (*models.Expense, error) {
	query := `
		SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.expense_date, e.created_at, e.updated_at, c.name
		FROM expense e
		JOIN categories c ON e.category_id = c.id
		WHERE e.id = ? AND e.user_id = ?
	`
	
	expense := &models.Expense{}
	err := r.db.QueryRow(query, id, userID).Scan(
		&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
		&expense.ExpenseDate, &expense.CreatedAt, &expense.UpdatedAt, &expense.CategoryName,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get expense by id: %w", err)
	}

	return expense, nil
}

// GetByUser retrieves all expense records for a user with optional filters
func (r *ExpenseRepository) GetByUser(userID int64, filters map[string]interface{}) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.expense_date, e.created_at, e.updated_at, c.name
		FROM expense e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?
	`
	
	args := []interface{}{userID}
	
	// Add filters
	if categoryID, ok := filters["category_id"].(int64); ok && categoryID > 0 {
		query += " AND e.category_id = ?"
		args = append(args, categoryID)
	}
	
	if date, ok := filters["date"].(string); ok && date != "" {
		query += " AND e.expense_date = ?"
		args = append(args, date)
	}
	
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query += " AND e.expense_date >= ?"
		args = append(args, startDate)
	}
	
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query += " AND e.expense_date <= ?"
		args = append(args, endDate)
	}
	
	query += " ORDER BY e.expense_date DESC, e.created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query expense: %w", err)
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.ExpenseDate, &expense.CreatedAt, &expense.UpdatedAt, &expense.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// GetTotalByUser calculates total expense for a user with optional filters
func (r *ExpenseRepository) GetTotalByUser(userID int64, filters map[string]interface{}) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM expense WHERE user_id = ?`
	args := []interface{}{userID}
	
	// Add filters
	conditions := []string{}
	if categoryID, ok := filters["category_id"].(int64); ok && categoryID > 0 {
		conditions = append(conditions, "category_id = ?")
		args = append(args, categoryID)
	}
	
	if date, ok := filters["date"].(string); ok && date != "" {
		conditions = append(conditions, "expense_date = ?")
		args = append(args, date)
	}
	
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		conditions = append(conditions, "expense_date >= ?")
		args = append(args, startDate)
	}
	
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		conditions = append(conditions, "expense_date <= ?")
		args = append(args, endDate)
	}
	
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	
	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total expense: %w", err)
	}

	return total, nil
}
