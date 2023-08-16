package redis_events

import (
	"context"
	"fmt"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/redis"
)

type SubscriberFunc func(playload events.NetworkEventPayload)

func subscribe(reference string, eventType string, subscribers []SubscriberFunc) func() {
	ctx := context.Background()
	channel := reference + "-" + eventType
	pubsub := redis.Client.Subscribe(ctx, channel)

	go func() {
		// Close the subscription when we are done.
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("Close Redis Channel Subscribe", channel)
				break
			}

			// queue.JobsQueue.QueueTask(func(ctx context.Context) error {
			for _, subscriber := range subscribers {
				subscriber(events.NetworkEventPayload{
					Type:      eventType,
					Reference: reference,
					Channel:   msg.Channel,
					Payload:   msg.Payload,
				})
			}
			// 	return nil
			// })
		}
	}()

	return func() {
		pubsub.Close()
	}
}

func ResourceProviningLogsEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, events.REDIS_EVENT_TYPE_LOGS, subscribers)
}

func ResourceProviningCredentialsEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, events.REDIS_EVENT_TYPE_CREDENTIALS, subscribers)
}

func ResourceProviningStatusEvents(reference string, subscribers []SubscriberFunc) func() {
	return subscribe(reference, events.REDIS_EVENT_TYPE_STATUS, subscribers)
}
