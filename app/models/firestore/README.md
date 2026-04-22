# Google Cloud Firestore Wrapper

A lightweight Go wrapper for Google Cloud Firestore. Provides a clean interface for initializing and managing the Firestore client.

## Overview

The Firestore wrapper provides:
- **Initialization**: One-time client setup with automatic authentication
- **Client Management**: Proper lifecycle management using sync.Once for singleton pattern
- **Resource Cleanup**: Graceful shutdown and resource release
- **Error Handling**: Consistent error handling

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
go get cloud.google.com/go/firestore
```

### Initialization

Initialize the Firestore client in your application:

```go
package main

import (
	"github.com/atharvyadav96k/gcp/app/models/firestore"
)

func main() {
	// Initialize Firestore client
	client, err := firestore.InitFirestore("your-gcp-project-id")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	
	// Create a service wrapper
	firestoreService := &firestore.Service{Client: client}
	
	// Now you can use the client to interact with Firestore
}
```

## Available Methods

### InitFirestore

Creates and returns a new Firestore client for the given project.

**Signature:**
```go
func InitFirestore(projectId string) (*firestore.Client, error)
```

**Parameters:**
- `projectId`: Google Cloud project ID

**Returns:**
- `*firestore.Client`: initialized Firestore client
- `error`: if client initialization fails

**Notes:**
- Uses application default credentials for authentication
- Creates client in context.Background()
- Caller is responsible for closing the client using `client.Close()`

---

### Service.Close

Shuts down the Firestore client and releases underlying resources.

**Signature:**
```go
func (f *Service) Close() error
```

**Returns:**
- `error`: nil on success, error if closing fails
- It is safe to call multiple times
- If the client is not initialized, returns nil

---

## Usage Examples

### Basic Setup

```go
package main

import (
	"context"
	"log"
	"github.com/atharvyadav96k/gcp/app/models/firestore"
)

func main() {
	// Initialize Firestore client
	client, err := firestore.InitFirestore("my-project-id")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	firestoreService := &firestore.Service{Client: client}
	
	// Use the client for Firestore operations
	ctx := context.Background()
	collection := firestoreService.Client.Collection("users")
	
	// Query documents
	docs, err := collection.Documents(ctx).GetAll()
	if err != nil {
		log.Fatal(err)
	}
	
	for _, doc := range docs {
		log.Printf("Document: %s", doc.Ref.ID)
	}
}
```

### Creating a Document

```go
ctx := context.Background()

// Define data structure
type User struct {
	Name  string
	Email string
	Age   int
}

user := User{
	Name:  "John Doe",
	Email: "john@example.com",
	Age:   30,
}

// Add document to Firestore
_, err := firestoreService.Client.Collection("users").Doc("user-123").Set(ctx, user)
if err != nil {
	log.Printf("Failed to add user: %v", err)
}
```

### Reading a Document

```go
ctx := context.Background()

// Get a single document
doc, err := firestoreService.Client.Collection("users").Doc("user-123").Get(ctx)
if err != nil {
	log.Printf("Failed to get user: %v", err)
	return
}

var user User
if err := doc.DataTo(&user); err != nil {
	log.Printf("Failed to unmarshal user: %v", err)
}

log.Printf("User: %+v", user)
```

### Querying Documents

```go
ctx := context.Background()

// Query for all users older than 25
query := firestoreService.Client.Collection("users").Where("age", ">=", 25)
docs, err := query.Documents(ctx).GetAll()
if err != nil {
	log.Printf("Query failed: %v", err)
	return
}

for _, doc := range docs {
	var user User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("Failed to unmarshal: %v", err)
		continue
	}
	log.Printf("User: %+v", user)
}
```

### Updating a Document

```go
ctx := context.Background()

// Update specific fields
_, err := firestoreService.Client.Collection("users").Doc("user-123").Update(ctx, []firestore.Update{
	{
		Path:  "age",
		Value: 31,
	},
	{
		Path:  "lastModified",
		Value: time.Now(),
	},
})
if err != nil {
	log.Printf("Failed to update user: %v", err)
}
```

### Deleting a Document

```go
ctx := context.Background()

