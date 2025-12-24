package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"myexpress-tracker/configs"
	"myexpress-tracker/internal/auth"
	"myexpress-tracker/internal/database"
	"myexpress-tracker/internal/handlers"
	"myexpress-tracker/internal/middleware"
	"myexpress-tracker/internal/repository"
)

func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	incomeRepo := repository.NewIncomeRepository(db.DB)
	expenseRepo := repository.NewExpenseRepository(db.DB)

	// Initialize auth service
	authService := auth.NewService(cfg.JWTSecret, cfg.JWTExpiration)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, authService)
	userHandler := handlers.NewUserHandler(userRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	incomeHandler := handlers.NewIncomeHandler(incomeRepo, categoryRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo, categoryRepo)
	dashboardHandler := handlers.NewDashboardHandler(db.DB)
	exportHandler := handlers.NewExportHandler(db.DB)

	// Create router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// Protected routes - User Profile & Settings
	userMux := http.NewServeMux()
	userMux.HandleFunc("/api/user/profile", userHandler.GetProfile)
	userMux.HandleFunc("/api/user/settings", userHandler.UpdateSettings)
	mux.Handle("/api/user/profile", middleware.AuthMiddleware(authService)(userMux))
	mux.Handle("/api/user/settings", middleware.AuthMiddleware(authService)(userMux))

	// Protected routes - Categories (can be accessed with auth)
	categoryMux := http.NewServeMux()
	categoryMux.HandleFunc("/api/categories", categoryHandler.GetCategories)
	mux.Handle("/api/categories", middleware.AuthMiddleware(authService)(categoryMux))

	// Protected routes - Income
	incomeMux := http.NewServeMux()
	incomeMux.HandleFunc("/api/income", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			incomeHandler.GetIncomes(w, r)
		} else if r.Method == http.MethodPost {
			incomeHandler.CreateIncome(w, r)
		} else {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	incomeMux.HandleFunc("/api/income/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			incomeHandler.UpdateIncome(w, r)
		} else if r.Method == http.MethodDelete {
			incomeHandler.DeleteIncome(w, r)
		} else {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/api/income", middleware.AuthMiddleware(authService)(incomeMux))
	mux.Handle("/api/income/", middleware.AuthMiddleware(authService)(incomeMux))

	// Protected routes - Expense
	expenseMux := http.NewServeMux()
	expenseMux.HandleFunc("/api/expense", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			expenseHandler.GetExpenses(w, r)
		} else if r.Method == http.MethodPost {
			expenseHandler.CreateExpense(w, r)
		} else {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	expenseMux.HandleFunc("/api/expense/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			expenseHandler.UpdateExpense(w, r)
		} else if r.Method == http.MethodDelete {
			expenseHandler.DeleteExpense(w, r)
		} else {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/api/expense", middleware.AuthMiddleware(authService)(expenseMux))
	mux.Handle("/api/expense/", middleware.AuthMiddleware(authService)(expenseMux))

	// Protected routes - Dashboard
	dashboardMux := http.NewServeMux()
	dashboardMux.HandleFunc("/api/dashboard", dashboardHandler.GetDashboard)
	mux.Handle("/api/dashboard", middleware.AuthMiddleware(authService)(dashboardMux))

	// Protected routes - Export
	exportMux := http.NewServeMux()
	exportMux.HandleFunc("/api/export/pdf", exportHandler.ExportToPDF)
	mux.Handle("/api/export/pdf", middleware.AuthMiddleware(authService)(exportMux))

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	// Apply CORS middleware
	handler := middleware.CORSMiddleware(mux)

	// Get current working directory
	cwd, _ := os.Getwd()
	log.Printf("Current directory: %s", cwd)
	log.Printf("Database path: %s", cfg.DatabasePath)
	log.Printf("Database full path: %s", filepath.Join(cwd, cfg.DatabasePath))

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Environment: %s", cfg.Environment)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
