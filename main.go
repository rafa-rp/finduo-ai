package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"finduo-ai/internal/handler"
	"finduo-ai/internal/repository/sqlite"
)

// loadEnv reads a local .env file and sets environment variables manually
// to avoid installing external dependencies.
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		// Ignore if .env is missing (e.g. in production/docker)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
				(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
			os.Setenv(key, val)
		}
	}
}

// corsMiddleware adds standard CORS headers for development/production web apps.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggerMiddleware prints basic request information.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "finduo.db"
		log.Println("INFO: DATABASE_URL is not set. Defaulting to local SQLite file 'finduo.db'.")
	}

	log.Printf("Starting Finduo AI Backend on port %s...", port)

	var db *sqlite.DB
	var err error

	db, err = sqlite.Connect(dbURL)
	if err != nil {
		log.Printf("WARNING: Could not connect to SQLite database: %v", err)
	} else {
		defer db.Close()
		log.Println("Successfully connected to SQLite database.")

		// Initialize DB tables automatically on startup
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = db.InitSchema(ctx)
		cancel()
		if err != nil {
			log.Fatalf("CRITICAL: Failed to initialize schema: %v", err)
		}
		log.Println("Database schema checked/initialized successfully.")
	}

	// Setup Repositories (will handle null DB gracefully or return error when called)
	var userRepo *sqlite.UserRepository
	var expenseRepo *sqlite.ExpenseRepository
	var settlementRepo *sqlite.SettlementRepository

	if db != nil {
		userRepo = sqlite.NewUserRepository(db)
		expenseRepo = sqlite.NewExpenseRepository(db)
		settlementRepo = sqlite.NewSettlementRepository(db)
	}

	// Router setup using Go 1.22+ new enhanced ServeMux pattern matching
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"UP","time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Register handlers if DB is initialized
	if db != nil {
		userHandler := handler.NewUserHandler(userRepo)
		expenseHandler := handler.NewExpenseHandler(expenseRepo, userRepo)
		summaryHandler := handler.NewSummaryHandler(userRepo, expenseRepo, settlementRepo)

		// Users
		mux.HandleFunc("GET /api/users", userHandler.List)
		mux.HandleFunc("POST /api/users", userHandler.Save)

		// Expenses
		mux.HandleFunc("POST /api/expenses", expenseHandler.Create)
		mux.HandleFunc("PUT /api/expenses/{id}", expenseHandler.Update)
		mux.HandleFunc("DELETE /api/expenses/{id}", expenseHandler.Delete)

		// Summary & Settlement
		mux.HandleFunc("GET /api/summary", summaryHandler.Get)
		mux.HandleFunc("POST /api/summary/settle", summaryHandler.Settle)
	} else {
		// Fallback handlers returning database connection error if DB was not reachable
		dbFallback := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"error":"Database is offline. Please check your SQLite connection."}`)
		}
		mux.HandleFunc("/api/", dbFallback)
	}

	// Wire middlewares
	handlerChain := loggerMiddleware(corsMiddleware(mux))

	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handlerChain); err != nil {
		log.Fatalf("CRITICAL: Server failed to start: %v", err)
	}
}
