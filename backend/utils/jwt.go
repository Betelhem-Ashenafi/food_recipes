package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback to Hasura-style JSON secret if present.
		hasuraSecret := os.Getenv("HASURA_GRAPHQL_JWT_SECRET")
		if hasuraSecret == "" {
			return nil, errors.New("JWT_SECRET is not set")
		}
		trimmed := strings.TrimSpace(hasuraSecret)
		if strings.HasPrefix(trimmed, "{") {
			var payload struct {
				Key string `json:"key"`
			}
			if err := json.Unmarshal([]byte(trimmed), &payload); err == nil && payload.Key != "" {
				return []byte(payload.Key), nil
			}
		}
		// Use raw value as a last resort.
		return []byte(hasuraSecret), nil
	}
	return []byte(secret), nil
}

// GetJWTSecret exposes the resolved JWT secret for other packages.
func GetJWTSecret() ([]byte, error) {
	return getJWTSecret()
}

// getTokenExpiration returns the token expiration duration
// Default: 7 days (168 hours)
// Can be overridden with JWT_EXPIRATION_HOURS environment variable
func getTokenExpiration() time.Duration {
	// Default to 7 days
	defaultExpiration := 7 * 24 * time.Hour

	// Check for environment variable
	if expHours := os.Getenv("JWT_EXPIRATION_HOURS"); expHours != "" {
		if hours, err := strconv.Atoi(expHours); err == nil && hours > 0 {
			return time.Duration(hours) * time.Hour
		}
	}

	return defaultExpiration
}

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID int, email string, name string) (string, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	expiration := getTokenExpiration()
	now := time.Now()
	expirationTime := now.Add(expiration)

	// Log token expiration for debugging
	log.Printf("[JWT] Generating token for user %d, expires in %v (at %v)",
		userID, expiration, expirationTime.Format(time.RFC3339))

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"name":    name,
		"iat":     now.Unix(), // Issued at time
		"exp":     expirationTime.Unix(),
		"https://hasura.io/jwt/claims": jwt.MapClaims{
			"x-hasura-allowed-roles": []string{"user"},
			"x-hasura-default-role":  "user",
			"x-hasura-user-id":       strconv.Itoa(userID),
			"x-hasura-user-name":     name,
			"x-hasura-user-email":    email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
