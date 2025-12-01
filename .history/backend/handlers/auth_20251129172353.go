package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/food-recipes-app/backend/models"
	"github.com/food-recipes-app/backend/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB // Set this in main.go

// LoginRequest is the expected input from Hasura
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the output for Hasura
type LoginResponse struct {
	Token string         `json:"token"`
	User  models.User    `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid request"}`))
		return
	}

	var user models.User
	err := DB.Get(&user, "SELECT id, email, password, name, avatar_url FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"user not found"}`))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"wrong password"}`))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"token error"}`))
		return
	}

	resp := LoginResponse{
		Token: token,
		User:  user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
