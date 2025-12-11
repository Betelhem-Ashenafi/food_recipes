package handlers

import (
	"encoding/json"
	"foodrecipes/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetRecipeIngredientsHandler returns ingredients for a recipe
func GetRecipeIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL: /recipes/{id}/ingredients
	// URL path: /recipes/1/ingredients
	// Split by "/" gives: ["", "recipes", "1", "ingredients"]
	path := strings.TrimPrefix(r.URL.Path, "/recipes/")
	path = strings.TrimSuffix(path, "/ingredients")
	recipeID, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("Error parsing recipe ID from path %s: %v", r.URL.Path, err)
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var ingredients []models.RecipeIngredient
	err = DB.Select(&ingredients, "SELECT id, recipe_id, name, quantity, unit FROM recipe_ingredients WHERE recipe_id=$1 ORDER BY id", recipeID)
	if err != nil {
		log.Printf("Error getting ingredients: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}

// GetRecipeStepsHandler returns steps for a recipe
func GetRecipeStepsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract recipe ID from URL: /recipes/{id}/steps
	// URL path: /recipes/1/steps
	// Split by "/" gives: ["", "recipes", "1", "steps"]
	path := strings.TrimPrefix(r.URL.Path, "/recipes/")
	path = strings.TrimSuffix(path, "/steps")
	recipeID, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("Error parsing recipe ID from path %s: %v", r.URL.Path, err)
		http.Error(w, "Invalid recipe ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	var steps []models.RecipeStep
	err = DB.Select(&steps, "SELECT id, recipe_id, step_number, instruction, COALESCE(image_url, '') as image_url FROM recipe_steps WHERE recipe_id=$1 ORDER BY step_number", recipeID)
	if err != nil {
		log.Printf("Error getting steps: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(steps)
}
