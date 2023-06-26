package lib

import "github.com/golang-queue/queue"

var Queue *queue.Queue

func init() {
	// Proccess only one queue
	Queue = queue.NewPool(0)

	Queue.Start()
}
