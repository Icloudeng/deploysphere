package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Ref      string `gorm:"index"`
	PostBody datatypes.JSON
	Logs     string
	Group    string
	Running  bool
	Success  bool
}

type JobRepository struct{}

func (r *JobRepository) GetByRef(ref string) *Job {
	res := &Job{
		Ref: ref,
	}

	db.Last(res)

	return res
}

func (r *JobRepository) Create(res *Job) {
	db.Create(res)
}

func (r *JobRepository) UpdateOrCreate(res *Job) {
	db.Save(res)
}

func (r *JobRepository) Delete(ID uint) {
	db.Delete(&Job{}, ID)
}

func init() {
	db.AutoMigrate(&Job{})
}
