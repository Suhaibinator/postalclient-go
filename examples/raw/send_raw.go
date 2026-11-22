package raw

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Suhaibinator/postalclient-go"
	"github.com/Suhaibinator/postalclient-go/models"
)

// SendExample demonstrates how to send a raw RFC2822 message using the Postal API
func SendExample(apiKey, baseURL string) {
	// Create a new client
	client := postalclient.NewClient(apiKey)

	// Set custom base URL if provided
	if baseURL != "" {
		client.BaseURL = baseURL
	}

	// Create a raw RFC2822 message
	rawMessage := `From: sender@yourdomain.com
To: recipient@example.com
Subject: Hello from Postal API
Content-Type: text/plain; charset=utf-8

This is a test email sent using the Postal API Go client with a raw RFC2822 message.
`

	// Encode the message in base64
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(rawMessage))

	// Create a raw message request
	req := &models.SendRawRequest{
		MailFrom: "sender@yourdomain.com",
		RcptTo:   []string{"recipient@example.com"},
		Data:     encodedMessage,
	}

	// Bound the requests with a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.SendRaw(ctx, req)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	fmt.Printf("Raw message sent successfully! Message-ID: %s\n", resp.MessageID)

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
