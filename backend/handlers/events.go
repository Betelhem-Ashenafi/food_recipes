package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
