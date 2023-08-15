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

	db.Last(res)

	return res
}

func (r *ResourcesStatesRepository) Create(res *ResourcesState) {
	db.Create(res)
}

func (r *ResourcesStatesRepository) UpdateOrCreate(res *ResourcesState) {
	db.Save(res)
}

func (r *ResourcesStatesRepository) Delete(ID uint) {
	db.Delete(&ResourcesState{}, ID)
}

func init() {
	db.AutoMigrate(&ResourcesState{})
}
