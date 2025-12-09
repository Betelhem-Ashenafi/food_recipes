package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ToggleLikeHandler - POST/DELETE /recipes/{id}/like
func ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	// Extract recipe ID from URL: /recipes/{id}/like
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "like"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "like" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		// Add like
		_, err := DB.Exec(`
			INSERT INTO likes (user_id, recipe_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, recipe_id) DO NOTHING
		`, userID, recipeID, time.Now())
		
		if err != nil {
			http.Error(w, "Failed to add like", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "liked"})
		
	} else if r.Method == http.MethodDelete {
		// Remove like
		_, err := DB.Exec(`
			DELETE FROM likes
			WHERE user_id = $1 AND recipe_id = $2
		`, userID, recipeID)
		
		if err != nil {
			http.Error(w, "Failed to remove like", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "unliked"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ToggleBookmarkHandler - POST/DELETE /recipes/{id}/bookmark
func ToggleBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	// Extract recipe ID from URL: /recipes/{id}/bookmark
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "bookmark"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "bookmark" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		// Add bookmark
		_, err := DB.Exec(`
			INSERT INTO bookmarks (user_id, recipe_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, recipe_id) DO NOTHING
		`, userID, recipeID, time.Now())
		
		if err != nil {
			http.Error(w, "Failed to add bookmark", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "bookmarked"})
		
	} else if r.Method == http.MethodDelete {
		// Remove bookmark
		_, err := DB.Exec(`
			DELETE FROM bookmarks
			WHERE user_id = $1 AND recipe_id = $2
		`, userID, recipeID)
		
		if err != nil {
			http.Error(w, "Failed to remove bookmark", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "unbookmarked"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Comment struct
type Comment struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"`
	RecipeID  int       `json:"recipe_id" db:"recipe_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// GetCommentsHandler - GET /recipes/{id}/comments
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL: /recipes/{id}/comments
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "comments"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "comments" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var comments []Comment
	err = DB.Select(&comments, `
		SELECT c.id, c.user_id, u.name as user_name, c.recipe_id, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.recipe_id = $1
		ORDER BY c.created_at DESC
	`, recipeID)
	
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// PostCommentHandler - POST /recipes/{id}/comments
func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	// Extract recipe ID from URL: /recipes/{id}/comments
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "comments"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "comments" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var payload struct {
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	var commentID int
	err = DB.QueryRow(`
		INSERT INTO comments (user_id, recipe_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, recipeID, payload.Content, time.Now()).Scan(&commentID)
	
	if err != nil {
		http.Error(w, "Failed to post comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": commentID})
}

// RateRecipeHandler - POST /recipes/{id}/rate
func RateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	// Extract recipe ID from URL: /recipes/{id}/rate
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "rate"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "rate" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var payload struct {
		Rating int `json:"rating"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Rating < 1 || payload.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	_, err = DB.Exec(`
		INSERT INTO ratings (user_id, recipe_id, rating, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, recipe_id)
		DO UPDATE SET rating = $3, updated_at = $4
	`, userID, recipeID, payload.Rating, time.Now())
	
	if err != nil {
		http.Error(w, "Failed to save rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "rated"})
}

// GetRatingHandler - GET /recipes/{id}/rate
func GetRatingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL: /recipes/{id}/rate
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// pathParts should be: ["recipes", "{id}", "rate"]
	if len(pathParts) < 3 || pathParts[0] != "recipes" || pathParts[2] != "rate" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var result struct {
		AverageRating float64 `json:"average_rating" db:"avg"`
		Count         int     `json:"count" db:"count"`
	}
	
	err = DB.Get(&result, `
		SELECT 
			COALESCE(AVG(rating), 0) as avg,
			COUNT(*) as count
		FROM ratings
		WHERE recipe_id = $1
	`, recipeID)
	
	if err != nil {
		http.Error(w, "Failed to fetch rating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
