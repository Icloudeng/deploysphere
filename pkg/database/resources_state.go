package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"

	tfjson "github.com/hashicorp/terraform-json"
)

type StateType map[string]*tfjson.StateResource

type ResourcesState struct {
	gorm.Model
	Ref         string `gorm:"index"`
	State       datatypes.JSONType[StateType]
	Credentials datatypes.JSON
	JobID       uint
	Job         Job `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ResourcesStatesRepository struct{}

func (r *ResourcesStatesRepository) GetByRef(ref string) *ResourcesState {
	res := &ResourcesState{
		Ref: ref,
	}

	dbConn.Last(res)

	if res.ID == 0 {
		return nil
	}

	return res
}

func (r *ResourcesStatesRepository) Get(ID uint) *ResourcesState {
	object := &ResourcesState{}

	dbConn.Last(object, ID)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (r *ResourcesStatesRepository) Create(res *ResourcesState) {
	dbConn.Create(res)
}

func (r *ResourcesStatesRepository) UpdateOrCreate(res *ResourcesState) {
	dbConn.Save(res)
}

func (r *ResourcesStatesRepository) Delete(ID uint) {
	dbConn.Delete(&ResourcesState{}, ID)
}

func init() {
	dbConn.AutoMigrate(&ResourcesState{})
}
