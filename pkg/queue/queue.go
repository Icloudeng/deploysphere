package queue

import (
	"context"
	"fmt"

	"github.com/golang-queue/queue"
)

var Queue *queue.Queue

func init() {
	// Proccess only one queue
	Queue = queue.NewPool(1)

	Queue.QueueTask(func(ctx context.Context) error {
		fmt.Println("Default queue task. Done!")
		return nil
	})
}
