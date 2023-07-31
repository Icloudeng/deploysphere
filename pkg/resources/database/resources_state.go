package database

import (
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/terraform"

	"gorm.io/datatypes"
)

func StoreOrUpdateResourceState(ref string) {
	stateModule := terraform.Tf.Show()
	repository := database.ResourcesStatesRepository{}

	if stateModule == nil {
		return
	}

	childModules := stateModule.ChildModules

	current_res := repository.GetByRef(ref)
	current_state := current_res.State.Data()

	for _, module := range childModules {
		address := module.Address
		for _, resource := range module.Resources {
			if resource.Name == ref {
				current_state[address] = resource
			}
		}
	}

	current_res.Ref = ref
	current_res.State = datatypes.NewJSONType(current_state)

	repository.UpdateOrCreate(current_res)
}
