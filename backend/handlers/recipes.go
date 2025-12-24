package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"foodrecipes/models"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

// contextKey is a custom type for context keys in this package
type contextKey string

var userIDKey = contextKey("user_id")

var DB *sqlx.DB

// validateTokenWithGracePeriod parses and validates a JWT token with a grace period for clock skew
// Returns the token, claims, and any error
func validateTokenWithGracePeriod(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// Parse token WITHOUT automatic expiration validation
	// We'll manually check expiration with a grace period to handle clock skew
	parser := jwt.Parser{
		SkipClaimsValidation: true, // Skip automatic expiration check
	}
	
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("token parse error: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token claims type")
	}

	// Manually validate expiration with 5-minute grace period for clock skew
	now := time.Now()
	gracePeriod := 5 * time.Minute
	
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		remaining := expTime.Sub(now)
		
		// Allow tokens that are expired but within grace period (for clock skew)
		if expTime.Before(now) && now.Sub(expTime) > gracePeriod {
			return nil, nil, fmt.Errorf("token expired at %v (beyond grace period)", expTime.Format(time.RFC3339))
		}
		
		// Log token expiration info for debugging (only in AuthMiddleware)
		if remaining > 0 {
			log.Printf("[AUTH] Token expires in %v (at %v)", remaining, expTime.Format(time.RFC3339))
		} else if now.Sub(expTime) <= gracePeriod {
			log.Printf("[AUTH] Token expired but within grace period: expired %v ago", now.Sub(expTime))
		}
	} else {
		return nil, nil, fmt.Errorf("token missing expiration claim")
	}

	return token, claims, nil
}

