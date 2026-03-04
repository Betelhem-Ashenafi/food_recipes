package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"foodrecipes/utils"

	"github.com/golang-jwt/jwt"
)

type HasuraUploadInput struct {
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
	Content  string `json:"content"` // base64
}

type HasuraUploadPayload struct {
	Input struct {
		Arg  *HasuraUploadInput `json:"arg"`
		File *HasuraUploadInput `json:"file"`
	} `json:"input"`
}

type HasuraUploadResponse struct {
	URL string `json:"url"`
}

// HasuraErrorResponse is already defined in auth.go; no need to redefine.
func HasuraUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Validate JWT token
	tokenStr := extractTokenFromHeader(r)
	if tokenStr == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header", "invalid_token")
		return
	}
	userID, err := validateAndExtractUserID(tokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", "invalid_token")
		return
	}
	_ = userID // optional logging

	// Parse JSON body
	var payload HasuraUploadPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_json")
		return
	}
	var input *HasuraUploadInput
	if payload.Input.Arg != nil {
		input = payload.Input.Arg
	} else if payload.Input.File != nil {
		input = payload.Input.File
	}
	if input == nil || input.Filename == "" || input.Mimetype == "" || input.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Missing upload input fields", "invalid_input")
		return
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(input.Content)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid base64 content", "invalid_base64")
		return
	}

	// Create a unique filename (preserve extension)
	ext := ""
	if dot := strings.LastIndex(input.Filename, "."); dot != -1 {
		ext = input.Filename[dot:]
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Upload to Cloudinary
	url, err := utils.UploadToCloudinary(r.Context(), bytes.NewReader(decoded), filename)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upload image: "+err.Error(), "cloudinary_error")
		return
	}

	// Return success
	json.NewEncoder(w).Encode(HasuraUploadResponse{URL: url})
}

// Helper functions (you may already have these in extra.go or auth.go)
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func validateAndExtractUserID(tokenStr string) (int, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	var userID float64
	if uid, ok := claims["user_id"]; ok {
		userID, _ = uid.(float64)
	} else if hasuraClaims, ok := claims["https://hasura.io/jwt/claims"].(map[string]interface{}); ok {
		if uid, ok := hasuraClaims["x-hasura-user-id"]; ok {
			userID, _ = uid.(float64)
		}
	}

	if userID == 0 {
		return 0, fmt.Errorf("user ID not found in token")
	}
	return int(userID), nil
}

func respondWithError(w http.ResponseWriter, status int, message, code string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(HasuraErrorResponse{
		Message: message,
		Code:    code,
	})
}
