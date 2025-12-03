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

	// Register file upload route
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.UploadFileHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Serve static files from uploads directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Register recipe routes (Protected)
	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AuthMiddleware(handlers.CreateRecipeHandler)(w, r)
		case http.MethodGet:
			log.Println("Calling GetRecipesHandler from main.go")
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
			if len(subPath) >= len("/images") && subPath[len(subPath)-len("/images"):] == "/images" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.UploadRecipeImagesHandler)(w, r)
					return
				}
			}
			// Check if it's a feature image request: /recipes/{id}/images/{imgId}/feature
			if len(subPath) >= len("/feature") && subPath[len(subPath)-len("/feature"):] == "/feature" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.FeatureRecipeImageHandler)(w, r)
					return
				}
			}

			// Check for social endpoints
			// /recipes/{id}/like
			if len(subPath) >= len("/like") && subPath[len(subPath)-len("/like"):] == "/like" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.LikeRecipeHandler)(w, r)
					return
				case http.MethodDelete:
					handlers.AuthMiddleware(handlers.UnlikeRecipeHandler)(w, r)
					return
				}
			}

			// /recipes/{id}/bookmark
			if len(subPath) >= len("/bookmark") && subPath[len(subPath)-len("/bookmark"):] == "/bookmark" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.BookmarkRecipeHandler)(w, r)
					return
				case http.MethodDelete:
					handlers.AuthMiddleware(handlers.UnbookmarkRecipeHandler)(w, r)
					return
				}
			}

			// /recipes/{id}/comments
			if len(subPath) >= len("/comments") && subPath[len(subPath)-len("/comments"):] == "/comments" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.CommentRecipeHandler)(w, r)
					return
				case http.MethodGet:
					handlers.GetCommentsHandler(w, r)
					return
				}
			}

			// /recipes/{id}/rate
			if len(subPath) >= len("/rate") && subPath[len(subPath)-len("/rate"):] == "/rate" {
				switch r.Method {
				case http.MethodPost:
					handlers.AuthMiddleware(handlers.RateRecipeHandler)(w, r)
					return
				case http.MethodGet:
					handlers.GetRecipeRatingHandler(w, r)
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

	// Register payment routes
	http.HandleFunc("/payment/initialize", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AuthMiddleware(handlers.InitializePaymentHandler)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/payment/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.AuthMiddleware(handlers.VerifyPaymentHandler)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register Hasura Action routes
	http.HandleFunc("/hasura/login", handlers.HasuraLoginHandler)

	// Register Hasura Event Trigger routes
	http.HandleFunc("/events/new-recipe", handlers.NewRecipeEventHandler)

	log.Println("Starting Food Recipes Backend on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
