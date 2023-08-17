package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ResourcesBackup struct {
	gorm.Model
	Resources datatypes.JSON
	State     datatypes.JSON
}

type ResourcesBackupRepository struct{}

func (r *ResourcesBackupRepository) Get(ID uint) *ResourcesBackup {
	object := &ResourcesBackup{}

	dbConn.Last(object)

	return object
}

func (r *ResourcesBackupRepository) Create(res *ResourcesBackup) {
	dbConn.Create(res)
}

func init() {
	dbConn.AutoMigrate(&ResourcesBackup{})
}
