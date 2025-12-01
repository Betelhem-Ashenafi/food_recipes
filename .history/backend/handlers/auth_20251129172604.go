package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/food-recipes-app/backend/models"
	"github.com/food-recipes-app/backend/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/jmoiron/sqlx"
)

// Struct for login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for login response
type LoginResponse struct {
	Token string         `json:"token"`
	User  models.User    `json:"user"`
	Error string         `json:"error,omitempty"`
}

var DB *sqlx.DB // Set this in main.go

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid request body"})
		return
	}

	var user models.User
	// Query user by email
	err := DB.Get(&user, "SELECT * FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid email or password"})
		return
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LoginResponse{Error: "Could not generate token"})
		return
	}

	user.Password = "" // Do not return password
	json.NewEncoder(w).Encode(LoginResponse{Token: token, User: user})
}
