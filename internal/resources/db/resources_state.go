package db

import (
	"encoding/base64"
	"encoding/json"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/pubsub"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"

	"gorm.io/datatypes"
)

type (
	resourceState struct{}
)

var ResourceState resourceState

func (resourceState) ResourceStateCreate(ref string, job entities.Job) *entities.ResourcesState {
	rep := entities.ResourcesStateRepository{}

	// Use credentials and state of the last ref object
	var credentials datatypes.JSON
	var state datatypes.JSON
	if last_rs := rep.GetByRef(ref); last_rs != nil {
		credentials = last_rs.Credentials
		state = last_rs.State
	}

	resource_state := &entities.ResourcesState{
		Ref:         ref,
		JobID:       job.ID,
		Job:         job,
		Credentials: credentials,
		State:       state,
	}

	rep.Create(resource_state)

	return resource_state
}

func (resourceState) ResourceStatePutTerraformState(params *entities.ResourcesState) {
	stateModule := terraform.Exec.Show()
	repository := entities.ResourcesStateRepository{}

	// Refresh Object
	resource_state := repository.Get(params.ID)

	if stateModule == nil {
		return
	}

	childModules := stateModule.ChildModules

	state := map[string]interface{}{}

	for _, module := range childModules {
		for _, resource := range module.Resources {
			if resource.Name == resource_state.Ref {
				state[resource.Type] = resource
			}
		}
	}

	state_encoded, _ := json.Marshal(state)

	resource_state.State = datatypes.JSON(state_encoded)

	repository.UpdateOrCreate(resource_state)
}

// =============== Redis Events Listener ============= //

func (resourceState) ResourceState_ListenResourceProviningCredentials(playload pubsub.NetworkEventPayload) {
	rep := entities.ResourcesStateRepository{}
	resource_state := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if resource_state == nil || err != nil {
		return
	}

	resource_state.Credentials = datatypes.JSON(decodedBytes)

	rep.UpdateOrCreate(resource_state)
}
