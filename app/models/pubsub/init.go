package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

func InitPubSub(projectId string) (*pubsub.Client, error) {
	var err error
	ctx := context.Background()
	Instance.Once.Do(func() {
		Instance.Client, err = pubsub.NewClient(ctx, projectId)
	})
	if err != nil {
		return nil, err
	}
	return Instance.Client, nil
}
