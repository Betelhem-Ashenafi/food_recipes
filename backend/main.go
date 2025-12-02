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

	// Register categories route
	http.HandleFunc("/categories", handlers.GetCategoriesHandler)

	// Register recipe routes (Protected)
	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
case http.MethodPost:
			handlers.AuthMiddleware(handlers.CreateRecipeHandler)(w, r)
		case http.MethodGet:
			handlers.GetRecipesHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register recipe detail routes (edit/delete)
	http.HandleFunc("/recipes/", func(w http.ResponseWriter, r *http.Request) {
		// Check for image upload/feature endpoints first
		if len(r.URL.Path) > len("/recipes/") {
			subPath := r.URL.Path[len("/recipes/"):]
			// Check if it's an image upload request: /recipes/{id}/images
			if len(subPath) > 0 && subPath[len(subPath)-len("/images"):] == "/images" {
				if r.Method == http.MethodPost {
					handlers.AuthMiddleware(handlers.UploadRecipeImagesHandler)(w, r)
					return
				}
			}
			// Check if it's a feature image request: /recipes/{id}/images/{imgId}/feature
			if len(subPath) > 0 && len(subPath) > len("/feature") && subPath[len(subPath)-len("/feature"):] == "/feature" {
				if r.Method == http.MethodPost {
					handlers.AuthMiddleware(handlers.FeatureRecipeImageHandler)(w, r)
					return
				}
			}
		}

		switch r.Method {
		case http.MethodPut:
			handlers.AuthMiddleware(handlers.EditRecipeHandler)(w, r)
		case http.MethodDelete:
			handlers.AuthMiddleware(handlers.DeleteRecipeHandler)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting Food Recipes Backend on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
