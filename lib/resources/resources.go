package resources

import (
	"smatflow/platform-installer/lib/resources/ovh"
	"smatflow/platform-installer/lib/resources/proxmox"
	"smatflow/platform-installer/lib/structs"
)

func GetOvhResource() *ovh.Resource {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()

	return ovh_resource.GetResource()
}

/** OVH Domain creation resource functions **/
func CreateOrWriteOvhResource(ref string, domain *structs.DomainZoneRecord) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()
	// Add domain to the resource
	ovh_resource.GetResource().AddDomainZoneRerord(ref, domain)
}

func DeleteOvhResource(ref string) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseResourcesJSON()
	// Write resource data
	defer ovh_resource.WriteResources()

	// Add domain to the resource
	ovh_resource.GetResource().DeleteDomainZoneRerord(ref)
}

func GetProxmoxResource() *proxmox.Resource {
	// Working on ovh resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()

	return proxmox_resource.GetResource()
}

/** Promox resource functions **/
func CreateOrWriteProxmoxResource(ref string, pm *structs.ProxmoxVmQemu) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().AddProxmoxVmQemu(ref, pm)
}

func DeleteProxmoxResource(ref string) {
	// Working on Proxmox resource
	proxmox_resource := proxmox.ResourceJSONData{}
	proxmox_resource.ParseResourcesJSON()
	// Write resource data
	defer proxmox_resource.WriteResources()
	// Add domain to the resource
	proxmox_resource.GetResource().DeleteProxmoxVmQemu(ref)
}

// Initiliaze some function
func init() {
	ovh_resource := &ovh.ResourceJSONData{}
	proxmox_resource := &proxmox.ResourceJSONData{}

	// Create file id not exist yet
	ovh_resource.InitResourcesFiles()
	proxmox_resource.InitResourcesFiles()
}
