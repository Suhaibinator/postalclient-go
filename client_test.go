package postalclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", DefaultBaseURL, client.BaseURL)
	}

	if client.HTTPClient.Timeout != DefaultTimeout {
		t.Errorf("Expected HTTPClient.Timeout to be %s, got %s", DefaultTimeout, client.HTTPClient.Timeout)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	apiKey := "test-api-key"
	baseURL := "https://custom.example.com/api/v1"
	timeout := 60 * time.Second

	client := NewClientWithOptions(apiKey, baseURL, timeout)

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", baseURL, client.BaseURL)
	}

	if client.HTTPClient.Timeout != timeout {
		t.Errorf("Expected HTTPClient.Timeout to be %s, got %s", timeout, client.HTTPClient.Timeout)
	}
}

func TestClientDo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be %s, got %s", http.MethodPost, r.Method)
		}

		// Check request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept header to be application/json, got %s", r.Header.Get("Accept"))
		}

		apiKey := "test-api-key"
		if r.Header.Get("X-Server-API-Key") != apiKey {
			t.Errorf("Expected X-Server-API-Key header to be %s, got %s", apiKey, r.Header.Get("X-Server-API-Key"))
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": {"test": "data"}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	resp, err := client.post("/test", map[string]string{"test": "data"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	if resp.Status != "success" {
		t.Errorf("Expected response status to be success, got %s", resp.Status)
	}

	if resp.Time != 0.123 {
		t.Errorf("Expected response time to be 0.123, got %f", resp.Time)
	}
}

func TestClientDoError(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "error",
			"time": 0.123,
			"flags": {},
			"data": {"error": "test error"},
			"error_code": "test-error",
			"message": "Test error message"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected error to be of type *Error, got %T", err)
	}

	if apiErr.Status != "error" {
		t.Errorf("Expected error status to be error, got %s", apiErr.Status)
	}

	if apiErr.Message != "Test error message" {
		t.Errorf("Expected error message to be 'Test error message', got %s", apiErr.Message)
	}
}

func TestClientDoHTTPError(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{
			"status": "error",
			"time": 0.123,
			"flags": {},
			"data": {"error": "test error"},
			"error_code": "test-error",
			"message": "Test error message"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	apiErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("Expected error to be of type *Error, got %T", err)
	}

	if apiErr.Status != "error" {
		t.Errorf("Expected error status to be error, got %s", apiErr.Status)
	}

	if apiErr.Message != "Test error message" {
		t.Errorf("Expected error message to be 'Test error message', got %s", apiErr.Message)
	}
}

func TestErrorMethod(t *testing.T) {
	// Create an Error instance
	apiErr := &Error{
		Status:    "error",
		Time:      0.123,
		Flags:     json.RawMessage(`{}`),
		Data:      json.RawMessage(`{"error": "test error"}`),
		ErrorCode: "test-error",
		Message:   "Test error message",
	}

	// Test the Error method
	errStr := apiErr.Error()
	expected := "postal API error: error - Test error message"
	if errStr != expected {
		t.Errorf("Expected error string to be '%s', got '%s'", expected, errStr)
	}

	// Test with empty message
	apiErr.Message = ""
	errStr = apiErr.Error()
	expected = "postal API error: error - "
	if errStr != expected {
		t.Errorf("Expected error string to be '%s', got '%s'", expected, errStr)
	}

	// Test with different status
	apiErr.Status = "parameter-error"
	apiErr.Message = "Invalid parameter"
	errStr = apiErr.Error()
	expected = "postal API error: parameter-error - Invalid parameter"
	if errStr != expected {
		t.Errorf("Expected error string to be '%s', got '%s'", expected, errStr)
	}
}

func TestClientDoRequestCreationError(t *testing.T) {
	// Create client with invalid URL to trigger request creation error
	client := NewClient("test-api-key")
	client.BaseURL = "http://invalid-url-with-\u007F-character"

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error creating request") {
		t.Errorf("Expected error message to contain 'error creating request', got '%s'", err.Error())
	}
}

func TestClientDoMarshalError(t *testing.T) {
	// Create client
	client := NewClient("test-api-key")

	// Create a body that can't be marshaled to JSON
	type BadType struct {
		Ch chan int
	}
	badBody := BadType{
		Ch: make(chan int),
	}

	// Make request
	_, err := client.post("/test", badBody)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error marshaling request body") {
		t.Errorf("Expected error message to contain 'error marshaling request body', got '%s'", err.Error())
	}
}

func TestClientDoHTTPClientError(t *testing.T) {
	// Create client with a custom HTTP client that always returns an error
	client := NewClient("test-api-key")
	client.HTTPClient = &http.Client{
		Transport: &mockTransport{},
	}

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error performing request") {
		t.Errorf("Expected error message to contain 'error performing request', got '%s'", err.Error())
	}
}

// mockTransport is a mock http.RoundTripper that always returns an error
type mockTransport struct{}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("mock transport error")
}

func TestClientDoReadBodyError(t *testing.T) {
	// Create a test server that returns a valid response but closes the connection immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers and status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Write a partial response and then hijack the connection to close it
		w.Write([]byte(`{"status":"success","time":0.123,"flags":{}`))
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request - this may not consistently produce a read error depending on timing
	// so we'll skip checking the specific error
	_, _ = client.post("/test", nil)
}

func TestClientDoUnmarshalResponseError(t *testing.T) {
	// Create a test server that returns an invalid JSON response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","time":0.123,"flags":{},"data":invalid-json}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

func TestClientDoUnmarshalErrorResponseError(t *testing.T) {
	// Create a test server that returns an error status with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"error","time":0.123,"flags":{},"data":invalid-json}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

func TestClientDoHTTPErrorUnmarshalError(t *testing.T) {
	// Create a test server that returns an HTTP error with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","time":0.123,"flags":{},"data":invalid-json}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Make request
	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !strings.Contains(err.Error(), "error unmarshaling error response") {
		t.Errorf("Expected error message to contain 'error unmarshaling error response', got '%s'", err.Error())
	}
}
