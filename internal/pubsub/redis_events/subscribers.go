package redis_events

import (
	"context"
	"fmt"
	"smatflow/platform-installer/internal/pubsub"
	"smatflow/platform-installer/internal/pubsub/redis"
)

type SubscriberFunc func(playload pubsub.NetworkEventPayload)

func subscribe(reference string, eventType string, subscribers []SubscriberFunc) func() {
	ctx := context.Background()
	channel := reference + "-" + eventType
	redis_pubsub := redis.Client.Subscribe(ctx, channel)

	go func() {
		// Close the subscription when we are done.
		defer redis_pubsub.Close()

		for {
			msg, err := redis_pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("Close Redis Channel Subscribe", channel)
				break
			}

			for _, subscriber := range subscribers {
				subscriber(pubsub.NetworkEventPayload{
					Type:      eventType,
					Reference: reference,
					Channel:   msg.Channel,
					Payload:   msg.Payload,
				})
			}
		}
	}()

	return func() {
		redis_pubsub.Close()
	}
}

func ResourceProviningLogsEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, pubsub.REDIS_EVENT_TYPE_LOGS, subscribers)
}

func ResourceProviningCredentialsEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, pubsub.REDIS_EVENT_TYPE_CREDENTIALS, subscribers)
}

func ResourceProviningStatusEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, pubsub.REDIS_EVENT_TYPE_STATUS, subscribers)
}
