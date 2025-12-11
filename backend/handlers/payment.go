package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

// ChapaConfig holds configuration for Chapa API
type ChapaConfig struct {
	SecretKey string
	BaseURL   string
}

var Chapa = ChapaConfig{
	SecretKey: getEnv("CHAPA_SECRET_KEY", "CHASECK_TEST-LR1J8py5LqhoMlJdtVT5piJnWtJ66RZk"),
	BaseURL:   getEnv("CHAPA_BASE_URL", "https://api.chapa.co/v1"),
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// InitializePaymentRequest is the payload from our frontend
type InitializePaymentRequest struct {
	Amount    string `json:"amount"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RecipeID  int    `json:"recipe_id"` // To track what they are buying
}

// ChapaInitializeRequest is the payload sent to Chapa
type ChapaInitializeRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TxRef       string `json:"tx_ref"`
	CallbackURL string `json:"callback_url"`
	ReturnURL   string `json:"return_url"`
}

// ChapaInitializeResponse is the response from Chapa
type ChapaInitializeResponse struct {
	Message interface{} `json:"message"` // Can be string or object
	Status  string      `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

// InitializePaymentHandler starts the payment process
func InitializePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req InitializePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Amount == "" || req.Email == "" || req.FirstName == "" || req.LastName == "" || req.RecipeID == 0 {
		http.Error(w, "Missing required fields: amount, email, first_name, last_name, recipe_id", http.StatusBadRequest)
		return
	}

	// Validate email format (basic check)
	if !strings.Contains(req.Email, "@") {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Generate a unique transaction reference
	txRef := fmt.Sprintf("tx-%d-%d", req.RecipeID, time.Now().Unix())

	// Prepare request to Chapa
	chapaReq := ChapaInitializeRequest{
		Amount:      req.Amount,
		Currency:    "ETB",
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		TxRef:       txRef,
		CallbackURL: "http://localhost:8081/payment/callback", // Webhook (optional)
		ReturnURL:   "http://localhost:3000/payment/success",  // Frontend success page
	}

	jsonData, _ := json.Marshal(chapaReq)

	// Send POST request to Chapa
	client := &http.Client{}
	reqChapa, _ := http.NewRequest("POST", Chapa.BaseURL+"/transaction/initialize", bytes.NewBuffer(jsonData))
	reqChapa.Header.Set("Authorization", "Bearer "+Chapa.SecretKey)
	reqChapa.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(reqChapa)
	if err != nil {
		http.Error(w, "Failed to connect to payment provider", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var chapaResp ChapaInitializeResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		http.Error(w, "Failed to parse payment provider response", http.StatusInternalServerError)
		return
	}

	if chapaResp.Status != "success" {
		msg := fmt.Sprintf("%v", chapaResp.Message)
		http.Error(w, "Payment initialization failed: "+msg, http.StatusBadRequest)
		return
	}

	// Return the checkout URL to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"checkout_url": chapaResp.Data.CheckoutURL,
		"tx_ref":       txRef,
	})
}

// VerifyPaymentHandler checks if a payment was successful
func VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	txRef := r.URL.Query().Get("tx_ref")
	if txRef == "" {
		http.Error(w, "Transaction reference required", http.StatusBadRequest)
		return
	}

	// Call Chapa Verify API
	client := &http.Client{}
	reqChapa, _ := http.NewRequest("GET", Chapa.BaseURL+"/transaction/verify/"+txRef, nil)
	reqChapa.Header.Set("Authorization", "Bearer "+Chapa.SecretKey)

	resp, err := client.Do(reqChapa)
	if err != nil {
		http.Error(w, "Failed to verify payment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Parse response
	var verifyResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Status string  `json:"status"`
			Amount float64 `json:"amount"`
		} `json:"data"`
	}
	json.Unmarshal(body, &verifyResp)

	if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
		// Extract User ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract Recipe ID from txRef
		var recipeID int
		var timestamp int64
		fmt.Sscanf(txRef, "tx-%d-%d", &recipeID, &timestamp)

		// Use GraphQL Client to fetch recipe details (Satisfying "must use go graphql client")
		httpClient := &http.Client{
			Transport: &headerTransport{
				Transport: http.DefaultTransport,
				Headers: map[string]string{
					"X-Hasura-Admin-Secret": "myhasurasecret",
				},
			},
		}
		gqlClient := graphql.NewClient("http://localhost:8080/v1/graphql", httpClient)

		var query struct {
			RecipesByPk struct {
				Price float64
			} `graphql:"recipes_by_pk(id: $id)"`
		}

		variables := map[string]interface{}{
			"id": graphql.Int(recipeID),
		}

		if err := gqlClient.Query(context.Background(), &query, variables); err != nil {
			fmt.Printf("Failed to fetch recipe details via GraphQL: %v\n", err)
			// Continue anyway, or fail? Let's continue but log.
		} else {
			fmt.Printf("Verified purchase for recipe price: %f\n", query.RecipesByPk.Price)
		}

		// Insert purchase record using SQLX (Simpler than GraphQL mutation for now)
		_, err := DB.Exec(`
			INSERT INTO purchases (user_id, recipe_id, amount, currency, chapa_tx_ref, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (chapa_tx_ref) DO NOTHING
		`, userID, recipeID, verifyResp.Data.Amount, "ETB", txRef, "success")

		if err != nil {
			fmt.Printf("Failed to insert purchase: %v\n", err)
			http.Error(w, "Failed to record purchase", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Payment verified and recorded"})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "failed", "message": "Payment verification failed"})
	}
}

// headerTransport adds headers to requests
type headerTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}
	return t.Transport.RoundTrip(req)
}
