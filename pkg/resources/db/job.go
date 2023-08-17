package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/pubsub"
	"time"

	"gorm.io/datatypes"
)

type JobCreateParam struct {
	Ref         string
	PostBody    interface{}
	Description string
	Group       string
	Status      string
	Handler     string
	Method      string
}

func JobCreate(data JobCreateParam) *database.Job {
	rep := database.JobRepository{}
	postBodyJson, _ := json.Marshal(data.PostBody)

	job := &database.Job{
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

func JobUpdateStatus(job *database.Job, Status string) *database.Job {
	rep := database.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	if job.Status != database.JOB_STATUS_FAILED {
		job.Status = Status
	}

	if Status == database.JOB_STATUS_COMPLETED || Status == database.JOB_STATUS_FAILED {
		job.FinishedAt = time.Now()
	}

	rep.UpdateOrCreate(job)

	return job
}

// =============== Redis Events Listener ============= //

func Job_ListenResourceProviningLogs(playload pubsub.NetworkEventPayload) {
	rep := database.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	job.Logs = fmt.Sprintf("%s%s\\n", job.Logs, string(decodedBytes))

	rep.UpdateOrCreate(job)
}

func Job_ListenResourceProviningStatus(playload pubsub.NetworkEventPayload) {
	rep := database.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	if string(decodedBytes) == "succeeded" {
		job.Status = database.JOB_STATUS_COMPLETED
	} else {
		job.Status = database.JOB_STATUS_FAILED
	}

	rep.UpdateOrCreate(job)
}
