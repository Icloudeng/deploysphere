package resources

import (
	"smatflow/platform-installer/resources/ovh"
	"smatflow/platform-installer/structs"
)

func CreateOrWriteOvhResource(ref string, domain structs.DomainZoneRecord) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseOVHresourcesJSON()
	// Write resource data
	defer ovh_resource.WriteOVHresources()

	// Add domain to the resource
	ovh_resource.GetResource().AddDomainZoneRerord(ref, &domain)

}

func DeleteOvhResource(ref string) {
	// Working on ovh resource
	ovh_resource := ovh.ResourceJSONData{}
	ovh_resource.ParseOVHresourcesJSON()
	// Write resource data
	defer ovh_resource.WriteOVHresources()

	// Add domain to the resource
	ovh_resource.GetResource().DeleteDomainZoneRerord(ref)

}

// Initiliaze some function
func init() {
	ovh_resource := &ovh.ResourceJSONData{}
	// Create file id not exist yet
	ovh_resource.InitResourcesFiles()
}
