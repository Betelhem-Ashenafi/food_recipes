package main

import (
	"foodrecipes/handlers"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	connStr := "postgres://fooduser:foodpass@127.0.0.1:5433/foodrecipes?sslmode=disable"
	if envConn := os.Getenv("DATABASE_URL"); envConn != "" {
		connStr = envConn
	}

	// Connect to Postgres
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	// Assign db to handlers.DB
	handlers.DB = db

	// Register login route
	http.HandleFunc("/login", handlers.LoginHandler)
	// Register signup route
	http.HandleFunc("/signup", handlers.SignupHandler)

	// Register recipe routes (Protected)
	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AuthMiddleware(handlers.CreateRecipeHandler)(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetRecipesHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting Food Recipes Backend on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
