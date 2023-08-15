package database

import (
	"encoding/json"
	"smatflow/platform-installer/pkg/database"

	"gorm.io/datatypes"
)

func JobCreate(ref string, postBody interface{}) *database.Job {
	rep := database.JobRepository{}
	postBodyJson, _ := json.Marshal(postBody)

	job := &database.Job{
		Ref:      ref,
		Running:  true,
		Success:  true,
		PostBody: datatypes.JSON(postBodyJson),
	}

	rep.Create(job)

	return job
}

func JobPutRunningDone(job *database.Job, Success bool) {
	rep := database.JobRepository{}

	job.Running = false

	job.Success = Success

	rep.UpdateOrCreate(job)
}
