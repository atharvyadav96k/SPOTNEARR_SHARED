package app

import "context"

func (a *App) PublishMessage(topic string, payload struct{}) {
	ctx := context.Background()
	a.PubSub.Publish(ctx, topic, payload)
}
