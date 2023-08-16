package db

import (
	"encoding/base64"
	"encoding/json"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/terraform"

	"gorm.io/datatypes"
)

func ResourceStateCreate(ref string, job database.Job) *database.ResourcesState {
	rep := database.ResourcesStatesRepository{}

	// Use credentials of the last ref object
	var credentials datatypes.JSON
	if last_rs := rep.GetByRef(ref); last_rs != nil {
		credentials = last_rs.Credentials
	}

	resource_state := &database.ResourcesState{
		Ref:         ref,
		JobID:       job.ID,
		Job:         job,
		Credentials: credentials,
	}

	rep.Create(resource_state)

	return resource_state
}

func ResourceStatePutTerraformState(resource_state *database.ResourcesState) {
	stateModule := terraform.Tf.Show()
	repository := database.ResourcesStatesRepository{}

	// Refresh Object
	resource_state = repository.Get(resource_state.ID)

	if stateModule == nil {
		return
	}

	childModules := stateModule.ChildModules

	state := map[string]interface{}{}

	for _, module := range childModules {
		address := module.Address
		for _, resource := range module.Resources {
			if resource.Name == resource_state.Ref {
				state[address] = resource
			}
		}
	}

	state_encoded, _ := json.Marshal(state)

	resource_state.State = datatypes.JSON(state_encoded)

	repository.UpdateOrCreate(resource_state)
}

// =============== Redis Events Listener ============= //

func ResourceState_ListenResourceProviningCredentials(playload events.NetworkEventPayload) {
	rep := database.ResourcesStatesRepository{}
	resource_state := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if resource_state == nil || err != nil {
		return
	}

	resource_state.Credentials = datatypes.JSON(decodedBytes)

	rep.UpdateOrCreate(resource_state)
}
