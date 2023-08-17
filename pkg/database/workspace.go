package database

import "gorm.io/gorm"

type Workspace struct {
	gorm.Model
	Name string `gorm:"uniqueIndex,unique"`
}

type WorkspaceRepository struct{}

func (r *WorkspaceRepository) Get() []Workspace {
	var objects []Workspace

	dbConn.Find(&objects)

	return objects
}

func (r *WorkspaceRepository) Create(res *Workspace) {
	dbConn.Create(res)
}

func init() {
	dbConn.AutoMigrate(&Workspace{})
}