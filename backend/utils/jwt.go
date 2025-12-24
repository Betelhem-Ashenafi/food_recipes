package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your-secret-key")

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
