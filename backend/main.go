package main

import (
	"foodrecipes/handlers"
	"log"
	"net/http"
	"os"
	"strings"

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

	// Register user profile routes - MUST be before /recipes/ to avoid conflicts
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Users route hit: %s %s", r.Method, r.URL.Path)
		// Extract subpath after /users/
		path := strings.TrimPrefix(r.URL.Path, "/users/")
		pathParts := strings.Split(path, "/")
		log.Printf("Path parts after /users/: %v", pathParts)

		// pathParts should be: ["{id}", "bookmarks"] or ["{id}", "purchases"]
		if len(pathParts) >= 2 {
			// Check if it's bookmarks: /users/{id}/bookmarks
			if pathParts[1] == "bookmarks" && r.Method == http.MethodGet {
				log.Printf("Calling GetUserBookmarksHandler for user %s", pathParts[0])
				handlers.AuthMiddleware(handlers.GetUserBookmarksHandler)(w, r)
				return
			}

			// Check if it's purchases: /users/{id}/purchases
			if pathParts[1] == "purchases" && r.Method == http.MethodGet {
				log.Printf("Calling GetUserPurchasesHandler for user %s", pathParts[0])
				handlers.AuthMiddleware(handlers.GetUserPurchasesHandler)(w, r)
				return
			}
		}
		log.Printf("No matching route for /users/ - returning 404")
		http.Error(w, "Not found", http.StatusNotFound)
	})

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
		// Check for ingredients/steps endpoints first
		if len(r.URL.Path) > len("/recipes/") {
			subPath := r.URL.Path[len("/recipes/"):]

			// Check if it's ingredients request: /recipes/{id}/ingredients
			if strings.HasSuffix(subPath, "/ingredients") {
				if r.Method == http.MethodGet {
					handlers.GetRecipeIngredientsHandler(w, r)
					return
				}
			}

			// Check if it's steps request: /recipes/{id}/steps
			if strings.HasSuffix(subPath, "/steps") {
				if r.Method == http.MethodGet {
					handlers.GetRecipeStepsHandler(w, r)
					return
				}
			}

			// Check if it's like request: /recipes/{id}/like
			if strings.HasSuffix(subPath, "/like") {
				handlers.AuthMiddleware(handlers.ToggleLikeHandler)(w, r)
				return
			}

			// Check if it's bookmark request: /recipes/{id}/bookmark
			if strings.HasSuffix(subPath, "/bookmark") {
				handlers.AuthMiddleware(handlers.ToggleBookmarkHandler)(w, r)
				return
			}

			// Check if it's comments request: /recipes/{id}/comments
			if strings.HasSuffix(subPath, "/comments") {
				if r.Method == http.MethodGet {
					handlers.GetCommentsHandler(w, r)
					return
				} else if r.Method == http.MethodPost {
					handlers.AuthMiddleware(handlers.PostCommentHandler)(w, r)
					return
				}
			}

			// Check if it's rate request: /recipes/{id}/rate
			if strings.HasSuffix(subPath, "/rate") {
				if r.Method == http.MethodGet {
					handlers.GetRatingHandler(w, r)
					return
				} else if r.Method == http.MethodPost {
					handlers.AuthMiddleware(handlers.RateRecipeHandler)(w, r)
					return
				}
			}

			// Check if it's like check: /recipes/{id}/like/check
			if strings.HasSuffix(subPath, "/like/check") {
				if r.Method == http.MethodGet {
					handlers.AuthMiddleware(handlers.CheckLikeHandler)(w, r)
					return
				}
			}

			// Check if it's bookmark check: /recipes/{id}/bookmark/check
			if strings.HasSuffix(subPath, "/bookmark/check") {
				if r.Method == http.MethodGet {
					handlers.AuthMiddleware(handlers.CheckBookmarkHandler)(w, r)
					return
				}
			}

			// Check if it's purchase check: /recipes/{id}/purchase/check
			if strings.HasSuffix(subPath, "/purchase/check") {
				if r.Method == http.MethodGet {
					handlers.AuthMiddleware(handlers.CheckPurchaseHandler)(w, r)
					return
				}
			}

			// Check for image endpoints: /recipes/{id}/images
			if strings.HasSuffix(subPath, "/images") {
				if r.Method == http.MethodGet {
					// GET /recipes/{id}/images - Public endpoint to fetch all images
					handlers.GetRecipeImagesHandler(w, r)
					return
				} else if r.Method == http.MethodPost {
					// POST /recipes/{id}/images - Upload images (requires auth)
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

		}

		switch r.Method {
		case http.MethodGet:
			// GET /recipes/{id} - Get single recipe (public)
			handlers.GetRecipeByIDHandler(w, r)
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

	// Payment callback (webhook) - no auth required, Chapa calls this
	http.HandleFunc("/payment/callback", handlers.PaymentCallbackHandler)

	// Register Hasura Action routes
	http.HandleFunc("/hasura/login", handlers.HasuraLoginHandler)
	http.HandleFunc("/hasura/signup", handlers.HasuraSignupHandler)
	http.HandleFunc("/hasura/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.HasuraUploadHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register Hasura Event Trigger routes
	http.HandleFunc("/events/new-recipe", handlers.NewRecipeEventHandler)

	log.Println("Starting Food Recipes Backend on :8081...")
	log.Fatal(http.ListenAndServe(":8081", corsMiddleware(http.DefaultServeMux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for development
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
