package entities

import (
	"time"

	"github.com/icloudeng/platform-installer/internal/database"

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
	var object Job

	database.Conn.Where(&Job{
		Ref: ref,
	}).Last(&object)

	if object.ID == 0 {
		return nil
	}

	return &object
}

func (JobRepository) Get(ID uint) *Job {
	object := &Job{}

	database.Conn.Last(object, ID)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (JobRepository) Create(object *Job) {
	database.Conn.Create(object)
}

func (JobRepository) UpdateOrCreate(object *Job) {
	database.Conn.Save(object)
}

func (JobRepository) Delete(ID uint) {
	database.Conn.Delete(&Job{}, ID)
}

func init() {
	database.Conn.AutoMigrate(&Job{})
}
