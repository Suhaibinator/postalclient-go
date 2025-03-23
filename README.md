# postalclient-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Suhaibinator/postalclient-go.svg)](https://pkg.go.dev/github.com/Suhaibinator/postalclient-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Suhaibinator/postalclient-go)](https://goreportcard.com/report/github.com/Suhaibinator/postalclient-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Coverage](https://img.shields.io/badge/coverage-92.6%25-brightgreen.svg)](https://github.com/Suhaibinator/postalclient-go)
[![Tests](https://github.com/Suhaibinator/postalclient-go/actions/workflows/tests.yml/badge.svg)](https://github.com/Suhaibinator/postalclient-go/actions/workflows/tests.yml)
[![Security](https://github.com/Suhaibinator/postalclient-go/actions/workflows/security.yml/badge.svg)](https://github.com/Suhaibinator/postalclient-go/actions/workflows/security.yml)

A comprehensive Go client library for the [Postal](https://github.com/postalserver/postal) API. This SDK provides a simple, idiomatic, and well-documented way to interact with the Postal API from Go applications.

## Features

- **Complete API Coverage**: Support for all Postal API endpoints
- **Type Safety**: Strongly typed request and response structures
- **Comprehensive Documentation**: Detailed GoDoc comments and examples
- **Error Handling**: Robust error handling with detailed error information
- **Customization**: Configurable client with options for base URL and timeout
- **Testing**: Extensive test coverage

## Installation

```bash
go get github.com/Suhaibinator/postalclient-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/Suhaibinator/postalclient-go"
    "github.com/Suhaibinator/postalclient-go/models"
)

func main() {
    // Create a client
    client := postalclient.NewClient("your-api-key")

    // Send a message
    resp, err := client.SendMessage(&models.SendMessageRequest{
        To:        []string{"recipient@example.com"},
        From:      "sender@yourdomain.com",
        Subject:   "Hello from Postal API",
        PlainBody: "This is a test email sent using the Postal API Go client.",
    })
    if err != nil {
        log.Fatalf("Error sending message: %v", err)
    }

    fmt.Printf("Message sent! ID: %d, Token: %s\n", resp.MessageID, resp.Token)
}
```

## Usage

### Creating a Client

```go
import (
    "time"
    "github.com/Suhaibinator/postalclient-go"
)

// Create a client with default options
client := postalclient.NewClient("your-api-key")

// Or with custom options
client := postalclient.NewClientWithOptions(
    "your-api-key",
    "https://postal.yourdomain.com/api/v1",
    60 * time.Second,
)
```

### Sending a Message

```go
import (
    "fmt"
    "log"
    "github.com/Suhaibinator/postalclient-go"
    "github.com/Suhaibinator/postalclient-go/models"
)

// Create a client
client := postalclient.NewClient("your-api-key")

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

// Use the response
fmt.Printf("Message sent! ID: %d, Token: %s\n", resp.MessageID, resp.Token)
```

### Sending a Raw RFC2822 Message

```go
import (
    "encoding/base64"
    "fmt"
    "log"
    "github.com/Suhaibinator/postalclient-go"
    "github.com/Suhaibinator/postalclient-go/models"
)

// Create a client
client := postalclient.NewClient("your-api-key")

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

// Use the response
fmt.Printf("Raw message sent! ID: %d, Token: %s\n", resp.MessageID, resp.Token)
```

### Getting Message Details

```go
import (
    "fmt"
    "log"
    "github.com/Suhaibinator/postalclient-go"
)

// Create a client
client := postalclient.NewClient("your-api-key")

// Get message details
message, err := client.GetMessage(messageID)
if err != nil {
    log.Fatalf("Error getting message: %v", err)
}

// Use the message
fmt.Printf("Message details - ID: %d, Token: %s\n", message.ID, message.Token)

// Access message content if available
if message.PlainBody != "" {
    fmt.Printf("Plain body: %s\n", message.PlainBody)
}

if message.HTMLBody != "" {
    fmt.Printf("HTML body: %s\n", message.HTMLBody)
}

// Access message headers if available
if len(message.Headers) > 0 {
    fmt.Println("Headers:")
    for key, value := range message.Headers {
        fmt.Printf("  %s: %s\n", key, value)
    }
}
```

### Getting Message Deliveries

```go
import (
    "fmt"
    "log"
    "github.com/Suhaibinator/postalclient-go"
)

// Create a client
client := postalclient.NewClient("your-api-key")

// Get message deliveries
deliveries, err := client.GetMessageDeliveries(messageID)
if err != nil {
    log.Fatalf("Error getting message deliveries: %v", err)
}

// Use the deliveries
fmt.Printf("Message has %d deliveries\n", len(deliveries))
for i, delivery := range deliveries {
    fmt.Printf("Delivery %d:\n", i+1)
    fmt.Printf("  Status: %s\n", delivery.Status)
    fmt.Printf("  Details: %s\n", delivery.Details)
    fmt.Printf("  Timestamp: %s\n", delivery.Timestamp.Format(time.RFC3339))
    fmt.Printf("  Sent with SSL: %v\n", delivery.SentWithSSL)
    fmt.Printf("  Time taken: %.2f seconds\n", delivery.Time)
}
```

## Error Handling

The client returns detailed error information when API requests fail:

```go
resp, err := client.SendMessage(req)
if err != nil {
    // Check if it's a Postal API error
    if apiErr, ok := err.(*postalclient.Error); ok {
        fmt.Printf("API Error: %s - %s\n", apiErr.Status, apiErr.Message)
        // Handle specific error types
        if apiErr.Status == "parameter-error" {
            // Handle parameter validation errors
        } else if apiErr.Status == "error" {
            // Handle general API errors
        }
    } else {
        // Handle other errors (network, parsing, etc.)
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Examples

See the [examples](./examples) directory for complete examples of how to use this library:

- [Sending a Message](./examples/message/send_message.go)
- [Sending a Raw RFC2822 Message](./examples/raw/send_raw.go)

To run the examples:

```bash
# Send a regular message
go run examples/main.go --api-key=your-api-key --example=message

# Send a raw RFC2822 message
go run examples/main.go --api-key=your-api-key --example=raw
```

## API Documentation

For more information about the Postal API, see the [official documentation](https://github.com/postalserver/postal/wiki/HTTP-API).

## CI/CD Workflows

This project uses GitHub Actions for continuous integration and delivery:

- **CI**: Runs linting and code formatting checks on every push and pull request
- **Tests**: Runs unit tests and generates code coverage reports
- **Release**: Automates the release process when a new tag is pushed
- **Security**: Scans dependencies for vulnerabilities and checks for outdated packages

### Status

[![CI](https://github.com/Suhaibinator/postalclient-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Suhaibinator/postalclient-go/actions/workflows/ci.yml)
[![Tests](https://github.com/Suhaibinator/postalclient-go/actions/workflows/tests.yml/badge.svg)](https://github.com/Suhaibinator/postalclient-go/actions/workflows/tests.yml)
[![Security](https://github.com/Suhaibinator/postalclient-go/actions/workflows/security.yml/badge.svg)](https://github.com/Suhaibinator/postalclient-go/actions/workflows/security.yml)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on how to contribute to this project, including:

- Development workflow
- Setting up the development environment
- Running tests and linting
- Style guide
- Documentation guidelines

All contributors are expected to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md).

When contributing to this repository, please:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Write tests for your changes
4. Update documentation as needed
5. Submit a pull request

The CI workflow will automatically run tests and linting on your pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
