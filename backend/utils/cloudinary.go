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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UploadToCloudinary uploads a file to Cloudinary using the REST API and returns the URL or error.
// Supports file as: string (path), []byte, *os.File, *multipart.FileHeader, or any io.Reader.
func UploadToCloudinary(ctx context.Context, file interface{}, filename string) (string, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	cloudURL := os.Getenv("CLOUDINARY_URL")

	if cloudURL != "" {
		// CLOUDINARY_URL format: cloudinary://API_KEY:API_SECRET@CLOUD_NAME
		u, err := url.Parse(cloudURL)
		if err == nil {
			if u.User != nil {
				if k := u.User.Username(); k != "" {
					apiKey = k
				}
				if s, ok := u.User.Password(); ok {
					apiSecret = s
				}
			}
			if h := u.Hostname(); h != "" {
				cloudName = h
			}
		}
	}

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("CLOUDINARY_CLOUD_NAME or API credentials not set")
	}

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
		if filename == "" {
			filename = filepath.Base(v.Name())
		}
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

	if closer != nil {
		defer closer.Close()
	}

	if filename == "" {
		filename = fmt.Sprintf("upload_%d", time.Now().Unix())
	}

	// Prepare multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, reader); err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// signature: sha1(public_id=<public_id>&timestamp=<timestamp><api_secret>)
	sigBase := fmt.Sprintf("public_id=%s&timestamp=%s", filename, timestamp)
	h := sha1.New()
	h.Write([]byte(sigBase + apiSecret))
	signature := hex.EncodeToString(h.Sum(nil))

	// add required fields
	_ = w.WriteField("api_key", apiKey)
	_ = w.WriteField("timestamp", timestamp)
	_ = w.WriteField("signature", signature)
	_ = w.WriteField("public_id", filename)

	if err := w.Close(); err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/auto/upload", cloudName)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var resBody bytes.Buffer
	if _, err := io.Copy(&resBody, resp.Body); err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("cloudinary upload failed: status %d: %s", resp.StatusCode, resBody.String())
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(resBody.Bytes(), &parsed); err != nil {
		return "", err
	}
	if s, ok := parsed["secure_url"].(string); ok && s != "" {
		return s, nil
	}
	if s, ok := parsed["url"].(string); ok && s != "" {
		return s, nil
	}
	return "", fmt.Errorf("no url in cloudinary response")
}
