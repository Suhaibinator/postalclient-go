package message

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// Bound the requests with a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Printf("Message sent successfully! Message-ID: %s\n", resp.MessageID)

	// Postal creates one message per recipient; look each one up by its numeric ID
	for rcpt, m := range resp.Messages {
		message, err := client.GetMessage(ctx, m.ID, models.ExpansionStatus)
		if err != nil {
			log.Fatalf("Error getting message: %v", err)
		}

		fmt.Printf("Message for %s - ID: %d, Token: %s, Status: %s\n", rcpt, message.ID, message.Token, message.Status.Status)

		deliveries, err := client.GetMessageDeliveries(ctx, m.ID)
		if err != nil {
			log.Fatalf("Error getting message deliveries: %v", err)
		}

		fmt.Printf("Message has %d deliveries\n", len(deliveries))
		for i, delivery := range deliveries {
			fmt.Printf("Delivery %d - Status: %s, Details: %s\n", i+1, delivery.Status, delivery.Details)
		}
	}
}
