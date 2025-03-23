// Package postalclient provides a Go client for the Postal API.
//
// Postal is a mail delivery platform that allows you to send emails through
// its API. This package provides a simple and idiomatic way to interact with
// the Postal API from Go applications.
//
// Basic usage:
//
//	client := postalclient.NewClient("your-api-key")
//
//	// Send a message
//	req := &models.SendMessageRequest{
//	    To:        []string{"recipient@example.com"},
//	    From:      "sender@yourdomain.com",
//	    Subject:   "Hello from Postal API",
//	    PlainBody: "This is a test email",
//	}
//	resp, err := client.SendMessage(req)
package postalclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Postal API.
	// This should be replaced with your Postal server's URL.
	DefaultBaseURL = "https://postal.example.com/api/v1"

	// DefaultTimeout is the default timeout for API requests.
	// Requests that take longer than this will be cancelled.
	DefaultTimeout = 30 * time.Second
)

// Client represents a client for the Postal API.
// It handles authentication, request creation, and response parsing.
type Client struct {
	// BaseURL is the base URL for the Postal API.
	// This should include the protocol, domain, and API version path.
	// Example: "https://postal.yourdomain.com/api/v1"
	BaseURL string

	// APIKey is the API key for authenticating with the Postal API.
	// This is required for all requests and should be obtained from
	// your Postal server's administration interface.
	APIKey string

	// HTTPClient is the HTTP client used to make requests.
	// This can be customized to add features like request tracing,
	// custom transport options, or different timeout values.
	HTTPClient *http.Client
}

// NewClient creates a new Postal API client with the given API key.
// It uses default values for the base URL and timeout.
//
// Example:
//
//	client := postalclient.NewClient("your-api-key")
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL:    DefaultBaseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
	}
}

// NewClientWithOptions creates a new Postal API client with the given options.
// This allows for customization of the base URL and timeout.
//
// Example:
//
//	client := postalclient.NewClientWithOptions(
//	    "your-api-key",
//	    "https://postal.yourdomain.com/api/v1",
//	    60 * time.Second,
//	)
func NewClientWithOptions(apiKey, baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

// Response represents a successful response from the Postal API.
// All API responses follow this structure, with the actual data
// contained in the Data field as JSON.
type Response struct {
	// Status indicates the result of the request, typically "success".
	Status string `json:"status"`

	// Time indicates how long the request took to process on the server.
	Time float64 `json:"time"`

	// Flags contains additional metadata about the request.
	// This may include pagination information or other context.
	Flags json.RawMessage `json:"flags"`

	// Data contains the actual response data as JSON.
	// This will be unmarshaled into the appropriate type based on the request.
	Data json.RawMessage `json:"data"`
}

// Error represents an error response from the Postal API.
// It implements the error interface for easy error handling.
type Error struct {
	// Status indicates the result of the request, typically "error" or "parameter-error".
	Status string `json:"status"`

	// Time indicates how long the request took to process on the server.
	Time float64 `json:"time"`

	// Flags contains additional metadata about the request.
	Flags json.RawMessage `json:"flags"`

	// Data contains additional error details as JSON.
	Data json.RawMessage `json:"data"`

	// ErrorCode is a machine-readable code identifying the error.
	ErrorCode string `json:"error_code,omitempty"`

	// Message is a human-readable description of the error.
	Message string `json:"message,omitempty"`
}

// Error returns a string representation of the error.
// This implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("postal API error: %s - %s", e.Status, e.Message)
}

// do performs an HTTP request and returns the response.
// It handles the details of creating the request, setting headers,
// performing the request, and parsing the response.
//
// This is an internal method used by other client methods.
func (c *Client) do(method, path string, body interface{}) (*Response, error) {
	// Create the request URL by combining the base URL and path
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	// Marshal the body to JSON if it's not nil
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create the HTTP request with the specified method, URL, and body
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set required headers for the Postal API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Server-API-Key", c.APIKey)

	// Perform the request using the client's HTTP client
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	// Read the entire response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check if the HTTP status code indicates an error
	if resp.StatusCode != http.StatusOK {
		var apiError Error
		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return nil, fmt.Errorf("error unmarshaling error response: %w", err)
		}
		return nil, &apiError
	}

	// Unmarshal the response body into a Response struct
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check if the API response status indicates an error
	if apiResp.Status != "success" {
		var apiError Error
		if err := json.Unmarshal(respBody, &apiError); err != nil {
			return nil, fmt.Errorf("error unmarshaling error response: %w", err)
		}
		return nil, &apiError
	}

	return &apiResp, nil
}

// post performs a POST request to the given path with the given body.
// It's a convenience wrapper around the do method.
//
// This is an internal method used by other client methods.
func (c *Client) post(path string, body interface{}) (*Response, error) {
	return c.do(http.MethodPost, path, body)
}
