package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/atharvyadav96k/gcp/common"
	bus_error "github.com/atharvyadav96k/gcp/common/error"
)

func (p *PubSub) Publish(ctx context.Context, topicName string, payload any) error {
	if p.Client == nil {
		return bus_error.ErrPubSubClientNotInitialized
	}

	topic := p.Client.Topic(topicName)

	bytes, err := common.ToJSON(payload)
	if err != nil {
		return err
	}

	res := topic.Publish(ctx, &pubsub.Message{
		Data: bytes,
	})

	_, err = res.Get(ctx)
	return err
}

func (p *PubSub) Close() error {
	if p.Client != nil {
		return p.Client.Close()
	}
	return nil
}
