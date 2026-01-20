package hmac

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kerimovok/go-pkg-utils/crypto"
)

const (
	// HeaderSignature is the header containing the HMAC signature
	HeaderSignature = "X-Signature"
	// HeaderTimestamp is the header containing the request timestamp
	HeaderTimestamp = "X-Timestamp"
)

// Client handles HMAC-authenticated HTTP requests
type Client struct {
	BaseURL    string
	HMACSecret string
	HTTPClient *http.Client
}

// Config holds configuration for HMAC client
type Config struct {
	BaseURL    string
	HMACSecret string
	Timeout    time.Duration
}

// NewClient creates a new HMAC HTTP client
func NewClient(config Config) *Client {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		BaseURL:    config.BaseURL,
		HMACSecret: config.HMACSecret,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ComputeSignature computes the HMAC-SHA256 signature for an HTTP request
// The signature is computed as: HMAC-SHA256(method + path + query + timestamp + body, secret)
func ComputeSignature(method, path, query, timestamp string, body []byte, secret string) string {
	// Build canonical message: method + path + query + timestamp + body
	message := method
	message += path
	if query != "" {
		message += "?" + query
	}
	message += timestamp
	message += string(body)

	// Compute HMAC-SHA256
	hash := crypto.HMACSHA256([]byte(message), []byte(secret))
	return hex.EncodeToString(hash)
}

// DoRequest makes an HMAC-authenticated HTTP request
func (c *Client) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.BaseURL + path

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req.Header.Set(HeaderTimestamp, timestamp)

	// Compute signature
	query := req.URL.RawQuery
	signature := ComputeSignature(method, path, query, timestamp, bodyBytes, c.HMACSecret)
	req.Header.Set(HeaderSignature, signature)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// DoRequestWithBody makes an HMAC-authenticated HTTP request with raw body bytes
func (c *Client) DoRequestWithBody(method, path string, bodyBytes []byte) (*http.Response, error) {
	url := c.BaseURL + path

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	req.Header.Set(HeaderTimestamp, timestamp)

	// Compute signature
	query := req.URL.RawQuery
	signature := ComputeSignature(method, path, query, timestamp, bodyBytes, c.HMACSecret)
	req.Header.Set(HeaderSignature, signature)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// ParseJSONResponse parses a JSON response from an HTTP response
func ParseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// ValidateSignature validates an HMAC signature
func ValidateSignature(method, path, query, timestamp string, body []byte, signature, secret string) bool {
	expectedSignature := ComputeSignature(method, path, query, timestamp, body, secret)
	signatureBytes, _ := hex.DecodeString(signature)
	expectedBytes, _ := hex.DecodeString(expectedSignature)
	return hmac.Equal(signatureBytes, expectedBytes)
}
