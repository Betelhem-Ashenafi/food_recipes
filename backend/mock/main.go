package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/transaction/verify/", func(w http.ResponseWriter, r *http.Request) {
		// Extract txRef
		parts := strings.Split(r.URL.Path, "/")
		txRef := parts[len(parts)-1]
		log.Printf("Verifying transaction: %s", txRef)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Payment verified",
			"data": map[string]interface{}{
				"status": "success",
				"amount": 100.00,
			},
		})
	})

	log.Println("Mock Chapa Server running on :8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
