package terraform

import (
	"errors"
	"fmt"
	"strings"

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
	ovh_resource.GetResource().DeleteDomainZoneRecord(ref)
	// Write resource data
	ovh_resource.WriteResources()
}

func (r resources) GetOvhDomainZoneFromUrl(fqdn string) (string, *structs.DomainZoneRecord, error) {
	subdomain, rootDomain := utilities.Helpers.ExtractSubdomainAndRootDomain(fqdn)
	resources := r.GetOvhResource().GetOVHDomainZoneRecord()

	for ref, value := range resources {
		if value.Subdomain == subdomain && value.Zone == rootDomain {
			return ref, value, nil
		}
	}

	return "", nil, errors.New("unable to find the resource associated with this URL")
}

func (r resources) GetOvhDomainZoneRefsLinkedSubdomains(ref string) []string {
	var refs []string
	resources := r.GetOvhResource().GetOVHDomainZoneRecord()

	for resourceRef := range resources {
		suffix := fmt.Sprintf("-subdomain-%s", ref)
		if strings.HasSuffix(resourceRef, suffix) {
			refs = append(refs, resourceRef)
		}
	}

	return refs
}

func init() {
	ovh_resource := &ovh.ResourceJSONData{}

	ovh_resource.InitResourcesFiles()
}
