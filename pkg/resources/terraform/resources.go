package terraform

import (
	"smatflow/platform-installer/pkg/resources/terraform/ovh"
	"smatflow/platform-installer/pkg/resources/terraform/proxmox"
	"smatflow/platform-installer/pkg/structs"
)

type resources struct{}

var Resources resources

//######################## OVH Resources ############################//

func (r resources) GetOvhResource() *ovh.Resource {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()

	return ovh_resource.GetResource()
}

func (r resources) GetOvhDomainZoneResource(ref string) *structs.DomainZoneRecord {
	return r.GetOvhResource().GetDomainZoneRerordStruct(ref)
}

/* OVH Domain creation resource functions */
func (r resources) WriteOvhDomainZoneResource(ref string, domain *structs.DomainZoneRecord) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()
	// Add domain to the resource
	ovh_resource.GetResource().AddDomainZoneRerord(ref, domain)
}

func (r resources) DeleteOvhDomainZoneResource(ref string) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()

	// Add domain to the resource
	ovh_resource.GetResource().DeleteDomainZoneRerord(ref)
}

//######################## Promox Resources ############################//

func (r resources) GetProxmoxResource() *proxmox.Resource {
	// Working on prpxmox vm resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()

	return proxmox_resource.GetResource()
}

func (r resources) GetProxmoxVmQemuResource(ref string) *structs.ProxmoxVmQemu {
	return r.GetProxmoxResource().GetProxmoxVmQemuStruct(ref)
}

func (r resources) WriteProxmoxVmQemuResource(ref string, pm *structs.ProxmoxVmQemu) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().AddProxmoxVmQemu(ref, pm)
}

func (r resources) DeleteProxmoxVmQemuResource(ref string) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().DeleteProxmoxVmQemu(ref)
}

// Store resources state

// ############## Init resources #####################
func init() {
	ovh_resource := &ovh.ResourceJSONData{}
	proxmox_resource := &proxmox.ResourceJSONData{}

	// Create file id not exist yet
	ovh_resource.InitResourcesFiles()
	proxmox_resource.InitResourcesFiles()
}
