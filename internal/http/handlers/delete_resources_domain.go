package handlers

import (
	"fmt"

	"github.com/icloudeng/platform-installer/internal/resources/terraform"
)

func deleteResourcesDomain(ref string) {
	terraform.Resources.DeleteOvhDomainZoneResource(ref)
	terraform.Resources.DeleteOvhDomainZoneResource(fmt.Sprintf("mx-%s", ref))

	linkedRefs := terraform.Resources.GetOvhDomainZoneRefsLinkedSubdomains(ref)
	fmt.Println(linkedRefs)
	for _, linkedRef := range linkedRefs {
		terraform.Resources.DeleteOvhDomainZoneResource(linkedRef)
	}
}
