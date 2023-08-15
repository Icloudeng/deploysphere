package queue

import (
	"context"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/resources/db"

	"github.com/golang-queue/queue"
)

type ResourceJob struct {
	Ref           string
	Task          queue.TaskFunc
	PostBody      interface{}
	ResourceState bool
	Description   string
}

func ResourceJobTask(task ResourceJob) {
	// Create new JOB
	job := db.JobCreate(task.Ref, task.PostBody, task.Description)

	Queue.QueueTask(func(ctx context.Context) error {
		res_state := &database.ResourcesState{}

		// Create Resource State
		if task.ResourceState {
			res_state = db.ResourceStateCreate(task.Ref, *job)
		}

		// Run task
		err := task.Task(ctx)
		if err == nil && task.ResourceState {
			db.ResourceStatePutTerraformState(res_state)
		}

		db.JobPutRunningDone(job, err == nil)

		return nil
	})

}
