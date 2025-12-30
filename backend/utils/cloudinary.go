package utils

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UploadToCloudinary uploads a file to Cloudinary using the REST API and returns the URL or error.
// It supports unsigned uploads via CLOUDINARY_UPLOAD_PRESET or signed uploads using CLOUDINARY_API_KEY and CLOUDINARY_API_SECRET.
func UploadToCloudinary(ctx context.Context, file interface{}, filename string) (string, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")

	if cloudName == "" {
		return "", fmt.Errorf("CLOUDINARY_CLOUD_NAME not set")
	}

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/auto/upload", cloudName)

	var reader io.Reader
	var closer io.Closer

	switch v := file.(type) {
	case nil:
		return "", fmt.Errorf("file is nil")
	case string:
		f, err := os.Open(v)
		if err != nil {
			return "", err
		}
		reader = f
		closer = f
		if filename == "" {
			filename = filepath.Base(v)
		}
	case []byte:
		reader = bytes.NewReader(v)
	case *os.File:
		reader = v
		closer = v
	case *multipart.FileHeader:
		f, err := v.Open()
		if err != nil {
			return "", err
		}
		reader = f
		closer = f
		if filename == "" {
			filename = v.Filename
		}
	default:
		if r, ok := v.(io.Reader); ok {
			reader = r
		} else {
			return "", fmt.Errorf("unsupported file type %T", file)
		}
	}

	if filename == "" {
		filename = fmt.Sprintf("upload_%d", time.Now().Unix())
	}

	if closer != nil {
		defer closer.Close()
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, reader); err != nil {
		mw.Close()
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if uploadPreset != "" {
		_ = mw.WriteField("upload_preset", uploadPreset)
	} else {
		_ = mw.WriteField("timestamp", timestamp)
		if filename != "" {
			_ = mw.WriteField("public_id", filename)
		}
		// create signature: sha1(<params_to_sign> + api_secret)
		var toSign string
		if filename != "" {
			toSign = fmt.Sprintf("public_id=%s&timestamp=%s", filename, timestamp)
		} else {
			toSign = fmt.Sprintf("timestamp=%s", timestamp)
		}
		toSign = toSign + apiSecret
		h := sha1.New()
		h.Write([]byte(toSign))
		signature := hex.EncodeToString(h.Sum(nil))
		_ = mw.WriteField("api_key", apiKey)
		_ = mw.WriteField("signature", signature)
	}

	mw.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(b))
	}

	var res struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.SecureURL != "" {
		return res.SecureURL, nil
	}
	if res.URL != "" {
		return res.URL, nil
	}

	return "", fmt.Errorf("no url in cloudinary response")
}
