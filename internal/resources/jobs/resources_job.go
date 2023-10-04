package jobs

import (
	"context"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/queue"
	"github.com/icloudeng/platform-installer/internal/resources/db"
	"github.com/icloudeng/platform-installer/internal/resources/websocket"
)

type TaskFunc func(context.Context, entities.Job) error

type ResourcesJob struct {
	Ref           string
	Task          TaskFunc
	PostBody      interface{}
	ResourceState bool
	Method        string
	Description   string
	Group         string
	Handler       string
}

func ResourcesJobTask(task ResourcesJob) *entities.Job {
	// Create new JOB
	job := db.Jobs.JobCreate(db.JobCreateParam{
		Ref:         task.Ref,
		PostBody:    task.PostBody,
		Description: task.Description,
		Group:       task.Group,
		Handler:     task.Handler,
		Method:      task.Method,
		Status:      entities.JOB_STATUS_IDLE,
	})

	//Emit ws events
	websocket.EmitJobEvent(job)

	queue.Queue.QueueTask(func(ctx context.Context) error {
		res_state := &entities.ResourcesState{}

		job = db.Jobs.JobUpdateStatus(job, entities.JOB_STATUS_RUNNING)
		//Emit ws events
		websocket.EmitJobEvent(job)

		// Listen to Redis Provisining events
		close := redis_pub_listeners(task.Ref)
		defer close()

		// Create Resource State
		if task.ResourceState {
			res_state = db.ResourceState.ResourceStateCreate(task.Ref, *job)
		}

		// Run task
		err := task.Task(ctx, *job)

		if err == nil && task.ResourceState {
			db.ResourceState.ResourceStatePutTerraformState(res_state)
		}

		if err == nil {
			job = db.Jobs.JobUpdateStatus(job, entities.JOB_STATUS_COMPLETED)
		} else {
			job = db.Jobs.JobUpdateLogs(job, err.Error())
			job = db.Jobs.JobUpdateStatus(job, entities.JOB_STATUS_FAILED)
		}

		//Emit ws events
		websocket.EmitJobEvent(job)

		// Allocate memory for resource db backup
		go db.ResourcesBackup.CreateNewResourcesBackup()

		return nil
	})

	return job
}
