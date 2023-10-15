package terraform

import (
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type (
	ResourceState struct {
		Module string
	}

	resources struct{}
)

var Resources resources

func (f ResourceState) getResourceByRef(module *tfjson.StateModule, ref string) *tfjson.StateResource {
	if strings.Contains(module.Address, f.Module) {
		for _, resource := range module.Resources {
			if resource.Name == ref {
				return resource
			}
		}
	}

	return nil
}

func (f ResourceState) GetResourceState(ref string) *tfjson.StateResource {
	stateModule := Exec.Show()

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
