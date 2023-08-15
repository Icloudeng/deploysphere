package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events/redis_events"

	"gorm.io/datatypes"
)

func JobCreate(ref string, postBody interface{}, Description string) *database.Job {
	rep := database.JobRepository{}
	postBodyJson, _ := json.Marshal(postBody)

	job := &database.Job{
		Ref:         ref,
		Running:     true,
		Success:     true,
		PostBody:    datatypes.JSON(postBodyJson),
		Description: Description,
	}

	rep.Create(job)

	return job
}

func JobPutRunningDone(job *database.Job, Success bool) {
	rep := database.JobRepository{}
	// refresh the job
	job = rep.Get(job.ID)

	job.Running = false

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

	// job.Logs = fmt.Sprintf("%s%s\\n", job.Logs, string(decodedBytes))
	job.Logs = fmt.Sprintf("%s%s\n", job.Logs, string(decodedBytes))

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
