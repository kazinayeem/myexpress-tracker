package handlers

import (
	"database/sql"
	"fmt"
	"myexpress-tracker/internal/middleware"
	"myexpress-tracker/internal/models"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ExportHandler handles export requests
type ExportHandler struct {
	db *sql.DB
}

// NewExportHandler creates a new export handler
func NewExportHandler(db *sql.DB) *ExportHandler {
	return &ExportHandler{db: db}
}

// ExportToPDF exports income and expense data to PDF
func (h *ExportHandler) ExportToPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse query parameters for date filtering
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02") // Last month
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02") // Today
	}

	// Fetch data
	incomes, err := h.getIncomesForExport(userID, startDate, endDate)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch income data"}`, http.StatusInternalServerError)
		return
	}

	expenses, err := h.getExpensesForExport(userID, startDate, endDate)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch expense data"}`, http.StatusInternalServerError)
		return
	}

	// Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Income & Expense Report")
	pdf.Ln(8)

	// Date range
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 6, fmt.Sprintf("Period: %s to %s", startDate, endDate))
	pdf.Ln(10)

	// Summary
	totalIncome := 0.0
	for _, inc := range incomes {
		totalIncome += inc.Amount
	}
	totalExpense := 0.0
	for _, exp := range expenses {
		totalExpense += exp.Amount
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Summary")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(95, 6, fmt.Sprintf("Total Income: $%.2f", totalIncome))
	pdf.Cell(95, 6, fmt.Sprintf("Total Expense: $%.2f", totalExpense))
	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(190, 6, fmt.Sprintf("Balance: $%.2f", totalIncome-totalExpense))
	pdf.Ln(10)

	// Income section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Income Details")
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(30, 6, "Date")
	pdf.Cell(40, 6, "Category")
	pdf.Cell(30, 6, "Amount")
	pdf.Cell(90, 6, "Description")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 9)
	for _, inc := range incomes {
		pdf.Cell(30, 6, inc.IncomeDate)
		pdf.Cell(40, 6, inc.CategoryName)
		pdf.Cell(30, 6, fmt.Sprintf("$%.2f", inc.Amount))
		pdf.Cell(90, 6, inc.Description)
		pdf.Ln(6)
	}

	pdf.Ln(5)

	// Expense section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Expense Details")
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(30, 6, "Date")
	pdf.Cell(40, 6, "Category")
	pdf.Cell(30, 6, "Amount")
	pdf.Cell(90, 6, "Description")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 9)
	for _, exp := range expenses {
		pdf.Cell(30, 6, exp.ExpenseDate)
		pdf.Cell(40, 6, exp.CategoryName)
		pdf.Cell(30, 6, fmt.Sprintf("$%.2f", exp.Amount))
		pdf.Cell(90, 6, exp.Description)
		pdf.Ln(6)
	}

	// Output PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s_to_%s.pdf", startDate, endDate))
	
	if err := pdf.Output(w); err != nil {
		http.Error(w, `{"error":"failed to generate PDF"}`, http.StatusInternalServerError)
		return
	}
}

// getIncomesForExport retrieves income data for export
func (h *ExportHandler) getIncomesForExport(userID int64, startDate, endDate string) ([]models.Income, error) {
	query := `
		SELECT i.id, i.user_id, i.category_id, i.amount, i.description, i.income_date, i.created_at, i.updated_at, c.name
		FROM income i
		JOIN categories c ON i.category_id = c.id
		WHERE i.user_id = ? AND i.income_date >= ? AND i.income_date <= ?
		ORDER BY i.income_date DESC
	`

	rows, err := h.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		if err := rows.Scan(
			&income.ID, &income.UserID, &income.CategoryID, &income.Amount, &income.Description,
			&income.IncomeDate, &income.CreatedAt, &income.UpdatedAt, &income.CategoryName,
		); err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

// getExpensesForExport retrieves expense data for export
func (h *ExportHandler) getExpensesForExport(userID int64, startDate, endDate string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.expense_date, e.created_at, e.updated_at, c.name
		FROM expense e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ? AND e.expense_date >= ? AND e.expense_date <= ?
		ORDER BY e.expense_date DESC
	`

	rows, err := h.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Amount, &expense.Description,
			&expense.ExpenseDate, &expense.CreatedAt, &expense.UpdatedAt, &expense.CategoryName,
		); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}
