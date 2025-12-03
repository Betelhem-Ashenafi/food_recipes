package handlers

import (
	"encoding/json"
	"foodrecipes/models"
	"net/http"
	"strconv"
	"strings"
)

// LikeRecipeHandler handles liking a recipe
func LikeRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL
	// URL pattern: /recipes/{id}/like
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2] // parts[0] is "", parts[1] is "recipes", parts[2] is id, parts[3] is "like"

	// If the URL is like /recipes/123/like, split gives ["", "recipes", "123", "like"]
	// Let's be more robust. The main.go passes the request.
	// If main.go strips prefix, we need to know.
	// In main.go: http.HandleFunc("/recipes/", ...)
	// So r.URL.Path is full path.

	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	_, err = DB.Exec("INSERT INTO likes (user_id, recipe_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, recipeID)
	if err != nil {
		http.Error(w, "Failed to like recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe liked"})
}

// UnlikeRecipeHandler handles unliking a recipe
func UnlikeRecipeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	_, err = DB.Exec("DELETE FROM likes WHERE user_id = $1 AND recipe_id = $2", userID, recipeID)
	if err != nil {
		http.Error(w, "Failed to unlike recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe unliked"})
}

// BookmarkRecipeHandler handles bookmarking a recipe
func BookmarkRecipeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	_, err = DB.Exec("INSERT INTO bookmarks (user_id, recipe_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, recipeID)
	if err != nil {
		http.Error(w, "Failed to bookmark recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe bookmarked"})
}

// UnbookmarkRecipeHandler handles unbookmarking a recipe
func UnbookmarkRecipeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	_, err = DB.Exec("DELETE FROM bookmarks WHERE user_id = $1 AND recipe_id = $2", userID, recipeID)
	if err != nil {
		http.Error(w, "Failed to unbookmark recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe unbookmarked"})
}

// CommentRecipeHandler handles adding a comment to a recipe
func CommentRecipeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// /recipes/{id}/comments
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var commentID int
	err = DB.QueryRow("INSERT INTO comments (user_id, recipe_id, content) VALUES ($1, $2, $3) RETURNING id", userID, recipeID, req.Content).Scan(&commentID)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Comment added",
		"id":      commentID,
	})
}

// GetCommentsHandler handles fetching comments for a recipe
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// /recipes/{id}/comments
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	query := `
		SELECT c.id, c.user_id, c.recipe_id, c.content, c.created_at, u.name as user_name
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.recipe_id = $1
		ORDER BY c.created_at DESC
	`
	err = DB.Select(&comments, query, recipeID)
	if err != nil {
		http.Error(w, "Failed to fetch comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// RateRecipeHandler handles rating a recipe (1-5 stars)
func RateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// /recipes/{id}/rate
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var req models.RateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	// Upsert rating: Insert or Update if exists
	query := `
		INSERT INTO ratings (user_id, recipe_id, rating)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, recipe_id)
		DO UPDATE SET rating = EXCLUDED.rating, created_at = CURRENT_TIMESTAMP
	`
	_, err = DB.Exec(query, userID, recipeID, req.Rating)
	if err != nil {
		http.Error(w, "Failed to rate recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe rated successfully"})
}

// GetRecipeRatingHandler returns the average rating and count for a recipe
func GetRecipeRatingHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// /recipes/{id}/rate
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeIDStr := parts[2]
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var result struct {
		AverageRating float64 `db:"avg_rating" json:"average_rating"`
		Count         int     `db:"count" json:"count"`
	}

	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count
		FROM ratings
		WHERE recipe_id = $1
	`
	err = DB.Get(&result, query, recipeID)
	if err != nil {
		http.Error(w, "Failed to get ratings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
