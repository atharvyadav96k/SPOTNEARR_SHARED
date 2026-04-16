package pubsub

import (
	"sync"

	"cloud.google.com/go/pubsub"
)

type PubSub struct {
	Client *pubsub.Client
	Once   sync.Once
}

// Instance is the global singleton
var Instance PubSub
