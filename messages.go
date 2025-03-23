// This file contains methods for interacting with the Postal API's message endpoints.
// It provides functionality for retrieving message details, message deliveries,
// and sending messages in both standard and raw formats.
package postalclient

import (
	"encoding/json"
	"fmt"

	"github.com/Suhaibinator/postalclient-go/models"
)

// GetMessage retrieves details about a message with the given ID.
//
// This method calls the /messages/message endpoint to retrieve information
// about a specific message, including its status, content, and metadata.
//
// Example:
//
//	message, err := client.GetMessage(123)
//	if err != nil {
//	    log.Fatalf("Error getting message: %v", err)
//	}
//	fmt.Printf("Message details - ID: %d, Token: %s\n", message.ID, message.Token)
func (c *Client) GetMessage(id int) (*models.Message, error) {
	// Create the request body with the message ID
	body := map[string]interface{}{
		"id": id,
	}

	// Make the request to the API
	resp, err := c.post("/messages/message", body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into a Message struct
	var message models.Message
	if err := json.Unmarshal(resp.Data, &message); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &message, nil
}

// GetMessageDeliveries retrieves delivery attempts for a message with the given ID.
//
// This method calls the /messages/deliveries endpoint to retrieve information
// about all delivery attempts for a specific message, including their status,
// details, and timestamps.
//
// Example:
//
//	deliveries, err := client.GetMessageDeliveries(123)
//	if err != nil {
//	    log.Fatalf("Error getting message deliveries: %v", err)
//	}
//	fmt.Printf("Message has %d deliveries\n", len(deliveries))
//	for i, delivery := range deliveries {
//	    fmt.Printf("Delivery %d - Status: %s\n", i+1, delivery.Status)
//	}
func (c *Client) GetMessageDeliveries(id int) ([]models.Delivery, error) {
	// Create the request body with the message ID
	body := map[string]interface{}{
		"id": id,
	}

	// Make the request to the API
	resp, err := c.post("/messages/deliveries", body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into a slice of Delivery structs
	var deliveries []models.Delivery
	if err := json.Unmarshal(resp.Data, &deliveries); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return deliveries, nil
}

// SendMessage sends an email message using the Postal API.
//
// This method calls the /send/message endpoint to send an email with
// the specified recipients, content, and options. The request must include
// at least one recipient (To, CC, or BCC), a From address, and either
// plain text or HTML content.
//
// Example:
//
//	req := &models.SendMessageRequest{
//	    To:        []string{"recipient@example.com"},
//	    From:      "sender@yourdomain.com",
//	    Subject:   "Hello from Postal API",
//	    PlainBody: "This is a test email sent using the Postal API Go client.",
//	}
//	resp, err := client.SendMessage(req)
//	if err != nil {
//	    log.Fatalf("Error sending message: %v", err)
//	}
//	fmt.Printf("Message sent! ID: %d, Token: %s\n", resp.MessageID, resp.Token)
func (c *Client) SendMessage(req *models.SendMessageRequest) (*models.SendMessageResponse, error) {
	// Make the request to the API
	resp, err := c.post("/send/message", req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into a SendMessageResponse struct
	var sendResp models.SendMessageResponse
	if err := json.Unmarshal(resp.Data, &sendResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling send response: %w", err)
	}

	return &sendResp, nil
}

// SendRaw sends a raw RFC2822 message using the Postal API.
//
// This method calls the /send/raw endpoint to send a pre-formatted email
// message. The request must include the sender address (MailFrom), at least
// one recipient address (RcptTo), and the raw message data encoded in base64.
//
// This is similar to sending a message through SMTP, but using the HTTP API.
//
// Example:
//
//	rawMessage := `From: sender@yourdomain.com
//	To: recipient@example.com
//	Subject: Hello from Postal API
//	Content-Type: text/plain; charset=utf-8
//
//	This is a test email sent using the Postal API Go client.`
//
//	encodedMessage := base64.StdEncoding.EncodeToString([]byte(rawMessage))
//
//	req := &models.SendRawRequest{
//	    MailFrom: "sender@yourdomain.com",
//	    RcptTo:   []string{"recipient@example.com"},
//	    Data:     encodedMessage,
//	}
//
//	resp, err := client.SendRaw(req)
//	if err != nil {
//	    log.Fatalf("Error sending raw message: %v", err)
//	}
//	fmt.Printf("Raw message sent! ID: %d, Token: %s\n", resp.MessageID, resp.Token)
func (c *Client) SendRaw(req *models.SendRawRequest) (*models.SendMessageResponse, error) {
	// Make the request to the API
	resp, err := c.post("/send/raw", req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response data into a SendMessageResponse struct
	var sendResp models.SendMessageResponse
	if err := json.Unmarshal(resp.Data, &sendResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling send response: %w", err)
	}

	return &sendResp, nil
}
