package handlers

import (
	"database/sql"
	"encoding/json"
	"myexpress-tracker/internal/middleware"
	"myexpress-tracker/internal/models"
	"net/http"
	"time"
)

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	db *sql.DB
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetDashboard retrieves dashboard summary data
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	summary := models.DashboardSummary{
		CategoryBreakdown: models.CategoryBreakdown{
			IncomeByCategory:  make(map[string]float64),
			ExpenseByCategory: make(map[string]float64),
		},
	}

	// Get today's date
	today := time.Now().Format("2006-01-02")
	currentMonth := time.Now().Format("2006-01")

	// Total income
	err := h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM income WHERE user_id = ?`, userID).Scan(&summary.TotalIncome)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch total income"}`, http.StatusInternalServerError)
		return
	}

	// Total expense
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expense WHERE user_id = ?`, userID).Scan(&summary.TotalExpense)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch total expense"}`, http.StatusInternalServerError)
		return
	}

	// Today's income
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM income WHERE user_id = ? AND income_date = ?`, userID, today).Scan(&summary.TodayIncome)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch today income"}`, http.StatusInternalServerError)
		return
	}

	// Today's expense
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expense WHERE user_id = ? AND expense_date = ?`, userID, today).Scan(&summary.TodayExpense)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch today expense"}`, http.StatusInternalServerError)
		return
	}

	// Monthly income
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM income WHERE user_id = ? AND strftime('%Y-%m', income_date) = ?`, userID, currentMonth).Scan(&summary.MonthlyIncome)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch monthly income"}`, http.StatusInternalServerError)
		return
	}

	// Monthly expense
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expense WHERE user_id = ? AND strftime('%Y-%m', expense_date) = ?`, userID, currentMonth).Scan(&summary.MonthlyExpense)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch monthly expense"}`, http.StatusInternalServerError)
		return
	}

	// Calculate balance
	summary.Balance = summary.TotalIncome - summary.TotalExpense

	// Get last 30 days data for charts
	dailyData, err := h.getDailyData(userID, 30)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch daily data"}`, http.StatusInternalServerError)
		return
	}
	summary.DailyData = dailyData

	// Get category breakdown
	incomeByCategory, err := h.getCategoryBreakdown(userID, "income")
	if err != nil {
		http.Error(w, `{"error":"failed to fetch income breakdown"}`, http.StatusInternalServerError)
		return
	}
	summary.CategoryBreakdown.IncomeByCategory = incomeByCategory

	expenseByCategory, err := h.getCategoryBreakdown(userID, "expense")
	if err != nil {
		http.Error(w, `{"error":"failed to fetch expense breakdown"}`, http.StatusInternalServerError)
		return
	}
	summary.CategoryBreakdown.ExpenseByCategory = expenseByCategory

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// getDailyData retrieves daily income and expense for the last N days
func (h *DashboardHandler) getDailyData(userID int64, days int) ([]models.DailyData, error) {
	query := `
		WITH RECURSIVE dates(date) AS (
			SELECT date('now', '-' || ? || ' days')
			UNION ALL
			SELECT date(date, '+1 day')
			FROM dates
			WHERE date < date('now')
		)
		SELECT 
			d.date,
			COALESCE(SUM(i.amount), 0) as income,
			COALESCE(SUM(e.amount), 0) as expense
		FROM dates d
		LEFT JOIN income i ON i.user_id = ? AND i.income_date = d.date
		LEFT JOIN expense e ON e.user_id = ? AND e.expense_date = d.date
		GROUP BY d.date
		ORDER BY d.date
	`

	rows, err := h.db.Query(query, days-1, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dailyData []models.DailyData
	for rows.Next() {
		var data models.DailyData
		if err := rows.Scan(&data.Date, &data.Income, &data.Expense); err != nil {
			return nil, err
		}
		dailyData = append(dailyData, data)
	}

	return dailyData, nil
}

// getCategoryBreakdown retrieves spending/income breakdown by category
func (h *DashboardHandler) getCategoryBreakdown(userID int64, categoryType string) (map[string]float64, error) {
	var query string
	if categoryType == "income" {
		query = `
			SELECT c.name, COALESCE(SUM(i.amount), 0) as total
			FROM categories c
			LEFT JOIN income i ON i.category_id = c.id AND i.user_id = ?
			WHERE c.type = 'income'
			GROUP BY c.id, c.name
			HAVING total > 0
			ORDER BY total DESC
		`
	} else {
		query = `
			SELECT c.name, COALESCE(SUM(e.amount), 0) as total
			FROM categories c
			LEFT JOIN expense e ON e.category_id = c.id AND e.user_id = ?
			WHERE c.type = 'expense'
			GROUP BY c.id, c.name
			HAVING total > 0
			ORDER BY total DESC
		`
	}

	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	breakdown := make(map[string]float64)
	for rows.Next() {
		var name string
		var total float64
		if err := rows.Scan(&name, &total); err != nil {
			return nil, err
		}
		breakdown[name] = total
	}

	return breakdown, nil
}
