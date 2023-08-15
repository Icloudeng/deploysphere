package jobs

import (
	"context"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events/redis_events"
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
}

func redisPubListeners(Ref string) func() {
	// Redis Event Logs Listeners
	close1 := redis_events.ResourceProviningLogsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Job_ListenResourceProviningLogs,
		},
	)

	// Redis Event Status Listeners
	close2 := redis_events.ResourceProviningStatusEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Job_ListenResourceProviningStatus,
		},
	)

	// Redis Event Credentials Listeners
	close3 := redis_events.ResourceProviningCredentialsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.ResourceState_ListenResourceProviningCredentials,
		},
	)

	return func() {
		close1()
		close3()
		close2()
	}
}

func ResourcesJobTask(task ResourcesJob) {
	// Create new JOB
	job := db.JobCreate(task.Ref, task.PostBody, task.Description)

	queue.Queue.QueueTask(func(ctx context.Context) error {
		res_state := &database.ResourcesState{}

		// Listen to Redis Provisining events
		close := redisPubListeners(task.Ref)
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
