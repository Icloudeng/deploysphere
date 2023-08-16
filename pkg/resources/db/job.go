package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events/redis_events"
	"time"

	"gorm.io/datatypes"
)

type JobCreateParam struct {
	Ref         string
	PostBody    interface{}
	Description string
	Group       string
}

func JobCreate(data JobCreateParam) *database.Job {
	rep := database.JobRepository{}
	postBodyJson, _ := json.Marshal(data.PostBody)

	job := &database.Job{
		Ref:         data.Ref,
		Running:     true,
		Success:     true,
		PostBody:    datatypes.JSON(postBodyJson),
		Description: data.Description,
		Group:       data.Group,
	}

	rep.Create(job)

	return job
}

func JobPutRunningDone(job *database.Job, Success bool) {
	rep := database.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	job.Running = false
	job.FinishedAt = time.Now()

	if !Success {
		job.Success = Success
	}

	rep.UpdateOrCreate(job)
}

// =============== Redis Events Listener ============= //

func Job_ListenResourceProviningLogs(playload redis_events.ResourceRedisEventPayload) {
	rep := database.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	job.Logs = fmt.Sprintf("%s%s\\n", job.Logs, string(decodedBytes))

	rep.UpdateOrCreate(job)
}

func Job_ListenResourceProviningStatus(playload redis_events.ResourceRedisEventPayload) {
	rep := database.JobRepository{}
	job := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if job == nil || err != nil {
		return
	}

	job.Success = string(decodedBytes) == "succeeded"

	rep.UpdateOrCreate(job)
}
