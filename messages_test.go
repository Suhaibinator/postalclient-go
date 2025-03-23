package postalclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Suhaibinator/postalclient-go/models"
)

func TestGetMessage(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/messages/message" {
			t.Errorf("Expected request path to be /messages/message, got %s", r.URL.Path)
		}

		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be %s, got %s", http.MethodPost, r.Method)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": {
				"id": 123,
				"token": "test-token"
			}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get message
	message, err := client.GetMessage(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check message
	if message.ID != 123 {
		t.Errorf("Expected message ID to be 123, got %d", message.ID)
	}

	if message.Token != "test-token" {
		t.Errorf("Expected message token to be test-token, got %s", message.Token)
	}
}

func TestGetMessageError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "error",
			"time": 0.123,
			"flags": {},
			"data": {"error": "test error"},
			"error_code": "message-not-found",
			"message": "Message not found"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get message
	_, err := client.GetMessage(123)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "Message not found") {
		t.Errorf("Expected error message to contain 'Message not found', got '%s'", err.Error())
	}
}

func TestGetMessageUnmarshalError(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": invalid-json
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get message
	_, err := client.GetMessage(123)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

func TestGetMessageDeliveries(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/messages/deliveries" {
			t.Errorf("Expected request path to be /messages/deliveries, got %s", r.URL.Path)
		}

		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be %s, got %s", http.MethodPost, r.Method)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": [
				{
					"id": 123,
					"status": "delivered",
					"details": "test details",
					"output": "test output",
					"sent_with_ssl": true,
					"log_id": 456,
					"time": 0.456,
					"timestamp": "2023-01-01T00:00:00Z"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get deliveries
	deliveries, err := client.GetMessageDeliveries(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check deliveries
	if len(deliveries) != 1 {
		t.Fatalf("Expected 1 delivery, got %d", len(deliveries))
	}

	delivery := deliveries[0]
	if delivery.ID != 123 {
		t.Errorf("Expected delivery ID to be 123, got %d", delivery.ID)
	}

	if delivery.Status != "delivered" {
		t.Errorf("Expected delivery status to be delivered, got %s", delivery.Status)
	}

	if delivery.SentWithSSL != true {
		t.Errorf("Expected delivery SentWithSSL to be true, got %v", delivery.SentWithSSL)
	}
}

func TestGetMessageDeliveriesError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "error",
			"time": 0.123,
			"flags": {},
			"data": {"error": "test error"},
			"error_code": "message-not-found",
			"message": "Message not found"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get deliveries
	_, err := client.GetMessageDeliveries(123)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "Message not found") {
		t.Errorf("Expected error message to contain 'Message not found', got '%s'", err.Error())
	}
}

func TestGetMessageDeliveriesUnmarshalError(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": invalid-json
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get deliveries
	_, err := client.GetMessageDeliveries(123)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

func TestSendMessage(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/send/message" {
			t.Errorf("Expected request path to be /send/message, got %s", r.URL.Path)
		}

		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be %s, got %s", http.MethodPost, r.Method)
		}

		// Check request body
		var req models.SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}

		if req.From != "test@example.com" {
			t.Errorf("Expected From to be test@example.com, got %s", req.From)
		}

		if len(req.To) != 1 || req.To[0] != "recipient@example.com" {
			t.Errorf("Expected To to be [recipient@example.com], got %v", req.To)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": {
				"message_id": 123,
				"token": "test-token"
			}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendMessageRequest{
		From:      "test@example.com",
		To:        []string{"recipient@example.com"},
		Subject:   "Test Subject",
		PlainBody: "Test Body",
	}

	// Send message
	resp, err := client.SendMessage(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	if resp.MessageID != 123 {
		t.Errorf("Expected message ID to be 123, got %d", resp.MessageID)
	}

	if resp.Token != "test-token" {
		t.Errorf("Expected token to be test-token, got %s", resp.Token)
	}
}

func TestSendMessageError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "parameter-error",
			"time": 0.123,
			"flags": {},
			"data": {"errors": {"from": ["is required"]}},
			"error_code": "validation-error",
			"message": "The provided data was not sufficient to send an email"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendMessageRequest{
		To:        []string{"recipient@example.com"},
		Subject:   "Test Subject",
		PlainBody: "Test Body",
	}

	// Send message
	_, err := client.SendMessage(req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "The provided data was not sufficient to send an email") {
		t.Errorf("Expected error message to contain 'The provided data was not sufficient to send an email', got '%s'", err.Error())
	}
}

func TestSendMessageUnmarshalError(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": invalid-json
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendMessageRequest{
		From:      "test@example.com",
		To:        []string{"recipient@example.com"},
		Subject:   "Test Subject",
		PlainBody: "Test Body",
	}

	// Send message
	_, err := client.SendMessage(req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

func TestSendRaw(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/send/raw" {
			t.Errorf("Expected request path to be /send/raw, got %s", r.URL.Path)
		}

		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be %s, got %s", http.MethodPost, r.Method)
		}

		// Check request body
		var req models.SendRawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}

		if req.MailFrom != "test@example.com" {
			t.Errorf("Expected MailFrom to be test@example.com, got %s", req.MailFrom)
		}

		if len(req.RcptTo) != 1 || req.RcptTo[0] != "recipient@example.com" {
			t.Errorf("Expected RcptTo to be [recipient@example.com], got %v", req.RcptTo)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": {
				"message_id": 123,
				"token": "test-token"
			}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendRawRequest{
		MailFrom: "test@example.com",
		RcptTo:   []string{"recipient@example.com"},
		Data:     "base64-encoded-data",
	}

	// Send raw message
	resp, err := client.SendRaw(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	if resp.MessageID != 123 {
		t.Errorf("Expected message ID to be 123, got %d", resp.MessageID)
	}

	if resp.Token != "test-token" {
		t.Errorf("Expected token to be test-token, got %s", resp.Token)
	}
}

func TestSendRawError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "parameter-error",
			"time": 0.123,
			"flags": {},
			"data": {"errors": {"mail_from": ["is required"]}},
			"error_code": "validation-error",
			"message": "The provided data was not sufficient to send an email"
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendRawRequest{
		RcptTo: []string{"recipient@example.com"},
		Data:   "base64-encoded-data",
	}

	// Send raw message
	_, err := client.SendRaw(req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "The provided data was not sufficient to send an email") {
		t.Errorf("Expected error message to contain 'The provided data was not sufficient to send an email', got '%s'", err.Error())
	}
}

func TestSendRawUnmarshalError(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "success",
			"time": 0.123,
			"flags": {},
			"data": invalid-json
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Create request
	req := &models.SendRawRequest{
		MailFrom: "test@example.com",
		RcptTo:   []string{"recipient@example.com"},
		Data:     "base64-encoded-data",
	}

	// Send raw message
	_, err := client.SendRaw(req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}
