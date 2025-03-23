package message

import (
	"fmt"
	"log"

	"github.com/Suhaibinator/postalclient-go"
	"github.com/Suhaibinator/postalclient-go/models"
)

// SendExample demonstrates how to send a message using the Postal API
func SendExample(apiKey, baseURL string) {
	// Create a new client
	client := postalclient.NewClient(apiKey)

	// Set custom base URL if provided
	if baseURL != "" {
		client.BaseURL = baseURL
	}

	// Create a message request
	req := &models.SendMessageRequest{
		To:        []string{"recipient@example.com"},
		From:      "sender@yourdomain.com",
		Subject:   "Hello from Postal API",
		PlainBody: "This is a test email sent using the Postal API Go client.",
		HTMLBody:  "<p>This is a test email sent using the <strong>Postal API Go client</strong>.</p>",
		Headers: map[string]string{
			"X-Custom-Header": "Custom Value",
		},
	}

	// Send the message
	resp, err := client.SendMessage(req)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Printf("Message sent successfully! Message ID: %d, Token: %s\n", resp.MessageID, resp.Token)

	// Get message details
	message, err := client.GetMessage(resp.MessageID)
	if err != nil {
		log.Fatalf("Error getting message: %v", err)
	}

	fmt.Printf("Message details - ID: %d, Token: %s\n", message.ID, message.Token)

	// Get message deliveries
	deliveries, err := client.GetMessageDeliveries(resp.MessageID)
	if err != nil {
		log.Fatalf("Error getting message deliveries: %v", err)
	}

	fmt.Printf("Message has %d deliveries\n", len(deliveries))
	for i, delivery := range deliveries {
		fmt.Printf("Delivery %d - Status: %s, Details: %s\n", i+1, delivery.Status, delivery.Details)
	}
}
