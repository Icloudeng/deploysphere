package entities

import (
	"database/sql"
	"smatflow/platform-installer/internal/database"

	"gorm.io/gorm"
)

type Workspace struct {
	gorm.Model
	Name   string       `gorm:"uniqueIndex,unique"`
	Active sql.NullBool `gorm:"default:true"`
}

type WorkspaceRepository struct{}

func (WorkspaceRepository) Get() []Workspace {
	var objects []Workspace

	database.Conn.Find(&objects)

	return objects
}

func (WorkspaceRepository) Create(res *Workspace) {
	database.Conn.Create(res)
}

func init() {
	database.Conn.AutoMigrate(&Workspace{})
}
