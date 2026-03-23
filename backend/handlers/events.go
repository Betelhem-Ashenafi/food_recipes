package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type EventTriggerPayload struct {
	Event struct {
		Op   string `json:"op"`
		Data struct {
			Old interface{} `json:"old"`
			New interface{} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}

type PaymentEventRecord struct {
	ID         int             `json:"id"`
	PurchaseID int             `json:"purchase_id"`
	TxRef      string          `json:"tx_ref"`
	UserID     int             `json:"user_id"`
	RecipeID   int             `json:"recipe_id"`
	OldStatus  *string         `json:"old_status"`
	NewStatus  string          `json:"new_status"`
	Payload    json.RawMessage `json:"payload"`
	CreatedAt  string          `json:"created_at"`
}

type PaymentEventTriggerPayload struct {
	Event struct {
		ID   string `json:"id"`
		Op   string `json:"op"`
		Data struct {
			Old *PaymentEventRecord `json:"old"`
			New *PaymentEventRecord `json:"new"`
		} `json:"data"`
	} `json:"event"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}

func NewRecipeEventHandler(w http.ResponseWriter, r *http.Request) {
	var payload EventTriggerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("New Recipe Created! Data: %+v", payload.Event.Data.New)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event received"))
}

// PaymentEventHandler receives Hasura event trigger payloads from payment_events table.
func PaymentEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var payload PaymentEventTriggerPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var raw map[string]interface{}
	_ = json.Unmarshal(body, &raw)

	if payload.Event.Data.New == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("No new event payload"))
		return
	}

	rec := payload.Event.Data.New

	// Normalize values so this handler works for either:
	// 1) payment_events table rows (purchase_id, tx_ref, new_status)
	// 2) purchases table rows (id, chapa_tx_ref, status)
	purchaseID := rec.PurchaseID
	if purchaseID == 0 {
		purchaseID = rec.ID
	}
	txRef := strings.TrimSpace(rec.TxRef)
	newStatus := strings.TrimSpace(strings.ToLower(rec.NewStatus))
	if newStatus == "" {
		newStatus = "(missing-new-status)"
	}

	oldStatus := ""
	if rec.OldStatus != nil {
		oldStatus = strings.TrimSpace(strings.ToLower(*rec.OldStatus))
	}

	// If core values are still empty, decode from raw event payload to support multiple table shapes.
	if purchaseID == 0 || txRef == "" || newStatus == "(missing-new-status)" {
		event, _ := raw["event"].(map[string]interface{})
		data, _ := event["data"].(map[string]interface{})
		newRec, _ := data["new"].(map[string]interface{})
		oldRec, _ := data["old"].(map[string]interface{})

		if purchaseID == 0 {
			if v, ok := newRec["purchase_id"]; ok {
				purchaseID = toInt(v)
			}
			if purchaseID == 0 {
				if v, ok := newRec["id"]; ok {
					purchaseID = toInt(v)
				}
			}
		}

		if txRef == "" {
			if v, ok := newRec["tx_ref"]; ok {
				txRef = strings.TrimSpace(toString(v))
			}
			if txRef == "" {
				if v, ok := newRec["chapa_tx_ref"]; ok {
					txRef = strings.TrimSpace(toString(v))
				}
			}
		}

		if newStatus == "(missing-new-status)" {
			if v, ok := newRec["new_status"]; ok {
				newStatus = strings.TrimSpace(strings.ToLower(toString(v)))
			}
			if newStatus == "" || newStatus == "(missing-new-status)" {
				if v, ok := newRec["status"]; ok {
					newStatus = strings.TrimSpace(strings.ToLower(toString(v)))
				}
			}
			if newStatus == "" {
				newStatus = "(missing-new-status)"
			}
		}

		if oldStatus == "" && oldRec != nil {
			if v, ok := oldRec["old_status"]; ok {
				oldStatus = strings.TrimSpace(strings.ToLower(toString(v)))
			}
			if oldStatus == "" && oldRec["status"] != nil {
				oldStatus = strings.TrimSpace(strings.ToLower(toString(oldRec["status"])))
			}
		}
	}

	currentDBStatus := ""
	if purchaseID > 0 {
		if status, err := currentPurchaseStatus(purchaseID); err == nil {
			currentDBStatus = status
		}
	}

	// Prefer explicit event status, but fall back to current DB status when event payload is incomplete.
	effectiveStatus := normalizeEventStatus(newStatus)
	if effectiveStatus == "unknown" {
		effectiveStatus = normalizeEventStatus(currentDBStatus)
	}

	if currentDBStatus != "" {
		log.Printf("[HASURA EVENT][PAYMENT %s] table=%s.%s op=%s purchase_id=%d tx_ref=%s old=%s new=%s current_db=%s", strings.ToUpper(effectiveStatus), payload.Table.Schema, payload.Table.Name, payload.Event.Op, purchaseID, txRef, oldStatus, newStatus, currentDBStatus)
	} else {
		log.Printf("[HASURA EVENT][PAYMENT %s] table=%s.%s op=%s purchase_id=%d tx_ref=%s old=%s new=%s", strings.ToUpper(effectiveStatus), payload.Table.Schema, payload.Table.Name, payload.Event.Op, purchaseID, txRef, oldStatus, newStatus)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Payment event received"))
}

func toInt(v interface{}) int {
	s := strings.TrimSpace(toString(v))
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.Itoa(int(t))
	case int:
		return strconv.Itoa(t)
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return strings.Trim(string(b), "\"")
	}
}

func currentPurchaseStatus(purchaseID int) (string, error) {
	if DB == nil {
		return "", nil
	}
	var status string
	if err := DB.Get(&status, `SELECT LOWER(COALESCE(status, '')) FROM purchases WHERE id = $1`, purchaseID); err != nil {
		return "", err
	}
	return strings.TrimSpace(status), nil
}

func normalizeEventStatus(s string) string {
	status := strings.ToLower(strings.TrimSpace(s))
	switch status {
	case "success", "paid", "completed":
		return "success"
	case "pending", "processing", "created", "initiated":
		return "pending"
	case "failed", "cancelled", "canceled", "error":
		return "failed"
	default:
		return "unknown"
	}
}
