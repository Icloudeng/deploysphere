package proxmox

import (
	"encoding/json"
	"log"
	"path"

	"smatflow/platform-installer/pkg/filesystem"
	structs "smatflow/platform-installer/pkg/structs"
)

type ResourceJSONData struct {
	Resources []*Resource `json:"resource"`
}

type Resource struct {
	ProxmoxVmQemu []ProxmoxVmQemu `json:"proxmox_vm_qemu"`
}

type ProxmoxVmQemu map[string]*structs.ProxmoxVmQemu

func (j *ResourceJSONData) GetResource() *Resource {
	if len(j.Resources) > 0 {
		return j.Resources[0]
	}

	res := &Resource{}

	j.Resources = append(j.Resources, res)

	return res
}

func (j *ResourceJSONData) GetFile() string {
	return path.Join(filesystem.TerraformDir, "modules/proxmox/resource_auto.tf.json")
}

func (r *ResourceJSONData) InitResourcesFiles() {
	// Create file if not exist
	for _, file := range r.GetFiles() {
		filesystem.CreateIfNotExistsWithContent(file, "{}")
	}
}

func (r *ResourceJSONData) GetFiles() [1]string {
	return [...]string{r.GetFile()}
}

func (r *ResourceJSONData) ParseResourcesJSON() error {
	err := json.Unmarshal(filesystem.ReadFile(r.GetFile()), &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceJSONData) WriteResources() {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatalln(err)
	}
	var isEmpty = true

	for _, res := range r.Resources {
		if len(res.ProxmoxVmQemu) == 0 {
			break
		} else {
			isEmpty = false
		}
	}

	if isEmpty {
		filesystem.WriteInFile(r.GetFile(), "{}")
		return
	}
	filesystem.WriteInFile(r.GetFile(), string(data))
}

// Resource methods
func (r *Resource) GetProxmoxVmQemu() ProxmoxVmQemu {
	if len(r.ProxmoxVmQemu) > 0 {
		return r.ProxmoxVmQemu[0]
	}
	return nil
}

func (r *Resource) AddProxmoxVmQemu(ref string, vm *structs.ProxmoxVmQemu) {
	pm := r.GetProxmoxVmQemu()

	if pm == nil {
		r.ProxmoxVmQemu = append(r.ProxmoxVmQemu, ProxmoxVmQemu{ref: vm})
		return
	}

	pm[ref] = vm
}

func (r *Resource) DeleteProxmoxVmQemu(ref string) {
	ozr := r.GetProxmoxVmQemu()

	delete(ozr, ref)

	if len(ozr) == 0 {
		r.ProxmoxVmQemu = nil
	}
}

func (r *Resource) GetProxmoxVmQemuStruct(ref string) *structs.ProxmoxVmQemu {
	ozr := r.GetProxmoxVmQemu()

	found, exist := ozr[ref]
	if !exist {
		return nil
	}

	return found
}
