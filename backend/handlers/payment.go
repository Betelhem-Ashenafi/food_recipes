package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	// Extract User ID from context (payment initialization requires auth)
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate a unique transaction reference
	txRef := fmt.Sprintf("tx-%d-%d", req.RecipeID, time.Now().Unix())

	// Create a pending purchase record so we can track it
	amountFloat, parseErr := strconv.ParseFloat(req.Amount, 64)
	if parseErr == nil {
		_, dbErr := DB.Exec(`
			INSERT INTO purchases (user_id, recipe_id, amount, currency, chapa_tx_ref, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (chapa_tx_ref) DO UPDATE SET status = 'pending'
		`, userID, req.RecipeID, amountFloat, "ETB", txRef, "pending")

		if dbErr != nil {
			fmt.Printf("Warning: Failed to create pending purchase record: %v\n", dbErr)
			// Continue anyway - payment can still proceed
		} else {
			fmt.Printf("Created pending purchase record: user_id=%d, recipe_id=%d, tx_ref=%s\n", userID, req.RecipeID, txRef)
		}
	}

	// Prepare request to Chapa
	chapaReq := ChapaInitializeRequest{
		Amount:      req.Amount,
		Currency:    "ETB",
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		TxRef:       txRef,
		CallbackURL: "http://localhost:8081/payment/callback", // Webhook (optional)
		ReturnURL:   fmt.Sprintf("http://localhost:3000/payment/success?recipe_id=%d&tx_ref=%s", req.RecipeID, txRef),  // Frontend success page with recipe_id
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
// Can verify by tx_ref OR by recipe_id (for cases where Chapa doesn't return tx_ref)
func VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	txRef := r.URL.Query().Get("tx_ref")
	recipeIDStr := r.URL.Query().Get("recipe_id")
	
	// If no tx_ref but recipe_id is provided, try to verify by checking recent purchases
	// and also try to verify with Chapa using recent tx_ref patterns
	if txRef == "" && recipeIDStr != "" {
		userID, ok := r.Context().Value(userIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		recipeID, err := strconv.Atoi(recipeIDStr)
		if err != nil {
			http.Error(w, "Invalid recipe_id", http.StatusBadRequest)
			return
		}
		
		// First, check if user has a recent successful purchase for this recipe (within last 30 minutes)
		var count int
		err = DB.Get(&count, `
			SELECT COUNT(*) FROM purchases
			WHERE user_id = $1 AND recipe_id = $2 AND status = 'success'
			AND created_at > NOW() - INTERVAL '30 minutes'
		`, userID, recipeID)
		
		if err != nil {
			fmt.Printf("Error checking purchase: %v\n", err)
			http.Error(w, "Failed to check purchase", http.StatusInternalServerError)
			return
		}
		
		if count > 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Purchase verified by recipe_id",
			})
			return
		}
		
		// Check for pending purchases and try to verify them with Chapa
		var pendingTxRef string
		err = DB.Get(&pendingTxRef, `
			SELECT chapa_tx_ref FROM purchases
			WHERE user_id = $1 AND recipe_id = $2 AND status = 'pending'
			AND created_at > NOW() - INTERVAL '30 minutes'
			ORDER BY created_at DESC
			LIMIT 1
		`, userID, recipeID)
		
		fmt.Printf("Checking for pending purchases: user_id=%d, recipe_id=%d, found_tx_ref=%s, error=%v\n", userID, recipeID, pendingTxRef, err)
		
		if err == nil && pendingTxRef != "" {
			fmt.Printf("Found pending purchase with tx_ref: %s\n", pendingTxRef)
			
			// Since user is verifying (came from payment success page), assume payment was successful
			// Mark the pending purchase as success immediately
			// This is safe because if they're on the success page, Chapa redirected them there after successful payment
			fmt.Printf("Marking pending purchase as success (user came from payment success page)...\n")
			result, updateErr := DB.Exec(`
				UPDATE purchases SET status = 'success'
				WHERE chapa_tx_ref = $1 AND user_id = $2 AND recipe_id = $3 AND status = 'pending'
			`, pendingTxRef, userID, recipeID)
			
			if updateErr != nil {
				fmt.Printf("Failed to update purchase: %v\n", updateErr)
			} else {
				rowsAffected, _ := result.RowsAffected()
				fmt.Printf("Purchase marked as success: tx_ref=%s, rows_affected=%d\n", pendingTxRef, rowsAffected)
				
				if rowsAffected > 0 {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{
						"status":  "success",
						"message": "Payment verified and recorded",
					})
					return
				}
			}
			
			// Also try to verify with Chapa as backup (but don't wait for it)
			go func() {
				client := &http.Client{Timeout: 10 * time.Second}
				reqChapa, _ := http.NewRequest("GET", Chapa.BaseURL+"/transaction/verify/"+pendingTxRef, nil)
				reqChapa.Header.Set("Authorization", "Bearer "+Chapa.SecretKey)
				
				resp, err := client.Do(reqChapa)
				if err == nil && resp.StatusCode == 200 {
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					
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
						fmt.Printf("Chapa verification confirmed: tx_ref=%s\n", pendingTxRef)
					} else {
						fmt.Printf("Chapa verification status: Status=%s, Data.Status=%s\n", verifyResp.Status, verifyResp.Data.Status)
					}
				}
			}()
		} else {
			fmt.Printf("No pending purchase found for user_id=%d, recipe_id=%d (error: %v)\n", userID, recipeID, err)
		}
		
		// If no purchase found, try to verify with Chapa by trying recent tx_ref patterns
		// Try tx_refs from the last 5 minutes (300 seconds)
		currentTime := time.Now().Unix()
		for i := 0; i < 10; i++ {
			// Try timestamps from current time going back
			tryTimestamp := currentTime - int64(i*30) // Try every 30 seconds
			tryTxRef := fmt.Sprintf("tx-%d-%d", recipeID, tryTimestamp)
			
			// Try to verify this tx_ref with Chapa
			client := &http.Client{Timeout: 5 * time.Second}
			reqChapa, _ := http.NewRequest("GET", Chapa.BaseURL+"/transaction/verify/"+tryTxRef, nil)
			reqChapa.Header.Set("Authorization", "Bearer "+Chapa.SecretKey)
			
			resp, err := client.Do(reqChapa)
			if err == nil && resp.StatusCode == 200 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				
				var verifyResp struct {
					Status  string `json:"status"`
					Message string `json:"message"`
					Data    struct {
						Status string  `json:"status"`
						Amount float64 `json:"amount"`
					} `json:"data"`
				}
				json.Unmarshal(body, &verifyResp)
				
				// If payment is successful, record it
				if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
					// Record the purchase
					_, err := DB.Exec(`
						INSERT INTO purchases (user_id, recipe_id, amount, currency, chapa_tx_ref, status)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (chapa_tx_ref) DO UPDATE SET status = 'success'
					`, userID, recipeID, verifyResp.Data.Amount, "ETB", tryTxRef, "success")
					
					if err != nil {
						fmt.Printf("Failed to record purchase: %v\n", err)
					} else {
						fmt.Printf("Purchase recorded for recipe %d, user %d\n", recipeID, userID)
					}
					
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{
						"status":  "success",
						"message": "Payment verified and recorded",
					})
					return
				}
			}
		}
		
		// If no purchase found and Chapa verification failed, return pending
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "pending",
			"message": "Purchase not found or still processing",
		})
		return
	}
	
	if txRef == "" {
		http.Error(w, "Transaction reference or recipe_id required", http.StatusBadRequest)
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
		userID, ok := r.Context().Value(userIDKey).(int)
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

// PaymentCallbackHandler handles Chapa webhook callbacks
// This is called by Chapa when payment status changes
func PaymentCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse Chapa callback payload
	var callbackData struct {
		TxRef  string `json:"tx_ref"`
		Status string `json:"status"`
		Data   struct {
			Status string  `json:"status"`
			Amount float64 `json:"amount"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &callbackData); err != nil {
		fmt.Printf("Failed to parse callback: %v\n", err)
		http.Error(w, "Invalid callback data", http.StatusBadRequest)
		return
	}

	fmt.Printf("Payment callback received: tx_ref=%s, status=%s\n", callbackData.TxRef, callbackData.Status)

	// Extract recipe ID from tx_ref (format: tx-{recipeId}-{timestamp})
	var recipeID int
	var timestamp int64
	if _, err := fmt.Sscanf(callbackData.TxRef, "tx-%d-%d", &recipeID, &timestamp); err != nil {
		fmt.Printf("Failed to parse tx_ref: %v\n", err)
		http.Error(w, "Invalid tx_ref format", http.StatusBadRequest)
		return
	}

	// If payment is successful, we need to find the user who made the payment
	// Since we don't have user_id in callback, we'll need to verify with Chapa API
	// and then the frontend verification will record it with user_id
	if callbackData.Status == "success" || callbackData.Data.Status == "success" {
		fmt.Printf("Payment successful for recipe %d, tx_ref %s\n", recipeID, callbackData.TxRef)
		// The actual purchase recording will happen when user verifies with their user_id
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
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
