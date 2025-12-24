package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// CheckLikeHandler - GET /recipes/{id}/like/check
func CheckLikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var count int
	err = DB.Get(&count, `
		SELECT COUNT(*) FROM likes
		WHERE user_id = $1 AND recipe_id = $2
	`, userID, recipeID)
	
	if err != nil {
		http.Error(w, "Failed to check like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"liked": count > 0})
}

// CheckBookmarkHandler - GET /recipes/{id}/bookmark/check
func CheckBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var count int
	err = DB.Get(&count, `
		SELECT COUNT(*) FROM bookmarks
		WHERE user_id = $1 AND recipe_id = $2
	`, userID, recipeID)
	
	if err != nil {
		http.Error(w, "Failed to check bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"bookmarked": count > 0})
}

// CheckPurchaseHandler - GET /recipes/{id}/purchase/check
func CheckPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var count int
	// Check for successful purchases only
	err = DB.Get(&count, `
		SELECT COUNT(*) FROM purchases
		WHERE user_id = $1 AND recipe_id = $2 AND status = 'success'
	`, userID, recipeID)
	
	if err != nil {
		http.Error(w, "Failed to check purchase", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"purchased": count > 0})
}

// CheckRatingHandler - GET /recipes/{id}/rate/check
func CheckRatingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	recipeID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var rating int
	err = DB.Get(&rating, `
		SELECT rating FROM ratings
		WHERE user_id = $1 AND recipe_id = $2
	`, userID, recipeID)
	
	if err != nil {
		// User hasn't rated yet
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"rated": false, "rating": 0})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"rated": true, "rating": rating})
}







