package terraform

import (
	"errors"

	"github.com/icloudeng/platform-installer/internal/resources/terraform/ovh"
	"github.com/icloudeng/platform-installer/internal/resources/utilities"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func (resources) GetOvhResource() *ovh.Resource {
	ovh_resource := ovh.NewResourceJSONData()
	return ovh_resource.GetResource()
}

func (r resources) GetOvhDomainZoneResource(ref string) *structs.DomainZoneRecord {
	return r.GetOvhResource().GetDomainZoneRerordStruct(ref)
}

/* OVH Domain creation resource functions */
func (resources) WriteOvhDomainZoneResource(ref string, domain *structs.DomainZoneRecord) {
	ovh_resource := ovh.NewResourceJSONData()
	// Add domain to the resource
	ovh_resource.GetResource().AddDomainZoneRerord(ref, domain)
	// ReWrite resources
	ovh_resource.WriteResources()
}

func (resources) DeleteOvhDomainZoneResource(ref string) {
	ovh_resource := ovh.NewResourceJSONData()
	// Add domain to the resource
	ovh_resource.GetResource().DeleteDomainZoneRerord(ref)
	// Write resource data
	ovh_resource.WriteResources()
}

func (r resources) GetOvhDomainZoneFromUrl(fqdn string) (string, *structs.DomainZoneRecord, error) {
	subdomain, rootdomain := utilities.Helpers.ExtractSubdomainAndRootDomain(fqdn)
	resouces := r.GetOvhResource().GetOVHDomainZoneRecord()

	for ref, value := range resouces {
		if value.Subdomain == subdomain && value.Zone == rootdomain {
			return ref, value, nil
		}
	}

	return "", nil, errors.New("unable to find the resource associated with this URL")
}

func init() {
	ovh_resource := &ovh.ResourceJSONData{}

	ovh_resource.InitResourcesFiles()
}
