package proxmox

import (
	"encoding/json"
	"log"
	"path"

	files "smatflow/platform-installer/lib/files"
	structs "smatflow/platform-installer/lib/structs"
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
	return path.Join(files.TerraformDir, "modules/proxmox/resource_auto.tf.json")
}

func (r *ResourceJSONData) InitResourcesFiles() {
	// Create file if not exist
	for _, file := range r.GetFiles() {
		files.CreateIfNotExistsWithContent(file, "{}")
	}
}

func (r *ResourceJSONData) GetFiles() [1]string {
	return [...]string{r.GetFile()}
}

func (r *ResourceJSONData) ParseResourcesJSON() error {
	err := json.Unmarshal(files.ReadFile(r.GetFile()), &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceJSONData) WriteResources() {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
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
		files.WriteInFile(r.GetFile(), "{}")
		return
	}
	files.WriteInFile(r.GetFile(), string(data))
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
