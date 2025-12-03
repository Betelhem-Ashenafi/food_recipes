package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"foodrecipes/models"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your-secret-key"), nil
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

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	}
}

func CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract User ID from Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your-secret-key"), nil
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

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)

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
	log.Printf("DEBUG: URL=%s", r.URL.String())
	// Get query params
	q := r.URL.Query()
	timeParam := q.Get("time")
	ingredientParam := q.Get("ingredient")
	titleParam := q.Get("title")
	creatorParam := q.Get("creator")

	log.Printf("DEBUG: Filters: time=%s, ingredient=%s, title=%s, creator=%s", timeParam, ingredientParam, titleParam, creatorParam)

	// Build base query and args
	query := "SELECT r.* FROM recipes r"
	args := []interface{}{}
	where := []string{}

	// Join with users table if filtering by creator
	if creatorParam != "" {
		query += " JOIN users u ON r.user_id = u.id"
		where = append(where, "u.name ILIKE $1")
		args = append(args, "%"+creatorParam+"%")
	}

	// Join with ingredients table if filtering by ingredient
	if ingredientParam != "" {
		query += " JOIN recipe_ingredients ri ON r.id = ri.recipe_id"
		// Adjust placeholder index based on current args length
		where = append(where, fmt.Sprintf("ri.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+ingredientParam+"%")
	}

	// Filter by time
	if timeParam != "" {
		where = append(where, fmt.Sprintf("r.preparation_time <= $%d", len(args)+1))
		args = append(args, timeParam)
	}

	// Filter by title
	if titleParam != "" {
		where = append(where, fmt.Sprintf("r.title ILIKE $%d", len(args)+1))
		args = append(args, "%"+titleParam+"%")
	}

	// Add WHERE clause if needed
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY r.created_at DESC"

	log.Printf("Query: %s", query)
	log.Printf("Args: %v", args)

	recipes := []models.Recipe{}
	err := DB.Select(&recipes, query, args...)
	if err != nil {
		log.Printf("Error getting recipes: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Fetch images for each recipe
	for i := range recipes {
		var images []models.RecipeImage
		err := DB.Select(&images, "SELECT id, recipe_id, url, is_featured FROM recipe_images WHERE recipe_id=$1", recipes[i].ID)
		if err == nil {
			recipes[i].Images = images
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	err := DB.Select(&categories, "SELECT id, name, COALESCE(image_url, '') as image_url FROM categories ORDER BY name ASC")
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// EditRecipeHandler allows the owner to update a recipe
func EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/")
	if idStr == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}
	// Convert to int
	var recipeID int
	_, err := fmt.Sscanf(idStr, "%d", &recipeID)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your-secret-key"), nil
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

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)

	// Check recipe ownership
	var ownerID int
	err = DB.Get(&ownerID, "SELECT user_id FROM recipes WHERE id=$1", recipeID)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden: not recipe owner", http.StatusForbidden)
		return
	}

	// Parse request body
	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := DB.Beginx()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Update recipe
	_, err = tx.Exec(`UPDATE recipes SET category_id=$1, title=$2, description=$3, preparation_time=$4, price=$5, thumbnail_url=$6 WHERE id=$7`,
		req.CategoryID, req.Title, req.Description, req.PreparationTime, req.Price, req.ThumbnailURL, recipeID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}

	// Delete old ingredients/steps
	_, _ = tx.Exec(`DELETE FROM recipe_ingredients WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM recipe_steps WHERE recipe_id=$1`, recipeID)

	// Insert new ingredients
	for _, ing := range req.Ingredients {
		_, err = tx.Exec(`INSERT INTO recipe_ingredients (recipe_id, name, quantity, unit) VALUES ($1, $2, $3, $4)`, recipeID, ing.Name, ing.Quantity, ing.Unit)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update ingredients", http.StatusInternalServerError)
			return
		}
	}

	// Insert new steps
	for i, step := range req.Steps {
		_, err = tx.Exec(`INSERT INTO recipe_steps (recipe_id, step_number, instruction, image_url) VALUES ($1, $2, $3, $4)`, recipeID, i+1, step.Instruction, step.ImageURL)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update steps", http.StatusInternalServerError)
			return
		}
	}

	// Commit
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": recipeID, "message": "Recipe updated successfully"})
}

// DeleteRecipeHandler allows the owner to delete a recipe
func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/")
	if idStr == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}
	var recipeID int
	_, err := fmt.Sscanf(idStr, "%d", &recipeID)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your-secret-key"), nil
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

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)

	// Check recipe ownership
	var ownerID int
	err = DB.Get(&ownerID, "SELECT user_id FROM recipes WHERE id=$1", recipeID)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	if ownerID != userID {
		http.Error(w, "Forbidden: not recipe owner", http.StatusForbidden)
		return
	}

	// Start transaction
	tx, err := DB.Beginx()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Delete ingredients and steps
	_, _ = tx.Exec(`DELETE FROM recipe_ingredients WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM recipe_steps WHERE recipe_id=$1`, recipeID)

	// Delete recipe
	_, err = tx.Exec(`DELETE FROM recipes WHERE id=$1`, recipeID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	// Commit
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": recipeID, "message": "Recipe deleted successfully"})
}
