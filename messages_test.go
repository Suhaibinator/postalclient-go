package postalclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	message, err := client.GetMessage(context.Background(), 123)
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
			"data": {"message": "Message not found", "id": 123, "code": "MessageNotFound"}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get message
	_, err := client.GetMessage(context.Background(), 123)
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
	_, err := client.GetMessage(context.Background(), 123)
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
					"log_id": "A1B2C3D4",
					"time": 0.46,
					"timestamp": 1672531200.123456
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get deliveries
	deliveries, err := client.GetMessageDeliveries(context.Background(), 123)
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

	if delivery.LogID != "A1B2C3D4" {
		t.Errorf("Expected delivery LogID to be A1B2C3D4, got %s", delivery.LogID)
	}

	if delivery.Time != 0.46 {
		t.Errorf("Expected delivery Time to be 0.46, got %v", delivery.Time)
	}

	wantTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 123456000, time.UTC)
	if !delivery.Timestamp.Equal(wantTimestamp) {
		t.Errorf("Expected delivery Timestamp to be %s, got %s", wantTimestamp, delivery.Timestamp)
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
			"data": {"message": "Message not found", "id": 123, "code": "MessageNotFound"}
		}`))
	}))
	defer server.Close()

	// Create client
	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	// Get deliveries
	_, err := client.GetMessageDeliveries(context.Background(), 123)
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
	_, err := client.GetMessageDeliveries(context.Background(), 123)
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
				"message_id": "8b5f9c3e-1d2a-4c5b-9e7f-0a1b2c3d4e5f@psrp.example.com",
				"messages": {
					"recipient@example.com": {"id": 123, "token": "test-token"}
				}
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
	resp, err := client.SendMessage(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	checkSendResponse(t, resp)
}

func TestSendMessageError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status": "error",
			"time": 0.123,
			"flags": {},
			"data": {"message": "The From address is missing and is required", "code": "FromAddressMissing"}
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
	_, err := client.SendMessage(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *Error, got %T", err)
	}
	if apiErr.Status != "error" || apiErr.ErrorCode != "FromAddressMissing" || apiErr.Message != "The From address is missing and is required" {
		t.Errorf("Unexpected error fields: %+v", apiErr)
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
	_, err := client.SendMessage(context.Background(), req)
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
				"message_id": "8b5f9c3e-1d2a-4c5b-9e7f-0a1b2c3d4e5f@psrp.example.com",
				"messages": {
					"recipient@example.com": {"id": 123, "token": "test-token"}
				}
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
	resp, err := client.SendRaw(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	checkSendResponse(t, resp)
}

func TestSendRawError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\":\"parameter-error\",\"time\":0.123,\"flags\":{}," +
			"\"data\":{\"message\":\"`mail_from` parameter is required but is missing\"}}"))
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
	_, err := client.SendRaw(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *Error, got %T", err)
	}
	if apiErr.Status != "parameter-error" || apiErr.ErrorCode != "" || apiErr.Message != "`mail_from` parameter is required but is missing" {
		t.Errorf("Unexpected error fields: %+v", apiErr)
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
	_, err := client.SendRaw(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error
	if !strings.Contains(err.Error(), "error unmarshaling response") {
		t.Errorf("Expected error message to contain 'error unmarshaling response', got '%s'", err.Error())
	}
}

// checkSendResponse checks the decoded shape Postal returns from both send endpoints.
func checkSendResponse(t *testing.T, resp *models.SendMessageResponse) {
	t.Helper()
	if resp.MessageID != "8b5f9c3e-1d2a-4c5b-9e7f-0a1b2c3d4e5f@psrp.example.com" {
		t.Errorf("Unexpected message ID %q", resp.MessageID)
	}
	m, ok := resp.Messages["recipient@example.com"]
	if !ok {
		t.Fatalf("Expected a message for recipient@example.com, got %v", resp.Messages)
	}
	if m.ID != 123 || m.Token != "test-token" {
		t.Errorf("Expected message {123 test-token}, got %+v", m)
	}
}

func TestGetMessageExpansions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		if string(body["id"]) != "42" {
			t.Errorf("Expected id 42, got %s", body["id"])
		}
		if got := string(body["_expansions"]); got != `["status","details","headers","attachments","activity_entries"]` {
			t.Errorf("Unexpected _expansions %s", got)
		}
		// Shape rendered by Postal's LegacyAPI::MessagesController#message
		_, _ = w.Write([]byte(`{"status":"success","time":0.01,"flags":{},"data":{
			"id": 42,
			"token": "AbCdEf123456",
			"status": {"status": "Held", "last_delivery_attempt": null, "held": true, "hold_expiry": 1700000000.5},
			"details": {
				"rcpt_to": "recipient@example.com",
				"mail_from": "sender@example.com",
				"subject": "Hi",
				"message_id": "abc@example.com",
				"timestamp": 1672531200.25,
				"direction": "outgoing",
				"size": "2048",
				"bounce": false,
				"bounce_for_id": 0,
				"tag": null,
				"received_with_ssl": true
			},
			"headers": {"received": ["from a", "from b"], "subject": ["Hi"]},
			"attachments": [{"filename": "a.txt", "content_type": "text/plain", "data": "aGVs\nbG8=\n", "size": 5, "hash": "aaf4c61d"}],
			"activity_entries": {"loads": [], "clicks": []}
		}}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.BaseURL = server.URL

	msg, err := client.GetMessage(context.Background(), 42,
		models.ExpansionStatus, models.ExpansionDetails, models.ExpansionHeaders,
		models.ExpansionAttachments, models.ExpansionActivityEntries)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if msg.Status == nil || msg.Status.Status != "Held" || !msg.Status.Held || msg.Status.LastDeliveryAttempt != nil {
		t.Errorf("Unexpected status %+v", msg.Status)
	}
	if msg.Status.HoldExpiry == nil || !msg.Status.HoldExpiry.Equal(time.Unix(1700000000, 500000000)) {
		t.Errorf("Unexpected hold expiry %v", msg.Status.HoldExpiry)
	}
	d := msg.Details
	if d == nil || d.RcptTo != "recipient@example.com" || d.MessageID != "abc@example.com" || d.Size != 2048 ||
		d.Direction != "outgoing" || !d.ReceivedWithSSL || !d.Timestamp.Equal(time.Unix(1672531200, 250000000)) {
		t.Errorf("Unexpected details %+v", d)
	}
	if got := msg.Headers["received"]; len(got) != 2 || got[1] != "from b" {
		t.Errorf("Unexpected headers %v", msg.Headers)
	}
	if len(msg.Attachments) != 1 || msg.Attachments[0].Filename != "a.txt" || msg.Attachments[0].Hash != "aaf4c61d" || msg.Attachments[0].Size != 5 {
		t.Errorf("Unexpected attachments %+v", msg.Attachments)
	}
	if string(msg.ActivityEntries) != `{"loads": [], "clicks": []}` {
		t.Errorf("Unexpected activity entries %s", msg.ActivityEntries)
	}
}

