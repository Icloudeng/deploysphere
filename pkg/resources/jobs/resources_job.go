package jobs

import (
	"context"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/queue"
	"smatflow/platform-installer/pkg/resources/db"

	goqueue "github.com/golang-queue/queue"
)

type ResourcesJob struct {
	Ref           string
	Task          goqueue.TaskFunc
	PostBody      interface{}
	ResourceState bool
	Description   string
	Group         string
}

func ResourcesJobTask(task ResourcesJob) {
	// Create new JOB
	job := db.JobCreate(db.JobCreateParam{
		Ref:         task.Ref,
		PostBody:    task.PostBody,
		Description: task.Description,
		Group:       task.Group,
	})

	queue.Queue.QueueTask(func(ctx context.Context) error {
		res_state := &database.ResourcesState{}

		// Listen to Redis Provisining events
		close := redis_pub_listeners(task.Ref)
		defer close()

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
