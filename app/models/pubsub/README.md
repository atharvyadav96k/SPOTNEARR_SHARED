# Google Cloud Pub/Sub Wrapper

A lightweight Go wrapper for Google Cloud Pub/Sub. Provides a clean, simple interface for publishing messages to Pub/Sub topics.

## Overview

The Pub/Sub wrapper provides:
- **Initialization**: One-time client setup with project ID
- **Publishing**: Send messages to Pub/Sub topics with automatic serialization
- **Resource Management**: Proper client lifecycle management
- **Error Handling**: Consistent error handling with custom error types

## Setup

### Prerequisites

1. Set up Google Cloud authentication. The SDK uses Application Default Credentials (ADC):

```bash
# If using gcloud CLI
gcloud auth application-default login

# Or set GOOGLE_APPLICATION_CREDENTIALS environment variable
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

2. Add the dependency to your `go.mod`:

```bash
go get cloud.google.com/go/pubsub
```

### Initialization

Initialize the Pub/Sub service in your application:

```go
package main

import (
	"context"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
)

func main() {
	ctx := context.Background()
	
	// Initialize Pub/Sub client
	client, err := pubsub.InitPubSub(ctx, "your-gcp-project-id")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	
	// Create a service wrapper
	pubsubService := &pubsub.Service{Client: client}
	
	// Now you can use the service to publish messages
}
```

## Available Methods

### InitPubSub

Initializes a new Pub/Sub client for the given project.

**Signature:**
```go
func InitPubSub(ctx context.Context, projectId string) (*pubsub.Client, error)
```

**Parameters:**
- `ctx`: context for initialization timeout and cancellation
- `projectId`: Google Cloud project ID

**Returns:**
- `*pubsub.Client`: initialized Pub/Sub client
- `error`: if initialization fails

**Note:** Caller is responsible for closing the client when done.

---

### Service.Publish

Publishes a message to a Pub/Sub topic.

**Signature:**
```go
func (s *Service) Publish(ctx context.Context, topicName string, payload any) error
```

**Parameters:**
- `ctx`: request context (controls timeout and cancellation)
- `topicName`: name of the Pub/Sub topic to publish to
- `payload`: data to be serialized and published (can be any type - will be converted to JSON)

**Returns:**
- `error`: nil on success, error if client is not initialized, serialization fails, or publish fails

**Notes:**
- The payload is automatically serialized to JSON before publishing
- This implementation waits for publish confirmation, making it synchronous
- For high-throughput systems, consider implementing an async version

---

### Service.Close

Shuts down the Pub/Sub client and releases resources.

**Signature:**
```go
func (s *Service) Close() error
```

**Returns:**
- `error`: nil on success, error if closing fails
- It is safe to call multiple times
- If the client is not initialized, returns nil

---

## Usage Examples

### Basic Message Publishing

```go
package main

import (
	"context"
	"log"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
)

func main() {
	ctx := context.Background()
	
	// Initialize client
	client, err := pubsub.InitPubSub(ctx, "my-project-id")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	pubsubService := &pubsub.Service{Client: client}
	
	// Publish a simple string message
	err = pubsubService.Publish(ctx, "my-topic", "Hello, World!")
	if err != nil {
		log.Fatal(err)
	}
}
```

### Publishing Structured Data

```go
type UserEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	EventType string `json:"event_type"`
	Timestamp int64  `json:"timestamp"`
}

// Publish user registration event
event := UserEvent{
	UserID:    "user-123",
	Email:     "user@example.com",
	EventType: "user.registered",
	Timestamp: time.Now().Unix(),
}

err := pubsubService.Publish(ctx, "user-events", event)
if err != nil {
	log.Printf("Failed to publish event: %v", err)
}
```

### Publishing with Timeout Context

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

message := map[string]interface{}{
	"order_id": 12345,
	"status":   "processing",
}

err := pubsubService.Publish(ctx, "order-events", message)
if err != nil {
	if err == context.DeadlineExceeded {
		log.Println("Publish operation timed out")
	} else {
		log.Printf("Publish failed: %v", err)
	}
}
```

### Multiple Topics

```go
// Publish to different topics
err1 := pubsubService.Publish(ctx, "user-events", userData)
err2 := pubsubService.Publish(ctx, "order-events", orderData)
err3 := pubsubService.Publish(ctx, "audit-logs", auditData)

// Check for errors
for _, err := range []error{err1, err2, err3} {
	if err != nil {
		log.Printf("Publish error: %v", err)
	}
}
```

## Error Handling

### Common Errors

The wrapper uses standard error types from `common.error` package:

- `ErrPubSubClientNotInitialized`: Client not initialized or nil
- `ErrInvalidTopic`: Topic does not exist
- Serialization errors: JSON encoding failures
- Context errors: Timeout or cancellation

### Error Handling Pattern

```go
err := pubsubService.Publish(ctx, "my-topic", myData)
if err != nil {
	switch err {
	case bus_error.ErrPubSubClientNotInitialized:
		log.Println("Pub/Sub client not initialized")
	case bus_error.ErrInvalidTopic:
		log.Println("Topic does not exist")
	case context.DeadlineExceeded:
		log.Println("Publish operation timed out")
	default:
		log.Printf("Publish failed: %v", err)
	}
}
```

## Best Practices

1. **Reuse Client**: Create the Pub/Sub client once and reuse it throughout your application
2. **Proper Cleanup**: Always call `Close()` when shutting down the application
3. **Structured Data**: Use structs to represent events for better type safety
4. **Error Checking**: Always check errors returned from `Publish()`
5. **Timeouts**: Use contexts with timeouts for production code
6. **Topic Names**: Ensure topics exist in your Pub/Sub instance before publishing

## Service Structure

```go
type Service struct {
	// Client is the underlying Pub/Sub client used for publishing
	// and subscribing to messages.
	Client *pubsub.Client

	// once ensures the client is initialized only once.
	once sync.Once
}
```

## See Also

- [Google Cloud Pub/Sub Documentation](https://cloud.google.com/pubsub/docs)
- [Go Client Library](https://cloud.google.com/go/docs/reference/cloud.google.com/go/pubsub)
