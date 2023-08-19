package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	JOB_STATUS_IDLE = "idle"

	JOB_STATUS_COMPLETED = "completed"

	JOB_STATUS_FAILED = "failed"

	JOB_STATUS_RUNNING = "running"
)

type Job struct {
	gorm.Model
	Ref         string `gorm:"index"`
	PostBody    datatypes.JSON
	Logs        string
	Group       string
	Description string
	Status      string
	Handler     string
	Method      string
	FinishedAt  time.Time
}

type JobRepository struct{}

func (JobRepository) GetByRef(ref string) *Job {
	object := &Job{
		Ref: ref,
	}

	dbConn.Last(object)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (JobRepository) Get(ID uint) *Job {
	object := &Job{}

	dbConn.Last(object, ID)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (JobRepository) Create(object *Job) {
	dbConn.Create(object)
}

func (JobRepository) UpdateOrCreate(object *Job) {
	dbConn.Save(object)
}

func (JobRepository) Delete(ID uint) {
	dbConn.Delete(&Job{}, ID)
}

func init() {
	dbConn.AutoMigrate(&Job{})
}
