package db

import (
	"encoding/json"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/resources"
	"smatflow/platform-installer/pkg/terraform"

	"gorm.io/datatypes"
)

func CreateNewResourcesBackup() {
	// Terraform state
	stateModule := terraform.Tf.Show()
	var state interface{} = map[string]interface{}{}
	if stateModule != nil {
		state = stateModule
	}

	// Local state
	ovh := resources.GetOvhResource()
	proxmox := resources.GetProxmoxResource()

	// Fill Reourses type here
	resources_list := map[string]interface{}{
		"Ovh":     ovh,
		"Proxmox": proxmox,
	}

	// Encode JSON
	resources_encoded, _ := json.Marshal(resources_list)
	state_encoded, _ := json.Marshal(state)

	// Store Database
	repository := database.ResourcesBackupRepository{}
	repository.Create(&database.ResourcesBackup{
		Resources: datatypes.JSON(resources_encoded),
		State:     datatypes.JSON(state_encoded),
	})
}
