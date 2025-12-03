package models

import "time"

type Like struct {
	UserID    int       `db:"user_id" json:"user_id"`
	RecipeID  int       `db:"recipe_id" json:"recipe_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Bookmark struct {
	UserID    int       `db:"user_id" json:"user_id"`
	RecipeID  int       `db:"recipe_id" json:"recipe_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Comment struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	RecipeID  int       `db:"recipe_id" json:"recipe_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserName  string    `db:"user_name" json:"user_name,omitempty"` // For display purposes
}

type Rating struct {
	UserID    int       `db:"user_id" json:"user_id"`
	RecipeID  int       `db:"recipe_id" json:"recipe_id"`
	Rating    int       `db:"rating" json:"rating"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type RateRecipeRequest struct {
	Rating int `json:"rating"`
}
