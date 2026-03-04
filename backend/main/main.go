package main

import (
	"log"
	"net/http"
	"os"

	"foodrecipes/handlers"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Database connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://fooduser:foodpass@127.0.0.1:5433/foodrecipes?sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Pass the database connection to the handlers package
	handlers.SetDB(db)

	// Set up routes for Hasura actions
	http.HandleFunc("/hasura/login", handlers.HasuraLoginHandler)
	http.HandleFunc("/hasura/signup", handlers.HasuraSignupHandler)
	http.HandleFunc("/hasura/upload", handlers.HasuraUploadHandler)

	// Add CORS middleware for the frontend
	handler := corsMiddleware(http.DefaultServeMux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// corsMiddleware adds CORS headers to allow requests from your frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
