package app

import (
	"github.com/atharvyadav96k/gcp/app/models/firestore"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
	"github.com/atharvyadav96k/gcp/app/models/secrets"
)

type App struct {
	Env       secrets.Env
	FireStore *firestore.Firestore
	PubSub    *pubsub.PubSub
}
