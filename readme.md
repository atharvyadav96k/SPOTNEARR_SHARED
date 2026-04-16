# GCP Library

## Overview
The Spotnearr Shared Library is a Go-based library designed to simplify the integration of Google Cloud Platform (GCP) services such as Firestore, Pub/Sub, and Secrets Management. This library provides a structured and reusable way to initialize and manage these services in your application.

## Features
- **Firestore Integration**: Easily initialize and manage Firestore clients.
- **Pub/Sub Integration**: Publish and subscribe to messages with Pub/Sub.
- **Secrets Management**: Load and manage environment variables securely.
- **Thread-Safe Initialization**: Uses `sync.Once` to ensure one-time initialization of clients.
- **Resource Cleanup**: Provides methods to close clients and release resources.

## Folder Structure
```
app/
  init.go          # Main entry point for initializing the library
  models.go        # Defines the App struct and its dependencies
  models/
    firestore/     # Firestore-related logic
      init.go      # Firestore client initialization
      methods.go   # Firestore utility methods
      models.go    # Firestore struct definition
    pubsub/        # Pub/Sub-related logic
      init.go      # Pub/Sub client initialization
      methods.go   # Pub/Sub utility methods
      modules.go   # Pub/Sub struct definition
    secrets/       # Secrets management logic
      getter.go    # Secrets getter methods
      init.go      # Secrets initialization
      module.go    # Secrets struct definition
```

## Getting Started

### Prerequisites
- Go 1.18 or higher
- Google Cloud SDK installed and authenticated
- Environment variable `GCP_PROJECT_ID` set to your GCP project ID

### Installation
Clone the repository:
```bash
git clone https://github.com/atharvyadav96k/gcp.git
```

Navigate to the project directory:
```bash
cd spotnearr-shared
```

### Importing the Library
To use this library in your Go project, import it as follows:
```go
import "github.com/atharvyadav96k/gcp/app"
```

### Usage

#### Initialize the Library
```go
package main

import (
	"github.com/atharvyadav96k/gcp/app/app"
)

func main() {
	// Initialize the App
	appInstance := app.Init()

	// Initialize Environment Variables
	appInstance.InitEnvironmentVariables()

	// Initialize Firestore
	if err := appInstance.InitFirestore(); err != nil {
		panic(err)
	}

	// Initialize Pub/Sub
	if err := appInstance.InitPubSub(appInstance.Env.GetProjectId()); err != nil {
		panic(err)
	}

	// Close resources on exit
	defer appInstance.Close()
}
```

#### Firestore Example
```go
// Use Firestore client
firestoreClient := appInstance.FireStore.FirestoreClient
// Perform Firestore operations...
```

#### Pub/Sub Example
```go
// Publish a message
err := appInstance.PubSub.Publish(ctx, "topic-name", payload)
if err != nil {
	log.Fatalf("Failed to publish message: %v", err)
}
```

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments
- [Google Cloud Go SDK](https://pkg.go.dev/cloud.google.com/go)