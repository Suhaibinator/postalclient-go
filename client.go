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
//	resp, err := client.SendMessage(ctx, req)
//
// Every API method takes a context.Context for cancellation and deadlines.
// HTTPClient's timeout still applies as well.
package postalclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	// A trailing slash is ignored.
	// Example: "https://postal.yourdomain.com/api/v1"
	BaseURL string

	// APIKey is the API key for authenticating with the Postal API.
	// This is required for all requests and should be obtained from
	// your Postal server's administration interface.
	APIKey string

	// HTTPClient is the HTTP client used to make requests.
	// This can be customized to add features like request tracing,
	// custom transport options, or different timeout values.
	// If nil, http.DefaultClient is used.
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
//
// The Postal legacy API answers HTTP 200 for failed requests and reports the
// failure in the JSON body: "status" is "error" (with "data.code" and,
// usually, "data.message") or "parameter-error" (with "data.message").
// ErrorCode and Message are copied from those fields.
type Error struct {
	// Status indicates the result of the request, typically "error" or "parameter-error".
	Status string `json:"status"`

	// Time indicates how long the request took to process on the server.
	Time float64 `json:"time"`

	// Flags contains additional metadata about the request.
	Flags json.RawMessage `json:"flags"`

	// Data contains the raw error details as JSON. For "InvalidServerAPIKey"
	// Postal echoes the rejected key in data.token, so avoid logging Data.
	Data json.RawMessage `json:"data"`

	// ErrorCode is the machine-readable code from data.code, such as
	// "InvalidServerAPIKey", "UnauthenticatedFromAddress" or "MessageNotFound".
	// It is empty for parameter errors.
	ErrorCode string `json:"-"`

	// Message is the human-readable description from data.message.
	// It may be empty; some errors (e.g. "ServerSuspended") carry only a code.
	Message string `json:"-"`
}

// Error returns a string representation of the error.
// This implements the error interface.
func (e *Error) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("postal API error: %s (%s) - %s", e.Status, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("postal API error: %s - %s", e.Status, e.Message)
}

// newError builds an Error from a non-success response envelope, copying
// data.code and data.message into ErrorCode and Message.
func newError(r *Response) *Error {
	apiErr := &Error{Status: r.Status, Time: r.Time, Flags: r.Flags, Data: r.Data}
	var details struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if json.Unmarshal(r.Data, &details) == nil {
		apiErr.ErrorCode = details.Code
		apiErr.Message = details.Message
	}
	return apiErr
}

// HTTPError is returned when the server answers with a non-200 HTTP status
// and a body that is not a Postal API response. The Postal legacy API always
// answers 200, so this usually means a proxy error page or a 404 from a
// BaseURL that does not point at the API (it must end in "/api/v1").
type HTTPError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int

	// Body holds up to the first 64 KiB of the response body.
	Body string
}

// Error returns a string representation of the error.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("postal API error: unexpected HTTP status %d %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// maxErrorBodySize bounds how much of a non-200 response body is read.
const maxErrorBodySize = 64 << 10

// do performs an HTTP request and returns the response.
// It handles the details of creating the request, setting headers,
// performing the request, and parsing the response.
//
// This is an internal method used by other client methods.
func (c *Client) do(ctx context.Context, method, path string, body interface{}) (*Response, error) {
	// Create the request URL by combining the base URL and path
	url := strings.TrimRight(c.BaseURL, "/") + path

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
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set required headers for the Postal API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Server-API-Key", c.APIKey)

	// Perform the request using the client's HTTP client
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	// The Postal API always answers 200; anything else comes from somewhere
	// else (a proxy, or a BaseURL that does not point at the API).
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
		var apiResp Response
		if json.Unmarshal(respBody, &apiResp) == nil && apiResp.Status != "" {
			return nil, newError(&apiResp)
		}
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	// Read the entire response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Unmarshal the response body into a Response struct
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// The API reports failures in the status field, not the HTTP status code.
	if apiResp.Status != "success" {
		return nil, newError(&apiResp)
	}

	return &apiResp, nil
}

// post performs a POST request to the given path with the given body.
// It's a convenience wrapper around the do method.
//
// This is an internal method used by other client methods.
func (c *Client) post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPost, path, body)
}