// Delete a document
_, err := firestoreService.Client.Collection("users").Doc("user-123").Delete(ctx)
if err != nil {
	log.Printf("Failed to delete user: %v", err)
}
```

### Batch Operations

```go
ctx := context.Background()

// Create a batch for multiple operations
batch := firestoreService.Client.Batch()

// Add multiple documents
batch.Set(firestoreService.Client.Collection("users").Doc("user-1"), 
	map[string]interface{}{"name": "John"})
batch.Set(firestoreService.Client.Collection("users").Doc("user-2"), 
	map[string]interface{}{"name": "Jane"})

// Commit the batch
_, err := batch.Commit(ctx)
if err != nil {
	log.Printf("Batch commit failed: %v", err)
}
```

### Transactions

```go
ctx := context.Background()

// Use transaction for multi-document operations
err := firestoreService.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
	// Read
	doc1, err := tx.Get(firestoreService.Client.Collection("users").Doc("user-1"))
	if err != nil {
		return err
	}
	
	// Write
	err = tx.Set(firestoreService.Client.Collection("users").Doc("user-2"), 
		map[string]interface{}{"name": "Jane"})
	if err != nil {
		return err
	}
	
	return nil
})

if err != nil {
	log.Printf("Transaction failed: %v", err)
}
```

## Firestore Client Operations

Once initialized, the `Service.Client` provides access to all standard Firestore operations:

### Collection Operations
- `Collection(path)` - Get a reference to a collection
- `Collections()` - List all collections

### Document Operations
- `Doc(path)` - Get a reference to a document
- `Get(ctx)` - Retrieve a document
- `Set(ctx, data)` - Create or overwrite a document
- `Update(ctx, updates)` - Update specific fields
- `Delete(ctx)` - Delete a document

### Query Operations
- `Where(path, op, value)` - Add filter condition
- `OrderBy(path)` - Order results
- `Limit(n)` - Limit number of results
- `Offset(n)` - Skip number of results
- `Documents(ctx)` - Get document snapshots
- `Count(ctx)` - Count matching documents

### Batch Operations
- `Batch()` - Create a batch for multiple operations

### Transactions
- `RunTransaction(ctx, fn)` - Execute operations in a transaction

## Error Handling

### Common Errors

- Connection errors: Project ID not found, authentication failed
- Document not found: Using `Status(err) == codes.NotFound`
- Permission denied: Insufficient Firestore permissions
- Invalid arguments: Incorrect field paths or query conditions

### Error Handling Pattern

```go
ctx := context.Background()
doc, err := firestoreService.Client.Collection("users").Doc("user-123").Get(ctx)
if err != nil {
	if status.Code(err) == codes.NotFound {
		log.Println("Document not found")
	} else {
		log.Printf("Failed to get document: %v", err)
	}
	return
}
```

## Best Practices

1. **Reuse Client**: Create the Firestore client once and reuse throughout your application
2. **Proper Cleanup**: Always call `Close()` when shutting down
3. **Use Transactions**: For multi-document operations that need consistency
4. **Batch Operations**: Use batch writes for better performance with multiple documents
5. **Structured Data**: Use Go structs with `firestore` tags for better type safety
6. **Context Timeouts**: Use contexts with timeouts for production code
7. **Error Checking**: Always check and handle errors from Firestore operations

## Service Structure

```go
type Service struct {
	// Client is the underlying Firestore client used to interact
	// with the database (collections, documents, queries, etc.).
	Client *firestore.Client

	// Once guarantees that the Firestore client is initialized only once,
	// even if Init is called multiple times from different goroutines.
	Once sync.Once
}
```

## See Also

- [Google Cloud Firestore Documentation](https://cloud.google.com/firestore/docs)
- [Go Client Library](https://cloud.google.com/go/docs/reference/cloud.google.com/go/firestore)
- [Firestore Best Practices](https://cloud.google.com/firestore/docs/best-practices)
