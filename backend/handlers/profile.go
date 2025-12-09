package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Recipe struct for profile responses
type ProfileRecipe struct {
	ID              int     `json:"id" db:"id"`
	Title           string  `json:"title" db:"title"`
	Description     string  `json:"description" db:"description"`
	ThumbnailURL    *string `json:"thumbnail_url" db:"thumbnail_url"`
	Price           float64 `json:"price" db:"price"`
	PreparationTime int     `json:"preparation_time" db:"preparation_time"`
	CreatedAt       string  `json:"created_at" db:"created_at"`
}

// GetUserBookmarksHandler - GET /users/{id}/bookmarks
func GetUserBookmarksHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var recipes []ProfileRecipe
	err = DB.Select(&recipes, `
		SELECT r.id, r.title, r.description, r.thumbnail_url, r.price, r.preparation_time, r.created_at
		FROM recipes r
		INNER JOIN bookmarks b ON r.id = b.recipe_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`, userID)
	
	if err != nil {
		http.Error(w, "Failed to fetch bookmarks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// GetUserPurchasesHandler - GET /users/{id}/purchases
func GetUserPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var recipes []ProfileRecipe
	err = DB.Select(&recipes, `
		SELECT r.id, r.title, r.description, r.thumbnail_url, r.price, r.preparation_time, r.created_at
		FROM recipes r
		INNER JOIN purchases p ON r.id = p.recipe_id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`, userID)
	
	if err != nil {
		http.Error(w, "Failed to fetch purchases", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

