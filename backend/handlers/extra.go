package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"foodrecipes/utils"
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

type hasuraActionEnvelope struct {
	Input            json.RawMessage        `json:"input"`
	SessionVariables map[string]interface{} `json:"session_variables"`
}

// HasuraErrorResponse is already defined in auth.go; no need to redefine.
func HasuraUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", "invalid_body")
		return
	}

	var envelope hasuraActionEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil || len(envelope.Input) == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_json")
		return
	}

	rawUserID, ok := envelope.SessionVariables["x-hasura-user-id"]
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Missing session user id", "invalid_session")
		return
	}
	if _, err := parseUserID(rawUserID); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid session user id", "invalid_session")
		return
	}

	var payload HasuraUploadPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid upload payload", "invalid_input")
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

func parseUserID(raw interface{}) (int, error) {
	switch v := raw.(type) {
	case float64:
		if v <= 0 {
			return 0, fmt.Errorf("invalid user id")
		}
		return int(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			return 0, fmt.Errorf("invalid user id")
		}
		return id, nil
	case int:
		if v <= 0 {
			return 0, fmt.Errorf("invalid user id")
		}
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported user id type")
	}
}

func respondWithError(w http.ResponseWriter, status int, message, code string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(HasuraErrorResponse{
		Message: message,
		Code:    code,
	})
}
