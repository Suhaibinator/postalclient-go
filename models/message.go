// Package models provides data structures for the Postal API.
//
// This package defines the request and response structures used by the
// postalclient package to interact with the Postal API. These structures
// are designed to match the JSON format expected by the API.
package models

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

// Expansions that GetMessage can request. Each one populates the
// matching field of Message.
const (
	ExpansionStatus          = "status"
	ExpansionDetails         = "details"
	ExpansionInspection      = "inspection"
	ExpansionPlainBody       = "plain_body"
	ExpansionHTMLBody        = "html_body"
	ExpansionAttachments     = "attachments"
	ExpansionHeaders         = "headers"
	ExpansionRawMessage      = "raw_message"
	ExpansionActivityEntries = "activity_entries"
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
	Attachments []MessageAttachment `json:"attachments,omitempty"`

	// Headers maps each lower-cased header name to its values, in order.
	// This is only included when the 'headers' expansion is requested.
	Headers map[string][]string `json:"headers,omitempty"`

	// RawMessage is the base64-encoded raw RFC2822 message.
	// This is only included when the 'raw_message' expansion is requested.
	RawMessage string `json:"raw_message,omitempty"`

	// ActivityEntries holds the message's open ("loads") and click records.
	// This is only included when the 'activity_entries' expansion is requested.
	ActivityEntries json.RawMessage `json:"activity_entries,omitempty"`
}

// MessageStatus represents the status of a message.
// This structure is populated when the 'status' expansion is requested.
type MessageStatus struct {
	// Status is the message status, such as "Pending", "Sent", "Held",
	// "SoftFail", "HardFail" or "Bounced".
	Status string `json:"status"`

	// LastDeliveryAttempt is when delivery was last attempted, or nil.
	LastDeliveryAttempt *time.Time `json:"last_delivery_attempt"`

	// Held reports whether the message is held.
	Held bool `json:"held"`

	// HoldExpiry is when a held message is released, or nil.
	HoldExpiry *time.Time `json:"hold_expiry"`
}

// UnmarshalJSON decodes a status whose times are Unix timestamps in seconds.
func (s *MessageStatus) UnmarshalJSON(b []byte) error {
	type plain MessageStatus
	aux := struct {
		*plain
		LastDeliveryAttempt *float64 `json:"last_delivery_attempt"`
		HoldExpiry          *float64 `json:"hold_expiry"`
	}{plain: (*plain)(s)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	s.LastDeliveryAttempt = optionalUnixTime(aux.LastDeliveryAttempt)
	s.HoldExpiry = optionalUnixTime(aux.HoldExpiry)
	return nil
}

// MessageDetails represents the details of a message.
// This structure is populated when the 'details' expansion is requested.
type MessageDetails struct {
	// RcptTo is the recipient address.
	RcptTo string `json:"rcpt_to"`

	// MailFrom is the envelope sender address.
	MailFrom string `json:"mail_from"`

	// Subject is the message subject.
	Subject string `json:"subject"`

	// MessageID is the RFC 5322 Message-ID, without angle brackets.
	MessageID string `json:"message_id"`

	// Timestamp is when the message was created.
	Timestamp time.Time `json:"timestamp"`

	// Direction is "incoming" or "outgoing".
	Direction string `json:"direction"`

	// Size is the size of the raw message in bytes.
	Size int64 `json:"size"`

	// Bounce reports whether the message is a bounce.
	Bounce bool `json:"bounce"`

	// BounceForID is the ID of the message this one bounces, or 0.
	BounceForID int `json:"bounce_for_id"`

	// Tag is the message's tag.
	Tag string `json:"tag"`

	// ReceivedWithSSL reports whether the message was received over TLS.
	ReceivedWithSSL bool `json:"received_with_ssl"`
}

// UnmarshalJSON decodes details whose timestamp is a Unix timestamp in
// seconds and whose size Postal stores, and returns, as a string.
func (d *MessageDetails) UnmarshalJSON(b []byte) error {
	type plain MessageDetails
	aux := struct {
		*plain
		Timestamp *float64        `json:"timestamp"`
		Size      json.RawMessage `json:"size"`
	}{plain: (*plain)(d)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	d.Timestamp = time.Time{}
	if t := optionalUnixTime(aux.Timestamp); t != nil {
		d.Timestamp = *t
	}
	d.Size, _ = strconv.ParseInt(strings.Trim(string(aux.Size), `"`), 10, 64)
	return nil
}

// Attachment is a file to attach to a message sent with SendMessage.
type Attachment struct {
	// Name is the filename of the attachment.
	Name string `json:"name"`

	// ContentType is the MIME type of the attachment.
	// Postal uses "application/octet-stream" when it is empty.
	ContentType string `json:"content_type,omitempty"`

	// Data is the base64-encoded content of the attachment.
	Data string `json:"data"`
}

// MessageAttachment is an attachment of a message returned by GetMessage.
// This structure is populated when the 'attachments' expansion is requested.
type MessageAttachment struct {
	// Filename is the filename of the attachment.
	Filename string `json:"filename"`

	// ContentType is the MIME type of the attachment.
	ContentType string `json:"content_type"`

	// Data is the base64-encoded content of the attachment. Postal wraps it
	// with newlines, which encoding/base64 ignores.
	Data string `json:"data"`

	// Size is the size of the attachment in bytes.
	Size int `json:"size"`

	// Hash is the hex SHA-1 of the attachment's content.
	Hash string `json:"hash"`
}

// Delivery represents a delivery attempt for a message.
// Each message may have multiple delivery attempts, one for each recipient.
type Delivery struct {
	// ID is the unique identifier for the delivery attempt.
	ID int `json:"id"`

	// Status is the status of the delivery attempt, such as "Sent",
	// "SoftFail", "HardFail", "Held" or "Bounced".
	Status string `json:"status"`

	// Details provides additional information about the delivery attempt.
	// This might include error messages or delivery confirmations.
	Details string `json:"details"`

	// Output contains the raw output from the receiving mail server.
	// This is useful for debugging delivery issues.
	Output string `json:"output"`

	// SentWithSSL indicates whether the delivery was made using SSL/TLS.
	SentWithSSL bool `json:"sent_with_ssl"`

	// LogID is the ID of the associated log entry, such as "A1B2C3D4".
	LogID string `json:"log_id"`

	// Time is the time taken for the delivery attempt in seconds.
	Time float64 `json:"time"`

	// Timestamp is when the delivery attempt was made. Postal sends it as a
	// Unix timestamp in seconds; json.Marshal writes it back as RFC 3339.
	Timestamp time.Time `json:"timestamp"`
}

// UnmarshalJSON decodes a delivery whose timestamp is a Unix timestamp in seconds.
func (d *Delivery) UnmarshalJSON(b []byte) error {
	type plain Delivery
	aux := struct {
		*plain
		Timestamp *float64 `json:"timestamp"`
	}{plain: (*plain)(d)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	d.Timestamp = time.Time{}
	if t := optionalUnixTime(aux.Timestamp); t != nil {
		d.Timestamp = *t
	}
	return nil
}

// optionalUnixTime converts Unix seconds with microsecond precision, as
// Postal stores them, to a UTC time. It returns nil for nil.
func optionalUnixTime(seconds *float64) *time.Time {
	if seconds == nil {
		return nil
	}
	t := time.UnixMicro(int64(math.Round(*seconds * 1e6))).UTC()
	return &t
}
