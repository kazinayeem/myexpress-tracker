package database

import (
	"fmt"
)

// RunMigrations creates all necessary tables
func (db *DB) RunMigrations() error {
	migrations := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			currency TEXT DEFAULT 'USD',
			theme TEXT DEFAULT 'light',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Create index on email and username for faster lookups
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,

		// Categories table
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Income table
		`CREATE TABLE IF NOT EXISTS income (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			amount REAL NOT NULL CHECK(amount > 0),
			description TEXT,
			income_date DATE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
		)`,

		// Create indexes for income table
		`CREATE INDEX IF NOT EXISTS idx_income_user_id ON income(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_income_date ON income(income_date)`,
		`CREATE INDEX IF NOT EXISTS idx_income_category ON income(category_id)`,

		// Expense table
		`CREATE TABLE IF NOT EXISTS expense (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			amount REAL NOT NULL CHECK(amount > 0),
			description TEXT,
			expense_date DATE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
		)`,

		// Create indexes for expense table
		`CREATE INDEX IF NOT EXISTS idx_expense_user_id ON expense(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_date ON expense(expense_date)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_category ON expense(category_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Insert default categories
	if err := db.insertDefaultCategories(); err != nil {
		return fmt.Errorf("failed to insert default categories: %w", err)
	}

	return nil
}

// insertDefaultCategories adds default income and expense categories
func (db *DB) insertDefaultCategories() error {
	categories := []struct {
		name string
		typ  string
	}{
		// Income categories
		{"Salary", "income"},
		{"Freelance", "income"},
		{"Investment", "income"},
		{"Other Income", "income"},
		
		// Expense categories
		{"Food", "expense"},
		{"Transport", "expense"},
		{"Rent", "expense"},
		{"Utilities", "expense"},
		{"Entertainment", "expense"},
		{"Healthcare", "expense"},
		{"Shopping", "expense"},
		{"Other Expense", "expense"},
	}

	for _, cat := range categories {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO categories (name, type) VALUES (?, ?)`,
			cat.name, cat.typ,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
