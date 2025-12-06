package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// UploadFileHandler handles single file upload and returns the URL
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 10MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Ensure uploads directory exists
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Create destination file
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Construct URL
	// Assuming server runs on localhost:8081
	// In production, this should be configured via env var
	fileURL := fmt.Sprintf("http://localhost:8081/uploads/%s", filename)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"url": fileURL,
	})
}

// UploadRecipeImagesHandler allows uploading multiple image URLs for a recipe
func UploadRecipeImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/recipes/")
	idStr = strings.TrimSuffix(idStr, "/images")
	var recipeID int
	_, err := fmt.Sscanf(idStr, "%d", &recipeID)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from token
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	userID := int(token.Claims.(jwt.MapClaims)["user_id"].(float64))

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

	// Parse image URLs from request body
	var req struct {
		Images []string `json:"images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert images
	for _, url := range req.Images {
		_, err := DB.Exec("INSERT INTO recipe_images (recipe_id, url, is_featured) VALUES ($1, $2, false)", recipeID, url)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Images uploaded successfully"})
}

// FeatureRecipeImageHandler sets a featured image for a recipe
func FeatureRecipeImageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe and image ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/recipes/")
	parts := strings.Split(path, "/images/")
	if len(parts) != 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	var recipeID, imageID int
	_, err := fmt.Sscanf(parts[0], "%d", &recipeID)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}
	_, err = fmt.Sscanf(parts[1], "%d/feature", &imageID)
	if err != nil {
		// Try without /feature
		_, err = fmt.Sscanf(parts[1], "%d", &imageID)
		if err != nil {
			http.Error(w, "Invalid image ID", http.StatusBadRequest)
			return
		}
	}

	// Extract user ID from token
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	userID := int(token.Claims.(jwt.MapClaims)["user_id"].(float64))

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

	// Set all images to not featured
	_, err = DB.Exec("UPDATE recipe_images SET is_featured=false WHERE recipe_id=$1", recipeID)
	if err != nil {
		http.Error(w, "Failed to update images", http.StatusInternalServerError)
		return
	}
	// Set selected image to featured
	_, err = DB.Exec("UPDATE recipe_images SET is_featured=true WHERE id=$1 AND recipe_id=$2", imageID, recipeID)
	if err != nil {
		http.Error(w, "Failed to set featured image", http.StatusInternalServerError)
		return
	}

	// Optionally update recipe table with featured image id
	_, err = DB.Exec("UPDATE recipes SET thumbnail_url=(SELECT url FROM recipe_images WHERE id=$1) WHERE id=$2", imageID, recipeID)
	if err != nil {
		http.Error(w, "Failed to update recipe thumbnail", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Featured image set successfully"})
}
