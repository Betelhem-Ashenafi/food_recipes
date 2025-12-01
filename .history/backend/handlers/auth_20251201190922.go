package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"foodrecipes/models"
	"foodrecipes/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Struct for login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
	Error string      `json:"error,omitempty"`
}

// Struct for signup request
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for signup response
type SignupResponse struct {
	User  models.User `json:"user"`
	Error string      `json:"error,omitempty"`
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

// SignupHandler handles user registration
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SignupResponse{Error: "Invalid request body"})
		return
	}

	// Check if email already exists
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{Error: "Database error"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(SignupResponse{Error: "Email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{Error: "Could not hash password"})
		return
	}

	// Insert new user into database
	var user models.User
	err = DB.Get(&user, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, COALESCE(avatar_url, '') as avatar_url
	`, req.Name, req.Email, string(hashedPassword))
	if err != nil {
		fmt.Println("Error creating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SignupResponse{Error: "Could not create user"})
		return
	}

	json.NewEncoder(w).Encode(SignupResponse{User: user})
}
