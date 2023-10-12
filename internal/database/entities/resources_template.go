package entities

import (
	"github.com/icloudeng/platform-installer/internal/database"
	"github.com/icloudeng/platform-installer/internal/structs"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ResourcesTemplate struct {
	gorm.Model
	PlatformName string `gorm:"index,unique"`
	Domain       datatypes.JSONType[*structs.DomainZoneRecord]
	Vm           datatypes.JSONType[*structs.ProxmoxVmQemu]
	Platform     datatypes.JSONType[*structs.Platform]
}

type ResourcesTemplateRepository struct{}

func (ResourcesTemplateRepository) GetByPlatform(platformName string) *ResourcesTemplate {
	var object ResourcesTemplate

	database.Conn.Where(&ResourcesTemplate{
		PlatformName: platformName,
	}).Last(&object)

	return &object
}

func (r ResourcesTemplateRepository) UpdateOrCreate(res *ResourcesTemplate) {
	if object := r.GetByPlatform(res.PlatformName); object != nil {
		res.ID = object.ID
	}

	database.Conn.Save(res)
}

func (ResourcesTemplateRepository) Delete(ID uint) {
	database.Conn.Delete(&ResourcesTemplate{}, ID)
}

func init() {
	database.Conn.AutoMigrate(&ResourcesTemplate{})
}
