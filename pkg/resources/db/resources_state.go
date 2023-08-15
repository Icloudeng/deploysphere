package db

import (
	"smatflow/platform-installer/pkg/database"
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
