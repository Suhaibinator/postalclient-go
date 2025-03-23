package raw

import (
	"encoding/base64"
	"fmt"
	"log"

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

	// Send the raw message
	resp, err := client.SendRaw(req)
	if err != nil {
		log.Fatalf("Error sending raw message: %v", err)
	}

	fmt.Printf("Raw message sent successfully! Message ID: %d, Token: %s\n", resp.MessageID, resp.Token)

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
