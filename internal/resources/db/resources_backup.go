package db

import (
	"encoding/json"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"

	"gorm.io/datatypes"
)

type (
	resourcesBackup struct{}
)

var ResourcesBackup resourcesBackup

func (resourcesBackup) CreateNewResourcesBackup() {
	// Terraform state
	stateModule := terraform.Exec.Show()
	var state interface{} = map[string]interface{}{}
	if stateModule != nil {
		state = stateModule
	}

	// Local state
	ovh := terraform.Resources.GetOvhResource()
	proxmox := terraform.Resources.GetProxmoxResource()

	// Fill Reourses type here
	resources_list := map[string]interface{}{
		"Ovh":     ovh,
		"Proxmox": proxmox,
	}

	// Encode JSON
	resources_encoded, _ := json.Marshal(resources_list)
	state_encoded, _ := json.Marshal(state)

	// Store Database
	repository := entities.ResourcesBackupRepository{}
	repository.Create(&entities.ResourcesBackup{
		Resources: datatypes.JSON(resources_encoded),
		State:     datatypes.JSON(state_encoded),
	})
}
