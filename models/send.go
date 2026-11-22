// Package models provides data structures for the Postal API.
//
// This file contains the request and response structures for sending messages
// through the Postal API, both in standard format and as raw RFC2822 messages.
package models

// SendMessageRequest represents a request to send a message through the Postal API.
// This structure is used with the SendMessage method to create and send an email.
type SendMessageRequest struct {
	// To is a list of recipient email addresses in the To field.
	// At least one recipient (To, CC, or BCC) is required.
	// Maximum 50 recipients.
	To []string `json:"to"`

	// CC is a list of carbon copy recipient email addresses.
	// Optional. Maximum 50 recipients.
	CC []string `json:"cc,omitempty"`

	// BCC is a list of blind carbon copy recipient email addresses.
	// Optional. Maximum 50 recipients.
	BCC []string `json:"bcc,omitempty"`

	// From is the sender's email address that appears in the From header.
	// This is required and must be authorized for the server.
	From string `json:"from"`

	// Sender is the email address for the Sender header.
	// Optional. If not provided, the From address is used.
	Sender string `json:"sender,omitempty"`

	// Subject is the subject line of the email.
	// This is required.
	Subject string `json:"subject"`

	// Tag is an optional tag for categorizing the email.
	// This can be used for filtering and analytics.
	Tag string `json:"tag,omitempty"`

	// ReplyTo is the email address that should receive replies.
	// Optional. If not provided, replies will go to the From address.
	ReplyTo string `json:"reply_to,omitempty"`

	// PlainBody is the plain text version of the email body.
	// Either PlainBody or HTMLBody (or both) must be provided.
	PlainBody string `json:"plain_body,omitempty"`

	// HTMLBody is the HTML version of the email body.
	// Either PlainBody or HTMLBody (or both) must be provided.
	HTMLBody string `json:"html_body,omitempty"`

	// Attachments is a list of files to attach to the email.
	// Optional.
	Attachments []Attachment `json:"attachments,omitempty"`

	// Headers is a map of additional headers to include in the email.
	// Optional.
	Headers map[string]string `json:"headers,omitempty"`

	// Bounce indicates whether this message is a bounce.
	// Optional. Default is false.
	Bounce bool `json:"bounce,omitempty"`
}

// SendRawRequest represents a request to send a raw RFC2822 message.
// This structure is used with the SendRaw method to send a pre-formatted email.
type SendRawRequest struct {
	// MailFrom is the address that should be logged as sending the message.
	// This is required and must be authorized for the server.
	MailFrom string `json:"mail_from"`

	// RcptTo is the list of recipient addresses for the message.
	// This is required.
	RcptTo []string `json:"rcpt_to"`

	// Data is the base64-encoded RFC2822 message to send.
	// This is required and must be a valid RFC2822 message.
	Data string `json:"data"` // Base64 encoded RFC2822 message

	// Bounce indicates whether this message is a bounce.
	// Optional. Default is false.
	Bounce bool `json:"bounce,omitempty"`
}

// SendMessageResponse represents a response from sending a message.
// This structure is returned by both SendMessage and SendRaw methods.
type SendMessageResponse struct {
	// MessageID is the RFC 5322 Message-ID of the sent message, without angle
	// brackets, e.g. "3f2c...-...@psrp.example.com". For SendRaw it is taken
	// from the raw message's Message-ID header and is empty if it has none.
	// It is not the numeric ID that GetMessage takes; see Messages for that.
	MessageID string `json:"message_id"`

	// Messages maps each recipient address to the message Postal created for
	// it. Postal stores one message per recipient, and its numeric ID is the
	// one GetMessage and GetMessageDeliveries take.
	Messages map[string]RecipientMessage `json:"messages"`
}

// RecipientMessage identifies the message Postal created for one recipient.
type RecipientMessage struct {
	// ID is the numeric message ID, for use with GetMessage and GetMessageDeliveries.
	ID int `json:"id"`

	// Token is the message's unique token.
	Token string `json:"token"`
}
