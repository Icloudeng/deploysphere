package terraform

import (
	"github.com/icloudeng/platform-installer/internal/resources/terraform/proxmox"
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

func init() {
	proxmox_resource := &proxmox.ResourceJSONData{}

	proxmox_resource.InitResourcesFiles()
}
