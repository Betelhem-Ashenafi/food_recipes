package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

// ==================== Configuration & Helpers ====================

type ChapaConfig struct {
	SecretKey string
	BaseURL   string
}

func getChapaConfig() ChapaConfig {
	return ChapaConfig{
		SecretKey: strings.TrimSpace(getEnv("CHAPA_SECRET_KEY", "")),
		BaseURL:   strings.TrimSuffix(getEnv("CHAPA_BASE_URL", "https://api.chapa.co/v1"), "/"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// URLBuilder constructs URLs from environment and request context.
type URLBuilder struct {
	apiURL      string
	frontendURL string
	req         *http.Request // optional, used to derive defaults
}

func NewURLBuilder(r *http.Request) *URLBuilder {
	apiURL := strings.TrimSuffix(getEnv("API_URL", ""), "/")
	if apiURL == "" {
		apiURL = getRequestBaseURL(r)
	}
	frontendURL := strings.TrimSuffix(getEnv("FRONTEND_URL", ""), "/")
	if frontendURL == "" {
		frontendURL = getFrontendBaseURL(r)
		if frontendURL == "" {
			frontendURL = "http://localhost:3000" // fallback for dev
		}
	}
	return &URLBuilder{
		apiURL:      apiURL,
		frontendURL: frontendURL,
		req:         r,
	}
}

func (b *URLBuilder) CallbackURL() string {
	if cb := strings.TrimSpace(getEnv("PAYMENT_CALLBACK_URL", "")); cb != "" {
		return cb
	}
	if b.apiURL != "" {
		return b.apiURL + "/hasura/payment/callback"
	}
	return ""
}

func (b *URLBuilder) ReturnURL(txRef string, recipeID int) string {
	q := url.Values{}
	q.Set("tx_ref", txRef)
	q.Set("recipe_id", strconv.Itoa(recipeID))
	return fmt.Sprintf("%s/payment/success?%s", b.frontendURL, q.Encode())
}

func (b *URLBuilder) ConfirmRedirectURL(recipeID int, txRef, status, message string) string {
	q := url.Values{}
	q.Set("recipe_id", strconv.Itoa(recipeID))
	q.Set("tx_ref", txRef)
	q.Set("status", status)
	if message != "" {
		q.Set("message", message)
	}
	return fmt.Sprintf("%s/payment/success?%s", b.frontendURL, q.Encode())
}

func getRequestBaseURL(r *http.Request) string {
	if r == nil {
		return ""
	}
	scheme := "http"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if r.TLS != nil {
		scheme = "https"
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	if host == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

func getFrontendBaseURL(r *http.Request) string {
	if r == nil {
		return ""
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		return strings.TrimSuffix(origin, "/")
	}
	return ""
}

// normalizePurchaseStatus maps various status strings to canonical values.
func normalizePurchaseStatus(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "success", "completed", "paid":
		return "success"
	case "pending", "processing", "created", "initiated":
		return "pending"
	default:
		return "failed"
	}
}

// parseHasuraInput extracts the action input from a Hasura Action request.
// It supports both direct {"input": {...}} and wrapped {"input": {"arg": {...}}} formats.
func parseHasuraInput[T any](body []byte) (T, map[string]interface{}, error) {
	var zero T
	var envelope struct {
		Input            json.RawMessage        `json:"input"`
		SessionVariables map[string]interface{} `json:"session_variables"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || len(envelope.Input) == 0 {
		return zero, nil, fmt.Errorf("invalid hasura action payload")
	}

	// Try wrapped format: {"arg": {...}}
	var argWrapper struct {
		Arg json.RawMessage `json:"arg"`
	}
	if err := json.Unmarshal(envelope.Input, &argWrapper); err == nil && len(argWrapper.Arg) > 0 && string(argWrapper.Arg) != "null" {
		var fromArg T
		if err := json.Unmarshal(argWrapper.Arg, &fromArg); err == nil {
			return fromArg, envelope.SessionVariables, nil
		}
	}

	// Try direct format
	var direct T
	if err := json.Unmarshal(envelope.Input, &direct); err == nil {
		return direct, envelope.SessionVariables, nil
	}

	return zero, nil, fmt.Errorf("invalid hasura action input")
}

// getUserIDFromSession extracts the user ID from Hasura session variables.
func getUserIDFromSession(session map[string]interface{}) (int, error) {
	raw, ok := session["x-hasura-user-id"]
	if !ok {
		return 0, fmt.Errorf("missing x-hasura-user-id")
	}
	switch v := raw.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("invalid user id type")
	}
}

// ==================== Service Layer ====================

type PaymentService struct {
	db         *sqlx.DB
	chapaCfg   ChapaConfig
	httpClient *http.Client
	logger     *log.Logger
}

func NewPaymentService(db *sqlx.DB, chapaCfg ChapaConfig, logger *log.Logger) *PaymentService {
	if logger == nil {
		logger = log.New(os.Stderr, "[payment] ", log.LstdFlags)
	}
	return &PaymentService{
		db:         db,
		chapaCfg:   chapaCfg,
		httpClient: &http.Client{Timeout: 20 * time.Second},
		logger:     logger,
	}
}

// ==================== Core Business Logic ====================

// InitializePayment starts a new payment or resumes an existing pending one.
func (s *PaymentService) InitializePayment(userID int, req *InitializePaymentRequest, urlBuilder *URLBuilder) (*InitializeResult, error) {
	// Validate input
	if err := s.validateInitializeRequest(userID, req); err != nil {
		return nil, err
	}

	// Check if already purchased
	if existing, _ := s.getSuccessfulPurchase(userID, req.RecipeID); existing != nil {
		return &InitializeResult{
			Status:  "success",
			Message: "Recipe already purchased",
			TxRef:   existing.TxRef,
		}, nil
	}

	// Check for pending purchase to resume
	if pending, _ := s.getPendingPurchase(userID, req.RecipeID); pending != nil && pending.CheckoutURL != "" {
		return &InitializeResult{
			Status:      "pending",
			Resumed:     true,
			CheckoutURL: pending.CheckoutURL,
			TxRef:       pending.TxRef,
		}, nil
	}

	// Create or update pending purchase record
	txRef := fmt.Sprintf("tx-%d-%d", req.RecipeID, time.Now().UnixNano())
	purchaseID, err := s.createPendingPurchase(userID, req, txRef)
	if err != nil {
		return nil, err
	}

	// Prepare Chapa request
	chapaReq := &ChapaInitializeRequest{
		Amount:    req.Amount,
		Currency:  "ETB",
		Email:     req.Email,
		FirstName: firstNameFromUserName(req.UserName),
		LastName:  "",
		TxRef:     txRef,
		ReturnURL: urlBuilder.ReturnURL(txRef, req.RecipeID),
	}
	if cb := urlBuilder.CallbackURL(); cb != "" {
		chapaReq.CallbackURL = cb
	}

	// Call Chapa API
	resp, err := s.callChapaInitialize(chapaReq)
	if err != nil {
		return nil, err
	}

	// Update purchase with checkout URL and provider payment ID
	if err := s.updatePurchaseWithProviderData(purchaseID, resp.Data.CheckoutURL); err != nil {
		s.logger.Printf("failed to update purchase with provider data: %v", err)
		// non-critical, continue
	}

	return &InitializeResult{
		Status:      "success",
		CheckoutURL: resp.Data.CheckoutURL,
		TxRef:       txRef,
		Resumed:     false,
	}, nil
}

// VerifyPayment checks payment status with Chapa and updates the database.
func (s *PaymentService) VerifyPayment(userID int, txRef string, recipeID int) (*VerifyResult, error) {
	// If only recipeID is given, try to find the latest transaction for that user+recipe
	if txRef == "" && recipeID > 0 {
		purchases, err := s.findPurchasesByUserAndRecipe(userID, recipeID)
		if err != nil {
			return nil, ErrNotFound
		}
		// If there's a successful one, return it immediately
		for _, p := range purchases {
			if p.Status == "success" {
				return &VerifyResult{
					Status:  "success",
					Message: "Payment already verified",
					TxRef:   p.TxRef,
				}, nil
			}
		}
		// Otherwise use the most recent txRef
		if len(purchases) > 0 && purchases[0].TxRef != "" {
			txRef = purchases[0].TxRef
		}
	}
	if txRef == "" {
		return nil, fmt.Errorf("tx_ref or recipe_id is required")
	}

	// Verify with Chapa
	status, amount, message, err := s.verifyChapaTransaction(txRef)
	if err != nil {
		return nil, err
	}
	s.logger.Printf("[PAYMENT VERIFY] tx_ref=%s recipe_id=%d chapa_status=%s amount=%.2f", txRef, recipeID, status, amount)

	// If recipeID not provided, try to parse from txRef
	if recipeID == 0 {
		if parsed, err := parseRecipeIDFromTxRef(txRef); err == nil {
			recipeID = parsed
		}
	}
	if recipeID == 0 {
		return nil, fmt.Errorf("unable to determine recipe_id")
	}

	// Update or insert purchase record
	if err := s.recordPurchase(userID, recipeID, txRef, status, amount); err != nil {
		return nil, err
	}
	if status == "success" {
		s.logger.Printf("[PAYMENT SUCCESS] user_id=%d recipe_id=%d tx_ref=%s status=%s", userID, recipeID, txRef, status)
	}

	return &VerifyResult{
		Status:  status,
		Message: message,
		TxRef:   txRef,
	}, nil
}

// ConfirmPayment is used by the redirect endpoint to verify and then redirect.
func (s *PaymentService) ConfirmPayment(txRef string, urlBuilder *URLBuilder) (redirectURL string, err error) {
	var purchase struct {
		ID       int    `db:"id"`
		RecipeID int    `db:"recipe_id"`
		Status   string `db:"status"`
	}
	err = s.db.Get(&purchase, `SELECT id, recipe_id, status FROM purchases WHERE chapa_tx_ref = $1`, txRef)
	if err != nil {
		return "", ErrNotFound
	}

	if purchase.Status == "success" {
		return urlBuilder.ConfirmRedirectURL(purchase.RecipeID, txRef, "success", "Payment already confirmed"), nil
	}

	status, amount, message, err := s.verifyChapaTransaction(txRef)
	if err != nil {
		return urlBuilder.ConfirmRedirectURL(purchase.RecipeID, txRef, "failed", message), nil
	}

	_, _ = s.db.Exec(`
		UPDATE purchases
		SET status = $1, amount = CASE WHEN $2 > 0 THEN $2 ELSE amount END
		WHERE chapa_tx_ref = $3
	`, status, amount, txRef)
	if status == "success" {
		s.logger.Printf("[PAYMENT SUCCESS] confirm tx_ref=%s status=%s", txRef, status)
	}

	return urlBuilder.ConfirmRedirectURL(purchase.RecipeID, txRef, status, message), nil
}

// HandleWebhook processes Chapa callbacks with signature verification.
func (s *PaymentService) HandleWebhook(body []byte, signature string) error {
	// Verify signature
	if !s.verifyWebhookSignature(body, signature) {
		return fmt.Errorf("invalid signature")
	}

	var callback struct {
		TxRef  string `json:"tx_ref"`
		Status string `json:"status"`
		Data   struct {
			Status string  `json:"status"`
			Amount float64 `json:"amount"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &callback); err != nil {
		return fmt.Errorf("invalid callback data: %v", err)
	}
	if callback.TxRef == "" {
		return nil // ignore if no tx_ref
	}

	status := normalizePurchaseStatus(callback.Data.Status)
	if status == "" {
		status = normalizePurchaseStatus(callback.Status)
	}
	_, err := s.db.Exec(`
		UPDATE purchases
		SET status = $1,
		    amount = CASE WHEN $2 > 0 THEN $2 ELSE amount END
		WHERE chapa_tx_ref = $3
	`, status, callback.Data.Amount, callback.TxRef)
	if err == nil && status == "success" {
		s.logger.Printf("[PAYMENT SUCCESS] webhook tx_ref=%s status=%s", callback.TxRef, status)
	}
	return err
}

// ==================== Internal helpers ====================

func (s *PaymentService) validateInitializeRequest(userID int, req *InitializePaymentRequest) error {
	if req.Amount == "" || req.Email == "" || req.RecipeID == 0 {
		return fmt.Errorf("missing required fields")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email")
	}
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("invalid amount")
	}
	// Check user exists
	var exists bool
	if err := s.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID); err != nil {
		return fmt.Errorf("failed to validate user: %v", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	if err := s.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM recipes WHERE id = $1)`, req.RecipeID); err != nil {
		return fmt.Errorf("failed to validate recipe: %v", err)
	}
	if !exists {
		return fmt.Errorf("recipe not found")
	}
	return nil
}

type purchaseInfo struct {
	ID          int    `db:"id"`
	TxRef       string `db:"chapa_tx_ref"`
	Status      string `db:"status"`
	CheckoutURL string `db:"checkout_url"`
}

func (s *PaymentService) getSuccessfulPurchase(userID, recipeID int) (*purchaseInfo, error) {
	var p purchaseInfo
	err := s.db.Get(&p, `
		SELECT id, chapa_tx_ref, status, checkout_url
		FROM purchases
		WHERE user_id = $1 AND recipe_id = $2 AND status = 'success'
		ORDER BY created_at DESC LIMIT 1
	`, userID, recipeID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PaymentService) getPendingPurchase(userID, recipeID int) (*purchaseInfo, error) {
	var p purchaseInfo
	err := s.db.Get(&p, `
		SELECT id, chapa_tx_ref, status, checkout_url
		FROM purchases
		WHERE user_id = $1 AND recipe_id = $2 AND status = 'pending'
		ORDER BY created_at DESC LIMIT 1
	`, userID, recipeID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PaymentService) findPurchasesByUserAndRecipe(userID, recipeID int) ([]purchaseInfo, error) {
	var purchases []purchaseInfo
	err := s.db.Select(&purchases, `
		SELECT id, chapa_tx_ref, status, checkout_url
		FROM purchases
		WHERE user_id = $1 AND recipe_id = $2
		ORDER BY CASE WHEN status = 'success' THEN 0 ELSE 1 END, created_at DESC
	`, userID, recipeID)
	return purchases, err
}

func (s *PaymentService) createPendingPurchase(userID int, req *InitializePaymentRequest, txRef string) (int, error) {
	amount, _ := strconv.ParseFloat(req.Amount, 64)
	var purchaseID int
	err := s.db.Get(&purchaseID, `
		INSERT INTO purchases (user_id, recipe_id, amount, currency, chapa_tx_ref, status, checkout_url)
		VALUES ($1, $2, $3, $4, $5, $6, NULL)
		ON CONFLICT (user_id, recipe_id) DO UPDATE
		SET amount = EXCLUDED.amount,
		    currency = EXCLUDED.currency,
		    chapa_tx_ref = EXCLUDED.chapa_tx_ref,
		    status = 'pending',
		    checkout_url = NULL,
		    created_at = CURRENT_TIMESTAMP
		RETURNING id
	`, userID, req.RecipeID, amount, "ETB", txRef, "pending")
	return purchaseID, err
}

func (s *PaymentService) updatePurchaseWithProviderData(purchaseID int, checkoutURL string) error {
	_, err := s.db.Exec(`
		UPDATE purchases
		SET checkout_url = $1
		WHERE id = $2
	`, checkoutURL, purchaseID)
	return err
}

func (s *PaymentService) callChapaInitialize(req *ChapaInitializeRequest) (*ChapaInitializeResponse, error) {
	data, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", s.chapaCfg.BaseURL+"/transaction/initialize", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.chapaCfg.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Chapa: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chapaResp ChapaInitializeResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to parse Chapa response: %v", err)
	}
	if chapaResp.Status != "success" || chapaResp.Data.CheckoutURL == "" {
		msg := stringFromAny(chapaResp.Message)
		return nil, fmt.Errorf("Chapa initialization failed: %s", msg)
	}
	return &chapaResp, nil
}

func (s *PaymentService) verifyChapaTransaction(txRef string) (status string, amount float64, message string, err error) {
	httpReq, err := http.NewRequest("GET", s.chapaCfg.BaseURL+"/transaction/verify/"+url.PathEscape(txRef), nil)
	if err != nil {
		return "failed", 0, "failed to create verification request", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.chapaCfg.SecretKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return "failed", 0, "failed to connect to Chapa", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var verifyResp ChapaVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return "failed", 0, "failed to parse verification response", err
	}

	msg := stringFromAny(verifyResp.Message)
	if verifyResp.Status != "success" {
		if msg == "" {
			msg = "verification failed"
		}
		return "failed", 0, msg, nil
	}

	status = normalizePurchaseStatus(verifyResp.Data.Status)
	if msg == "" {
		msg = "payment status: " + status
	}
	return status, verifyResp.Data.Amount, msg, nil
}

func (s *PaymentService) recordPurchase(userID, recipeID int, txRef, status string, amount float64) error {
	_, err := s.db.Exec(`
		INSERT INTO purchases (user_id, recipe_id, amount, currency, chapa_tx_ref, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (chapa_tx_ref) DO UPDATE
		SET status = EXCLUDED.status,
		    amount = CASE WHEN EXCLUDED.amount > 0 THEN EXCLUDED.amount ELSE purchases.amount END,
		    currency = EXCLUDED.currency
	`, userID, recipeID, amount, "ETB", txRef, status)
	return err
}

func (s *PaymentService) verifyWebhookSignature(body []byte, signatureHeader string) bool {
	// Chapa uses HMAC-SHA256 with the secret key.
	// The signature is in the header "x-chapa-signature".
	if s.chapaCfg.SecretKey == "" || signatureHeader == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(s.chapaCfg.SecretKey))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signatureHeader), []byte(expected))
}

// ==================== Request/Response Types ====================

type InitializePaymentRequest struct {
	Amount   string `json:"amount"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	RecipeID int    `json:"recipe_id"`
}

type VerifyPaymentRequest struct {
	TxRef    string `json:"tx_ref"`
	RecipeID int    `json:"recipe_id"`
}

type ChapaInitializeRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TxRef       string `json:"tx_ref"`
	CallbackURL string `json:"callback_url,omitempty"`
	ReturnURL   string `json:"return_url"`
}

type ChapaInitializeResponse struct {
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
	Data    struct {
		CheckoutURL string      `json:"checkout_url"`
		ID          interface{} `json:"id"`
	} `json:"data"`
}

type ChapaVerifyResponse struct {
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
	Data    struct {
		Status string      `json:"status"`
		Amount float64     `json:"amount"`
		ID     interface{} `json:"id"`
	} `json:"data"`
}

type InitializeResult struct {
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	Resumed     bool   `json:"resumed,omitempty"`
	CheckoutURL string `json:"checkout_url,omitempty"`
	TxRef       string `json:"tx_ref,omitempty"`
}

type VerifyResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	TxRef   string `json:"tx_ref,omitempty"`
}

// ==================== HTTP Handlers ====================

// withRecovery is middleware that recovers from panics and logs them.
func withRecovery(next http.HandlerFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Printf("panic: %v", rec)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// writeError sends a JSON error response in Hasura Action format.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// InitializePaymentHandler handles Hasura Action for initializing payment.
func InitializePaymentHandler(svc *PaymentService) http.HandlerFunc {
	return withRecovery(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		req, session, err := parseHasuraInput[InitializePaymentRequest](body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid hasura action payload")
			return
		}

		userID, err := getUserIDFromSession(session)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		urlBuilder := NewURLBuilder(r)
		result, err := svc.InitializePayment(userID, &req, urlBuilder)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}, svc.logger)
}

// VerifyPaymentHandler handles Hasura Action for verifying payment.
func VerifyPaymentHandler(svc *PaymentService) http.HandlerFunc {
	return withRecovery(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		req, session, err := parseHasuraInput[VerifyPaymentRequest](body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid hasura action payload")
			return
		}

		userID, err := getUserIDFromSession(session)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		result, err := svc.VerifyPayment(userID, req.TxRef, req.RecipeID)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}, svc.logger)
}

// PaymentCallbackHandler handles Chapa webhook callbacks.
func PaymentCallbackHandler(svc *PaymentService) http.HandlerFunc {
	return withRecovery(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		signature := r.Header.Get("x-chapa-signature")
		if err := svc.HandleWebhook(body, signature); err != nil {
			svc.logger.Printf("webhook processing error: %v", err)
			http.Error(w, "invalid signature or data", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "received"})
	}, svc.logger)
}

