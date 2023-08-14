package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type HistoryPostBody map[string]interface{}

type ResourcesHistory struct {
	gorm.Model
	Ref              string `gorm:"index,unique"`
	StateId          uint
	State            datatypes.JSONType[StateType]
	PostBody         datatypes.JSONType[HistoryPostBody]
	ProvisioningLogs string
}

type ResourcesHistoryRepository struct{}

func (r *ResourcesHistoryRepository) GetByRef(ref string) *ResourcesHistory {
	res := &ResourcesHistory{
		Ref: ref,
	}

	db.First(res)

	return res
}

func (r *ResourcesHistoryRepository) Create(res *ResourcesHistory) {
	db.Create(res)
}

func (r *ResourcesHistoryRepository) UpdateOrCreate(res *ResourcesHistory) {
	db.Save(res)
}

func (r *ResourcesHistoryRepository) Delete(ID uint) {
	db.Delete(&ResourcesHistory{}, ID)
}

func init() {
	db.AutoMigrate(&ResourcesHistory{})
}
