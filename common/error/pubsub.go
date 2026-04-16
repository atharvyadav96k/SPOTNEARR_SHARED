package bus_error

import "errors"

var (
	ErrPubSubClientNotInitialized = errors.New("pubsub client not initialized")
)