// ConfirmPaymentHandler handles the redirect endpoint after payment.
func ConfirmPaymentHandler(svc *PaymentService) http.HandlerFunc {
	return withRecovery(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 || parts[0] != "payment" || parts[2] != "confirm" {
			http.NotFound(w, r)
			return
		}

		txRef, err := url.PathUnescape(parts[1])
		if err != nil || txRef == "" {
			http.Error(w, "invalid tx_ref", http.StatusBadRequest)
			return
		}
		urlBuilder := NewURLBuilder(r)
		redirectURL, err := svc.ConfirmPayment(txRef, urlBuilder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}, svc.logger)
}

// ==================== Utility Functions ====================

func stringFromAny(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func firstNameFromUserName(userName string) string {
	name := strings.TrimSpace(userName)
	if name == "" {
		return "Customer"
	}
	parts := strings.Fields(name)
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return "Customer"
	}
	return parts[0]
}

func parseRecipeIDFromTxRef(txRef string) (int, error) {
	var recipeID int
	var ts int64
	_, err := fmt.Sscanf(txRef, "tx-%d-%d", &recipeID, &ts)
	return recipeID, err
}

// NewDefaultPaymentService creates a payment service using environment-based Chapa config.
func NewDefaultPaymentService(db *sqlx.DB, logger *log.Logger) *PaymentService {
	return NewPaymentService(db, getChapaConfig(), logger)
}
