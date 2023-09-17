package entities

import (
	"smatflow/platform-installer/pkg/database"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ResourcesBackup struct {
	gorm.Model
	Resources datatypes.JSON
	State     datatypes.JSON
}

type ResourcesBackupRepository struct{}

func (ResourcesBackupRepository) Get() *ResourcesBackup {
	object := &ResourcesBackup{}

	database.Conn.Last(object)

	return object
}

func (ResourcesBackupRepository) Create(res *ResourcesBackup) {
	database.Conn.Create(res)
}

func init() {
	database.Conn.AutoMigrate(&ResourcesBackup{})
}
