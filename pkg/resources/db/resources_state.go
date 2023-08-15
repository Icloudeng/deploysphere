package db

import (
	"encoding/base64"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/events/redis_events"
	"smatflow/platform-installer/pkg/terraform"

	"gorm.io/datatypes"
)

func ResourceStateCreate(ref string, job database.Job) *database.ResourcesState {
	rep := database.ResourcesStatesRepository{}

	resource_state := &database.ResourcesState{
		Ref:   ref,
		JobID: job.ID,
		Job:   job,
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

	current_state := resource_state.State.Data()

	for _, module := range childModules {
		address := module.Address
		for _, resource := range module.Resources {
			if resource.Name == resource_state.Ref {
				current_state[address] = resource
			}
		}
	}

	resource_state.State = datatypes.NewJSONType(current_state)

	repository.UpdateOrCreate(resource_state)
}

// =============== Redis Events Listener ============= //

func ResourceState_ListenResourceProviningCredentials(playload redis_events.ResourceRedisEventPayload) {
	rep := database.ResourcesStatesRepository{}
	resource_state := rep.GetByRef(playload.Reference)

	decodedBytes, err := base64.StdEncoding.DecodeString(playload.Payload)

	if resource_state == nil || err != nil {
		return
	}

	resource_state.Credentials = datatypes.JSON(decodedBytes)

	rep.UpdateOrCreate(resource_state)
}
