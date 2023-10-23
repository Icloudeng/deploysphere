package terraform

import (
	"errors"

	"github.com/icloudeng/platform-installer/internal/resources/terraform/proxmox"
	"github.com/icloudeng/platform-installer/internal/resources/utilities"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func (r resources) GetProxmoxResource() *proxmox.Resource {
	// Working on prpxmox vm resource
	proxmox_resource := proxmox.NewResourceJSONData()
	return proxmox_resource.GetResource()
}

func (r resources) GetProxmoxVmQemuResource(ref string) *structs.ProxmoxVmQemu {
	return r.GetProxmoxResource().GetProxmoxVmQemuStruct(ref)
}

func (r resources) WriteProxmoxVmQemuResource(ref string, pm *structs.ProxmoxVmQemu) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.NewResourceJSONData()

	// Add domain to the resource
	proxmox_resource.GetResource().AddProxmoxVmQemu(ref, pm)

	// Write resource data
	proxmox_resource.WriteResources()
}

func (r resources) DeleteProxmoxVmQemuResource(ref string) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.NewResourceJSONData()

	// Add domain to the resource
	proxmox_resource.GetResource().DeleteProxmoxVmQemu(ref)

	// Write resource data
	proxmox_resource.WriteResources()
}

func (r resources) GetPlatformNameByReference(ref string) (string, error) {
	resource := r.GetProxmoxVmQemuResource(ref)
	if resource == nil {
		return "", errors.New("unable to platform respond the passed reference")
	}

	provisioner := resource.Provisioner[0].(map[string]interface{})
	provision := provisioner["local-exec"].([]interface{})
	command := provision[0].(*structs.PmLocalExec)

	keyValueMap := utilities.Helpers.ExtractCommandKeyValuePairs(command.Command)

	if platform, ok := keyValueMap["platform"]; !ok {
		return "", errors.New("unable to platform respond the passed reference")
	} else {
		return platform, nil
	}
}

func init() {
	proxmox_resource := &proxmox.ResourceJSONData{}

	proxmox_resource.InitResourcesFiles()
}
