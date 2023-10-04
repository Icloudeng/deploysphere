package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/icloudeng/platform-installer/internal/pubsub"
	"github.com/icloudeng/platform-installer/internal/pubsub/redis"
	"github.com/icloudeng/platform-installer/internal/queue"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func init() {
	ctx := context.Background()
	sub := redis.Client.Subscribe(ctx, pubsub.REDIS_EVENT_TYPE_QUEUE_TASK_STATUS_REQUEST)

	go func() {
		// Close the subscription when we are done.
		defer sub.Close()

		for {
			msg, err := sub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println("Close Redis Channel Subscribe, ", err.Error())
				break
			}

			log.Println("Received redis event: ", msg.Channel)
			time.Sleep(time.Second)

			payload, err := json.MarshalIndent(structs.QueueStatusRequest{
				BusyWorkers:    queue.Queue.BusyWorkers(),
				FailureTasks:   queue.Queue.FailureTasks(),
				SubmittedTasks: queue.Queue.SubmittedTasks(),
				SuccessTasks:   queue.Queue.SuccessTasks(),
			}, "", "    ")

			if err == nil {
				// Send status queue
				redis.Client.Publish(ctx, pubsub.REDIS_EVENT_TYPE_QUEUE_TASK_STATUS_RESPONSE, payload)
			}
		}
	}()
}
