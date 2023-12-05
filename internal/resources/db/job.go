package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/pubsub"

	"gorm.io/datatypes"
)

type (
	JobCreateParam struct {
		Ref         string
		PostBody    interface{}
		Description string
		Group       string
		Status      string
		Handler     string
		Method      string
	}

	jobs struct{}
)

var Jobs jobs

func (jobs) JobCreate(data JobCreateParam) *entities.Job {
	rep := entities.JobRepository{}
	postBodyJson, _ := json.Marshal(data.PostBody)

	job := &entities.Job{
		Ref:         data.Ref,
		PostBody:    datatypes.JSON(postBodyJson),
		Description: data.Description,
		Group:       data.Group,
		Status:      data.Status,
		Handler:     data.Handler,
		Method:      data.Method,
	}

	rep.Create(job)

	return job
}

func (jobs) JobUpdateLogs(job *entities.Job, Logs string) *entities.Job {
	rep := entities.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	job.Logs = fmt.Sprintf("%s%s\\n", job.Logs, Logs)

	rep.UpdateOrCreate(job)

	return job
}

func (jobs) JobUpdateStatus(job *entities.Job, Status string) *entities.Job {
	rep := entities.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	if job.Status != entities.JOB_STATUS_FAILED {
		job.Status = Status
	}

	if Status == entities.JOB_STATUS_COMPLETED || Status == entities.JOB_STATUS_FAILED {
		job.FinishedAt = time.Now()
	}

	rep.UpdateOrCreate(job)

	return job
}

func (jobs) JobUpdatePostBody(job *entities.Job, PostBody interface{}) *entities.Job {
	rep := entities.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	// Update PostBody
	postBodyJson, _ := json.Marshal(PostBody)
	job.PostBody = datatypes.JSON(postBodyJson)

	rep.UpdateOrCreate(job)

	return job
}

func (jobs) JobGetByID(ID uint) *entities.Job {
	rep := entities.JobRepository{}
	return rep.Get(ID)
}

// =============== Redis Events Listener ============= //

func (jobs) Job_ListenResourceProviningLogs(playload pubsub.NetworkEventPayload) {
	rep := entities.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	job.Logs = fmt.Sprintf("%s%s\\n", job.Logs, string(decodedBytes))

	rep.UpdateOrCreate(job)
}

func (jobs) Job_ListenResourceProviningStatus(playload pubsub.NetworkEventPayload) {
	rep := entities.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	if string(decodedBytes) == "succeeded" {
		job.Status = entities.JOB_STATUS_COMPLETED
	} else {
		job.Status = entities.JOB_STATUS_FAILED
	}

	rep.UpdateOrCreate(job)
}
