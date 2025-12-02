package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"foodrecipes/models"

	"github.com/golang-jwt/jwt"
)

// Middleware to validate JWT and extract User ID
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // In production, use env variable
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user_id to context (simplified for now, we'll just pass it or use a global/context)
		// For this simple handler, we will extract it inside the handler for now to keep it simple
		// or we can set it in a context. Let's just validate here.
		next(w, r)
	}
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract User ID from Token
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	userID := int(token.Claims.(jwt.MapClaims)["user_id"].(float64))

	// 2. Parse Request
	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Start Transaction
	tx, err := DB.Beginx()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 4. Insert Recipe
	var recipeID int
	err = tx.QueryRow(`
		INSERT INTO recipes (user_id, category_id, title, description, preparation_time, price, thumbnail_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, req.CategoryID, req.Title, req.Description, req.PreparationTime, req.Price, req.ThumbnailURL).Scan(&recipeID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to save recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Insert Ingredients
	for _, ing := range req.Ingredients {
		_, err = tx.Exec(`
			INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit)
			VALUES ($1, $2, $3, $4)
		`, recipeID, ing.Name, ing.Quantity, ing.Unit)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to save ingredients", http.StatusInternalServerError)
			return
		}
	}

	// 6. Insert Steps
	for i, step := range req.Steps {
		_, err = tx.Exec(`
			INSERT INTO recipe_steps (recipe_id, step_number, instruction, image_url)
			VALUES ($1, $2, $3, $4)
		`, recipeID, i+1, step.Instruction, step.ImageURL)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to save steps", http.StatusInternalServerError)
			return
		}
	}

	// 7. Commit
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": recipeID, "message": "Recipe created successfully"})
}

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	var recipes []models.Recipe
	err := DB.Select(&recipes, "SELECT * FROM recipes ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
