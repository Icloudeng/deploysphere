package redis_events

import (
	"context"
	"smatflow/platform-installer/pkg/redis"
)

type ResourceRedisEventPayload struct {
	Type      string
	Reference string
	Channel   string
	Payload   string
}

type SubscriberFunc func(playload ResourceRedisEventPayload)

func events(reference string, eventType string, subscribers []SubscriberFunc) func() {
	ctx := context.Background()
	pubsub := redis.Client.Subscribe(ctx, reference+"-"+eventType)
	// Close the subscription when we are done.
	defer pubsub.Close()

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				break
			}

			for _, subscriber := range subscribers {
				subscriber(ResourceRedisEventPayload{
					Type:      eventType,
					Reference: reference,
					Channel:   msg.Channel,
					Payload:   msg.Payload,
				})
			}
		}
	}()

	return func() {
		pubsub.Close()
	}
}

func ResourceProviningLogsEvents(reference string, subscribers []SubscriberFunc) func() {
	return events(reference, "logs", subscribers)
}

func ResourceProviningCredentialsEvents(reference string, subscribers []SubscriberFunc) func() {
	return events(reference, "credentials", subscribers)
}

func ResourceProviningStatusEvents(reference string, subscribers []SubscriberFunc) func() {
	return events(reference, "status", subscribers)
}
