package models

import "time"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Currency     string    `json:"currency"` // User's preferred currency (USD, EUR, GBP, etc.)
	Theme        string    `json:"theme"` // User's theme preference (light/dark)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Category represents an income or expense category
type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "income" or "expense"
	CreatedAt time.Time `json:"created_at"`
}

// Income represents an income record
type Income struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CategoryID  int64     `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	IncomeDate  string    `json:"income_date"` // Date in YYYY-MM-DD format
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Joined fields
	CategoryName string `json:"category_name,omitempty"`
}

// Expense represents an expense record
type Expense struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CategoryID  int64     `json:"category_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	ExpenseDate string    `json:"expense_date"` // Date in YYYY-MM-DD format
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Joined fields
	CategoryName string `json:"category_name,omitempty"`
}

// DashboardSummary represents dashboard statistics
type DashboardSummary struct {
	TotalIncome      float64            `json:"total_income"`
	TotalExpense     float64            `json:"total_expense"`
	Balance          float64            `json:"balance"`
	TodayIncome      float64            `json:"today_income"`
	TodayExpense     float64            `json:"today_expense"`
	MonthlyIncome    float64            `json:"monthly_income"`
	MonthlyExpense   float64            `json:"monthly_expense"`
	DailyData        []DailyData        `json:"daily_data"`
	CategoryBreakdown CategoryBreakdown `json:"category_breakdown"`
}

// DailyData represents income and expense for a specific day
type DailyData struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// CategoryBreakdown represents spending by category
type CategoryBreakdown struct {
	IncomeByCategory  map[string]float64 `json:"income_by_category"`
	ExpenseByCategory map[string]float64 `json:"expense_by_category"`
}