// Middleware to validate JWT and extract User ID
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH] Missing Authorization header for %s %s", r.Method, r.URL.Path)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			log.Printf("[AUTH] Empty token string for %s %s", r.Method, r.URL.Path)
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Use helper function to validate token with grace period
		token, claims, err := validateTokenWithGracePeriod(tokenString)
		if err != nil {
			log.Printf("[AUTH] Token validation error for %s %s: %v", r.Method, r.URL.Path, err)
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}
		
		_ = token // Token is valid, we just need claims

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("[AUTH] Invalid user_id in token for %s %s, claims: %v", r.Method, r.URL.Path, claims)
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		log.Printf("[AUTH] Authenticated user %d for %s %s", userID, r.Method, r.URL.Path)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
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
	
	// Use helper function to validate token with grace period
	_, claims, err := validateTokenWithGracePeriod(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		log.Printf("[CREATE] Invalid user_id type in token claims: %v", claims["user_id"])
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	userID := int(userIDFloat)
	log.Printf("[CREATE] Extracted user_id from token: %d", userID)

	// 2. Verify user exists in database
	var userExists bool
	err = DB.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", userID)
	if err != nil {
		log.Printf("[CREATE] Error checking user existence: %v", err)
		http.Error(w, "Database error while verifying user", http.StatusInternalServerError)
		return
	}
	if !userExists {
		log.Printf("[CREATE] User ID %d from token does not exist in database", userID)
		http.Error(w, "User not found. Please log in again.", http.StatusUnauthorized)
		return
	}

	// 3. Parse Request
	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Start Transaction
	tx, err := DB.Beginx()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 5. Insert Recipe
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

	// 6. Insert Ingredients
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

	// 7. Insert Steps
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

	// 8. Insert Images (if provided)
	if len(req.Images) > 0 {
		// Find featured image (thumbnail_url)
		featuredURL := req.ThumbnailURL
		for _, imgURL := range req.Images {
			if imgURL == "" {
				continue // Skip empty URLs
			}
			isFeatured := (imgURL == featuredURL)
			_, err = tx.Exec(`
				INSERT INTO recipe_images (recipe_id, url, is_featured)
				VALUES ($1, $2, $3)
			`, recipeID, imgURL, isFeatured)
			if err != nil {
				log.Printf("Error inserting image %s: %v", imgURL, err)
				tx.Rollback()
				http.Error(w, "Failed to save images: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		log.Printf("Saved %d images for recipe %d", len(req.Images), recipeID)
	} else if req.ThumbnailURL != "" {
		// If no images array but thumbnail_url exists, save it as featured
		_, err = tx.Exec(`
			INSERT INTO recipe_images (recipe_id, url, is_featured)
			VALUES ($1, $2, true)
		`, recipeID, req.ThumbnailURL)
		if err != nil {
			log.Printf("Error inserting thumbnail image: %v", err)
			tx.Rollback()
			http.Error(w, "Failed to save featured image: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Saved thumbnail image for recipe %d", recipeID)
	}

	// Commit transaction
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
	categoryParam := q.Get("category")

	log.Printf("DEBUG: Filters: time=%s, ingredient=%s, title=%s, creator=%s, category=%s", timeParam, ingredientParam, titleParam, creatorParam, categoryParam)

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

	// Filter by category (category_id)
	if categoryParam != "" {
		var categoryID int
		if _, err := fmt.Sscanf(categoryParam, "%d", &categoryID); err == nil {
			where = append(where, fmt.Sprintf("r.category_id = $%d", len(args)+1))
			args = append(args, categoryID)
		}
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

// GetRecipeByIDHandler returns a single recipe by ID
func GetRecipeByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL: /recipes/{id}
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

	var recipe models.Recipe
	err = DB.Get(&recipe, "SELECT id, user_id, category_id, title, description, preparation_time, price, thumbnail_url, created_at FROM recipes WHERE id=$1", recipeID)
	if err != nil {
		log.Printf("Error getting recipe %d: %v", recipeID, err)
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Fetch images for the recipe
	var images []models.RecipeImage
	err = DB.Select(&images, "SELECT id, recipe_id, url, is_featured FROM recipe_images WHERE recipe_id=$1", recipeID)
	if err == nil {
		recipe.Images = images
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
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
	w.Header().Set("Content-Type", "application/json")
	
	// Extract recipe ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/")
	// Remove any trailing slashes or paths
	if idx := strings.Index(idStr, "/"); idx != -1 {
		idStr = idStr[:idx]
	}
	if idStr == "" {
		http.Error(w, "Recipe ID required", http.StatusBadRequest)
		return
	}
	// Convert to int
	var recipeID int
	_, err := fmt.Sscanf(idStr, "%d", &recipeID)
	if err != nil {
		log.Printf("[EDIT] Error parsing recipe ID from path %s: %v", r.URL.Path, err)
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	
	log.Printf("[EDIT] Updating recipe ID: %d", recipeID)

	// Get user ID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

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
		log.Printf("[EDIT] Error decoding request body: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	log.Printf("[EDIT] Request data - Title: %s, CategoryID: %d, Ingredients: %d, Steps: %d", 
		req.Title, req.CategoryID, len(req.Ingredients), len(req.Steps))

	// Start transaction
	tx, err := DB.Beginx()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Update recipe
	result, err := tx.Exec(`UPDATE recipes SET category_id=$1, title=$2, description=$3, preparation_time=$4, price=$5, thumbnail_url=$6 WHERE id=$7`,
		req.CategoryID, req.Title, req.Description, req.PreparationTime, req.Price, req.ThumbnailURL, recipeID)
	if err != nil {
		tx.Rollback()
		log.Printf("[EDIT] Error updating recipe: %v", err)
		http.Error(w, "Failed to update recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("[EDIT] Recipe update affected %d rows", rowsAffected)
	if rowsAffected == 0 {
		tx.Rollback()
		log.Printf("[EDIT] Warning: No rows updated for recipe ID %d", recipeID)
		http.Error(w, "Recipe not found or no changes made", http.StatusNotFound)
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
			log.Printf("[EDIT] Error inserting ingredient: %v", err)
			http.Error(w, "Failed to update ingredients: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Printf("[EDIT] Inserted %d ingredients", len(req.Ingredients))

	// Insert new steps
	for i, step := range req.Steps {
		_, err = tx.Exec(`INSERT INTO recipe_steps (recipe_id, step_number, instruction, image_url) VALUES ($1, $2, $3, $4)`, recipeID, i+1, step.Instruction, step.ImageURL)
		if err != nil {
			tx.Rollback()
			log.Printf("[EDIT] Error inserting step %d: %v", i+1, err)
			http.Error(w, "Failed to update steps: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	log.Printf("[EDIT] Inserted %d steps", len(req.Steps))

	// Update images (if provided)
	if len(req.Images) > 0 {
		// Delete old images
		_, _ = tx.Exec(`DELETE FROM recipe_images WHERE recipe_id=$1`, recipeID)
		log.Printf("[EDIT] Deleted old images for recipe %d", recipeID)

		// Insert new images
		featuredURL := req.ThumbnailURL
		for _, imgURL := range req.Images {
			if imgURL == "" {
				continue // Skip empty URLs
			}
			isFeatured := (imgURL == featuredURL)
			_, err = tx.Exec(`
				INSERT INTO recipe_images (recipe_id, url, is_featured)
				VALUES ($1, $2, $3)
			`, recipeID, imgURL, isFeatured)
			if err != nil {
				tx.Rollback()
				log.Printf("[EDIT] Error inserting image %s: %v", imgURL, err)
				http.Error(w, "Failed to update images: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		log.Printf("[EDIT] Inserted %d images for recipe %d", len(req.Images), recipeID)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		log.Printf("[EDIT] Error committing transaction: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	log.Printf("[EDIT] Recipe %d updated successfully", recipeID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": recipeID, 
		"message": "Recipe updated successfully",
		"title": req.Title,
	})
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
	
	// Use helper function to validate token with grace period
	_, claims, err := validateTokenWithGracePeriod(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
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

	// Delete all related data (cascade delete)
	_, _ = tx.Exec(`DELETE FROM recipe_images WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM recipe_ingredients WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM recipe_steps WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM likes WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM bookmarks WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM comments WHERE recipe_id=$1`, recipeID)
	_, _ = tx.Exec(`DELETE FROM ratings WHERE recipe_id=$1`, recipeID)

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
