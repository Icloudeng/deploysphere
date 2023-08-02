package resources

import (
	"smatflow/platform-installer/pkg/terraform"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformResourceState struct {
	Module string
}

func (f TerraformResourceState) getResourceByRef(module *tfjson.StateModule, ref string) *tfjson.StateResource {
	if strings.Contains(module.Address, f.Module) {
		for _, resource := range module.Resources {
			if resource.Name == ref {
				return resource
			}
		}
	}

	return nil
}

func (f TerraformResourceState) GetResourceState(ref string) *tfjson.StateResource {
	stateModule := terraform.Tf.Show()

	if stateModule == nil {
		return nil
	}

	if resource := f.getResourceByRef(stateModule, ref); resource != nil {
		return resource
	}

	for _, module := range stateModule.ChildModules {
		if resource := f.getResourceByRef(module, ref); resource != nil {
			return resource
		}
	}

	return nil
}