func TestGetMessageWithoutExpansionsOmitsParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		_ = json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["_expansions"]; ok {
			t.Errorf("Expected no _expansions, got %s", body["_expansions"])
		}
		_, _ = w.Write([]byte(`{"status":"success","time":0.01,"flags":{},"data":{"id":42,"token":"t"}}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.BaseURL = server.URL
	if _, err := client.GetMessage(context.Background(), 42); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetMessageDeliveriesNullFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"success","time":0.01,"flags":{},"data":[
			{"id": 1, "status": "SoftFail", "details": null, "output": null, "sent_with_ssl": false, "log_id": null, "time": null, "timestamp": 1672531200.0}
		]}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.BaseURL = server.URL
	deliveries, err := client.GetMessageDeliveries(context.Background(), 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(deliveries) != 1 || deliveries[0].LogID != "" || !deliveries[0].Timestamp.Equal(time.Unix(1672531200, 0)) {
		t.Errorf("Unexpected deliveries %+v", deliveries)
	}
}

func TestSendRawEmptyMessageID(t *testing.T) {
	// Raw messages without a Message-ID header yield a null message_id
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]json.RawMessage
		_ = json.NewDecoder(r.Body).Decode(&body)
		for _, key := range []string{"mail_from", "rcpt_to", "data"} {
			if _, ok := body[key]; !ok {
				t.Errorf("Expected %s in request body", key)
			}
		}
		_, _ = w.Write([]byte(`{"status":"success","time":0.01,"flags":{},"data":{"message_id":null,"messages":{"a@example.com":{"id":7,"token":"tok"}}}}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.BaseURL = server.URL
	resp, err := client.SendRaw(context.Background(), &models.SendRawRequest{MailFrom: "s@example.com", RcptTo: []string{"a@example.com"}, Data: "eA=="})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.MessageID != "" || resp.Messages["a@example.com"].ID != 7 {
		t.Errorf("Unexpected response %+v", resp)
	}
}

func TestSendMessageRequestShape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Headers     map[string]string            `json:"headers"`
			Attachments []map[string]json.RawMessage `json:"attachments"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("Error decoding request body: %v", err)
		}
		if body.Headers["X-Test"] != "1" {
			t.Errorf("Expected headers to carry X-Test, got %v", body.Headers)
		}
		if len(body.Attachments) != 1 {
			t.Fatalf("Expected one attachment, got %v", body.Attachments)
		}
		a := body.Attachments[0]
		if string(a["name"]) != `"a.txt"` || string(a["data"]) != `"aGVsbG8="` {
			t.Errorf("Unexpected attachment %v", a)
		}
		for _, key := range []string{"size", "content_type"} {
			if _, ok := a[key]; ok {
				t.Errorf("Expected no %s in attachment, got %s", key, a[key])
			}
		}
		_, _ = w.Write([]byte(`{"status":"success","time":0.01,"flags":{},"data":{"message_id":"x@example.com","messages":{}}}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.BaseURL = server.URL
	_, err := client.SendMessage(context.Background(), &models.SendMessageRequest{
		To:          []string{"a@example.com"},
		From:        "s@example.com",
		PlainBody:   "hi",
		Headers:     map[string]string{"X-Test": "1"},
		Attachments: []models.Attachment{{Name: "a.txt", Data: "aGVsbG8="}},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
