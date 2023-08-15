package jobs

import (
	"context"
	"fmt"
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

func redisPubListeners(Ref string) {
	// Redis Event Logs Listeners
	close1 := redis_events.ResourceProviningLogsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Job_ListenResourceProviningLogs,
		},
	)
	defer close1()

	// Redis Event Status Listeners
	close2 := redis_events.ResourceProviningStatusEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.Job_ListenResourceProviningStatus,
		},
	)
	defer close2()

	// Redis Event Credentials Listeners
	close3 := redis_events.ResourceProviningCredentialsEvents(
		Ref,
		[]redis_events.SubscriberFunc{
			db.ResourceState_ListenResourceProviningCredentials,
		},
	)
	defer close3()
}

func ResourcesJobTask(task ResourcesJob) {
	// Create new JOB
	job := db.JobCreate(task.Ref, task.PostBody, task.Description)

	queue.Queue.QueueTask(func(ctx context.Context) error {
		fmt.Printf("==== Start Job QueueTask, ref: %s ====", task.Ref)
		res_state := &database.ResourcesState{}

		fmt.Printf("==== redis Pub Listeners, ref: %s ====", task.Ref)
		redisPubListeners(task.Ref)

		// Create Resource State
		if task.ResourceState {
			fmt.Printf("==== DB Create Resource state, ref: %s ====", task.Ref)
			res_state = db.ResourceStateCreate(task.Ref, *job)
		}

		// Run task
		fmt.Printf("==== Run Job task, ref: %s ====", task.Ref)
		err := task.Task(ctx)

		if err == nil && task.ResourceState {
			db.ResourceStatePutTerraformState(res_state)
		}

		db.JobPutRunningDone(job, err == nil)

		return nil
	})

}
