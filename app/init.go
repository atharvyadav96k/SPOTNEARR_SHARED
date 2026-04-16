package app

import (
	"fmt"

	"github.com/atharvyadav96k/gcp/app/models/firestore"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
	"github.com/atharvyadav96k/gcp/app/models/secrets"
)

func (a *App) InitEnvironmentVariables() {
	a.Env = secrets.NewSecrets()
}

func (a *App) Init() *App {
	return &App{
		FireStore: &firestore.Firestore{},
		PubSub:    &pubsub.PubSub{},
	}
}

func (a *App) InitFirestore() error {
	var err error
	a.FireStore.Once.Do(func() {
		client, initErr := firestore.InitFirestore(a.Env.GCP_PROJECT_ID)
		if initErr != nil {
			err = initErr
			return
		}
		a.FireStore.FirestoreClient = client
	})
	return err
}

func (a *App) InitPubSub(projectId string) error {
	var err error
	a.PubSub.Once.Do(func() {
		client, initErr := pubsub.InitPubSub(projectId)
		if initErr != nil {
			err = initErr
			return
		}
		a.PubSub.Client = client
	})
	return err
}

func (a *App) Close() {
	err := a.FireStore.Close()
	if err != nil {
		fmt.Printf("Error closing Firestore client: %v\n", err)
	}
	err = a.PubSub.Close()
	if err != nil {
		fmt.Printf("Error closing Pubsub client %v\n", err)
	}
}
