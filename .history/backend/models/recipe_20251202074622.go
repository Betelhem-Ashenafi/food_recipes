package models

import "time"

type Category struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	ImageURL string `db:"image_url" json:"image_url"`
}

type Recipe struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	CategoryID      int       `db:"category_id" json:"category_id"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	PreparationTime int       `db:"preparation_time" json:"preparation_time"` // in minutes
	Price           float64   `db:"price" json:"price"`
	ThumbnailURL    string    `db:"thumbnail_url" json:"thumbnail_url"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type RecipeIngredient struct {
	ID       int    `db:"id" json:"id"`
	RecipeID int    `db:"recipe_id" json:"recipe_id"`
	Name     string `db:"name" json:"name"`
	Quantity string `db:"quantity" json:"quantity"`
	Unit     string `db:"unit" json:"unit"`
}

type RecipeStep struct {
	ID          int    `db:"id" json:"id"`
	RecipeID    int    `db:"recipe_id" json:"recipe_id"`
	StepNumber  int    `db:"step_number" json:"step_number"`
	Instruction string `db:"instruction" json:"instruction"`
	ImageURL    string `db:"image_url" json:"image_url"`
}

// Request struct for creating a full recipe
type CreateRecipeRequest struct {
	CategoryID      int                `json:"category_id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	PreparationTime int                `json:"preparation_time"`
	Price           float64            `json:"price"`
	ThumbnailURL    string             `json:"thumbnail_url"`
	Ingredients     []RecipeIngredient `json:"ingredients"`
	Steps           []RecipeStep       `json:"steps"`
}
