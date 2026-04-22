package app

import (
	"github.com/atharvyadav96k/gcp/app/models/firestore"
	"github.com/atharvyadav96k/gcp/app/models/neon"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
	"github.com/atharvyadav96k/gcp/app/models/secrets"
)

// App is the central container for all application-level dependencies.
// It holds initialized clients for external services like Firestore,
// Pub/Sub, Neon PostgreSQL, and environment configuration.
//
// This struct should be created once (singleton-style) and reused
// across the application to avoid re-initializing clients.
type App struct {

	// Env contains environment variables and secrets required
	// for configuring external services (e.g., project IDs, keys).
	Env secrets.Env

	// FireStore provides access to the Firestore database client.
	// It manages initialization and lifecycle of the Firestore connection.
	FireStore *firestore.Service

	// PubSub provides access to Google Cloud Pub/Sub client.
	// It is used for publishing and subscribing to messages.
	PubSub *pubsub.Service

	// Neon provides access to the Neon PostgreSQL database client.
	// It manages initialization and lifecycle of the database connection.
	Neon *neon.Service
}

type UniqueField struct {
	Name  string
	Value interface{}
}

type StoreEntity interface {
	SetID(string)
}
