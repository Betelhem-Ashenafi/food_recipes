package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadImage uploads an image to Supabase Storage and returns the public URL
func UploadImage(filePath string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_API_KEY")
	supabaseBucket := os.Getenv("SUPABASE_BUCKET")

	if supabaseURL == "" || supabaseKey == "" || supabaseBucket == "" {
		return "", fmt.Errorf("Supabase environment variables not set")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}
	writer.Close()

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, supabaseBucket, filepath.Base(filePath))
	req, err := http.NewRequest("POST", uploadURL, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Supabase: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Supabase upload failed: %s", string(respBody))
	}

	// Public URL format for Supabase Storage
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, supabaseBucket, filepath.Base(filePath))
	return publicURL, nil
}
