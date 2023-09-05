package main

import (
	"context"
	"fmt"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/pubsub/redis"
)

func main() {
	ctx := context.Background()
	data := make(chan string)

	go func() {
		sub := redis.Client.Subscribe(ctx, pubsub.REDIS_EVENT_TYPE_QUEUE_TASK_STATUS_RESPONSE)
		defer sub.Close()

		if msg, err := sub.ReceiveMessage(ctx); err != nil {
			fmt.Println("Failed to get queue task status, ", err.Error())
			return
		} else {
			data <- msg.Payload
		}
	}()

	redis.Client.Publish(ctx, pubsub.REDIS_EVENT_TYPE_QUEUE_TASK_STATUS_REQUEST, nil)

	fmt.Println(<-data)
}
