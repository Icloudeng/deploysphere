package resources

import (
	"smatflow/platform-installer/pkg/resources/ovh"
	"smatflow/platform-installer/pkg/resources/proxmox"
	"smatflow/platform-installer/pkg/structs"
)

//######################## OVH Resources ############################//

func GetOvhResource() *ovh.Resource {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()

	return ovh_resource.GetResource()
}

func GetOvhDomainZoneResource(ref string) *structs.DomainZoneRecord {
	return GetOvhResource().GetDomainZoneRerordStruct(ref)
}

/* OVH Domain creation resource functions */
func CreateOrWriteOvhResource(ref string, domain *structs.DomainZoneRecord) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()
	// Add domain to the resource
	ovh_resource.GetResource().AddDomainZoneRerord(ref, domain)
}

func DeleteOvhDomainZoneResource(ref string) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()

	// Add domain to the resource
	ovh_resource.GetResource().DeleteDomainZoneRerord(ref)
}

//######################## Promox Resources ############################//

func GetProxmoxResource() *proxmox.Resource {
	// Working on prpxmox vm resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()

	return proxmox_resource.GetResource()
}

func GetProxmoxVmQemuResource(ref string) *structs.ProxmoxVmQemu {
	return GetProxmoxResource().GetProxmoxVmQemuStruct(ref)
}

func CreateOrWriteProxmoxResource(ref string, pm *structs.ProxmoxVmQemu) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().AddProxmoxVmQemu(ref, pm)
}

func DeleteProxmoxVmQemuResource(ref string) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().DeleteProxmoxVmQemu(ref)
}

// ############## Init resources #####################
func init() {
	ovh_resource := &ovh.ResourceJSONData{}
	proxmox_resource := &proxmox.ResourceJSONData{}

	// Create file id not exist yet
	ovh_resource.InitResourcesFiles()
	proxmox_resource.InitResourcesFiles()
}
