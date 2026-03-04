package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"foodrecipes/models"
	"foodrecipes/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// ...existing code...

var DB *sqlx.DB

// SetDB sets the database connection for the handlers package.
func SetDB(db *sqlx.DB) {
	DB = db
}

// Request/response structs for the Hasura actions

type HasuraLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HasuraLoginResponse struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type HasuraSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HasuraSignupResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type HasuraErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// HasuraLoginHandler handles the login action from Hasura
func HasuraLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read the whole body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid request body"})
		return
	}

	var req HasuraLoginRequest
	parsed := false

	// Try to parse as Hasura action payload (input wrapper)
	var wrapper struct {
		Input json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err == nil && len(wrapper.Input) > 0 {
		// First try input.arg format
		var argWrapper struct {
			Arg HasuraLoginRequest `json:"arg"`
		}
		if err := json.Unmarshal(wrapper.Input, &argWrapper); err == nil {
			if argWrapper.Arg.Email != "" || argWrapper.Arg.Password != "" {
				req = argWrapper.Arg
				parsed = true
			}
		}
		// Then try input direct format
		if !parsed {
			if err := json.Unmarshal(wrapper.Input, &req); err == nil {
				parsed = true
			}
		}
	}

	// Fallback: try plain JSON (for testing)
	if !parsed {
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid request format"})
			return
		}
	}

	// Fetch user from DB
	var user models.User
	err = DB.Get(&user, "SELECT id, name, email, password, COALESCE(avatar_url, '') as avatar_url FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid email or password", Code: "invalid_credentials"})
		return
	}

	// Compare password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid email or password", Code: "invalid_credentials"})
		return
	}

	// Generate JWT with Hasura claims
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Could not generate token"})
		return
	}

	// Return success
	json.NewEncoder(w).Encode(HasuraLoginResponse{
		Token:  token,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	})
}

// HasuraSignupHandler handles the signup action from Hasura
func HasuraSignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid request body"})
		return
	}

	var req HasuraSignupRequest
	parsed := false

	// Try to parse as Hasura action payload (input wrapper)
	var wrapper struct {
		Input json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err == nil && len(wrapper.Input) > 0 {
		// First try input.arg format
		var argWrapper struct {
			Arg HasuraSignupRequest `json:"arg"`
		}
		if err := json.Unmarshal(wrapper.Input, &argWrapper); err == nil {
			if argWrapper.Arg.Email != "" || argWrapper.Arg.Password != "" || argWrapper.Arg.Name != "" {
				req = argWrapper.Arg
				parsed = true
			}
		}
		// Then try input direct format
		if !parsed {
			if err := json.Unmarshal(wrapper.Input, &req); err == nil {
				parsed = true
			}
		}
	}

	// Fallback: plain JSON
	if !parsed {
		if err := json.Unmarshal(bodyBytes, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Invalid request format"})
			return
		}
	}

	// Debug: log received values
	log.Printf("Signup request: name=%q, email=%q", req.Name, req.Email)

	// Check if email already exists
	var count int
	err = DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1", req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Database error"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Email already registered", Code: "email_exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Could not hash password"})
		return
	}

	// Insert user
	var user models.User
	err = DB.Get(&user, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email
	`, req.Name, req.Email, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(HasuraErrorResponse{Message: "Could not create user"})
		return
	}

	// Return the new user
	json.NewEncoder(w).Encode(HasuraSignupResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}