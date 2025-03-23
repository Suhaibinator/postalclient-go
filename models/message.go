// Package models provides data structures for the Postal API.
//
// This package defines the request and response structures used by the
// postalclient package to interact with the Postal API. These structures
// are designed to match the JSON format expected by the API.
package models

import (
	"encoding/json"
	"time"
)

// Message represents a message in the Postal API.
// It contains all the details about an email message, including its
// content, headers, and metadata.
type Message struct {
	// ID is the unique identifier for the message.
	ID int `json:"id"`

	// Token is a unique token that can be used to reference the message.
	Token string `json:"token"`

	// Status contains information about the current status of the message.
	// This is only included when the 'status' expansion is requested.
	Status *MessageStatus `json:"status,omitempty"`

	// Details contains additional details about the message.
	// This is only included when the 'details' expansion is requested.
	Details *MessageDetails `json:"details,omitempty"`

	// Inspection contains information about spam checks and other inspections.
	// This is only included when the 'inspection' expansion is requested.
	Inspection json.RawMessage `json:"inspection,omitempty"`

	// PlainBody is the plain text body of the message.
	// This is only included when the 'plain_body' expansion is requested.
	PlainBody string `json:"plain_body,omitempty"`

	// HTMLBody is the HTML body of the message.
	// This is only included when the 'html_body' expansion is requested.
	HTMLBody string `json:"html_body,omitempty"`

	// Attachments is a list of attachments included with the message.
	// This is only included when the 'attachments' expansion is requested.
	Attachments []Attachment `json:"attachments,omitempty"`

	// Headers is a map of headers included with the message.
	// This is only included when the 'headers' expansion is requested.
	Headers map[string]string `json:"headers,omitempty"`

	// RawMessage is the raw RFC2822 message.
	// This is only included when the 'raw_message' expansion is requested.
	RawMessage string `json:"raw_message,omitempty"`
}

// MessageStatus represents the status of a message.
// This structure is populated when the 'status' expansion is requested.
type MessageStatus struct {
	// Add status fields as needed based on API responses.
	// Common fields might include:
	// - Held (bool)
	// - HoldExpiry (time.Time)
	// - LastDeliveryAttempt (time.Time)
	// - DeliveryStatus (string)
}

// MessageDetails represents the details of a message.
// This structure is populated when the 'details' expansion is requested.
type MessageDetails struct {
	// Add details fields as needed based on API responses.
	// Common fields might include:
	// - SpamScore (float64)
	// - Size (int)
	// - Bounce (bool)
	// - Received (time.Time)
}

// Attachment represents an attachment in a message.
// This structure is populated when the 'attachments' expansion is requested.
type Attachment struct {
	// Name is the filename of the attachment.
	Name string `json:"name"`

	// ContentType is the MIME type of the attachment.
	ContentType string `json:"content_type"`

	// Data is the base64-encoded content of the attachment.
	Data string `json:"data"` // Base64 encoded

	// Size is the size of the attachment in bytes.
	Size int `json:"size"`
}

// Delivery represents a delivery attempt for a message.
// Each message may have multiple delivery attempts, one for each recipient.
type Delivery struct {
	// ID is the unique identifier for the delivery attempt.
	ID int `json:"id"`

	// Status is the current status of the delivery attempt.
	// Common values include "delivered", "failed", "pending", etc.
	Status string `json:"status"`

	// Details provides additional information about the delivery attempt.
	// This might include error messages or delivery confirmations.
	Details string `json:"details"`

	// Output contains the raw output from the receiving mail server.
	// This is useful for debugging delivery issues.
	Output string `json:"output"`

	// SentWithSSL indicates whether the delivery was made using SSL/TLS.
	SentWithSSL bool `json:"sent_with_ssl"`

	// LogID is the ID of the associated log entry.
	LogID int `json:"log_id"`

	// Time is the time taken for the delivery attempt in seconds.
	Time float64 `json:"time"`

	// Timestamp is when the delivery attempt was made.
	Timestamp time.Time `json:"timestamp"`
}
