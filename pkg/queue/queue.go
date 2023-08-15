package queue

import (
	"context"

	"github.com/golang-queue/queue"
)

var Queue *queue.Queue
var JobsQueue *queue.Queue

func init() {
	// Proccess only one queue
	Queue = queue.NewPool(1)
	JobsQueue = queue.NewPool(1)

	Queue.QueueTask(func(ctx context.Context) error {
		return nil
	})
	JobsQueue.QueueTask(func(ctx context.Context) error {
		return nil
	})
}
