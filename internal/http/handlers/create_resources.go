package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/http/validators"
	"github.com/icloudeng/platform-installer/internal/resources/db"
	"github.com/icloudeng/platform-installer/internal/resources/jobs"
	"github.com/icloudeng/platform-installer/internal/resources/proxmox"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
	util "github.com/icloudeng/platform-installer/internal/resources/utilities"
	"github.com/icloudeng/platform-installer/internal/structs"
)

type domainStruct struct {
	ref    string
	record *structs.DomainZoneRecord
}

type domains []domainStruct

func createResourceJob(ctx *gin.Context, json *resourcesBody) *entities.Job {
	domains := make(domains, 0)

	// If domain key doesn't exist in metadata platform
	// then auto fill with the passed domain resource
	metadata := json.Platform.Metadata
	domainStr := fmt.Sprintf("%s.%s", json.Domain.Subdomain, json.Domain.Zone)

	_, domain_exists := metadata["domain"]
	if !domain_exists {
		json.Platform.Metadata["domain"] = domainStr
	} else {
		domainStr = metadata["domain"].(string)
	}

	/** Domain resource **/
	domains = append(domains, domainStruct{
		ref:    json.Ref,
		record: json.Domain,
	})

	/** Mx Domain resource **/
	mxDomain := autoComposeMxDomain(domainStr, json)
	if mxDomain != nil {
		json.MxDomainValue = mxDomain

		domains = append(domains, domainStruct{
			ref:    fmt.Sprintf("mx-%s", json.Ref),
			record: mxDomain,
		})
	}

	if json.Platform != nil && mxDomain != nil {
		json.Platform.Metadata["mx_domain"] = util.Helpers.ConcatenateSubdomain(
			mxDomain.Subdomain,
			mxDomain.Zone,
		)
	}

	/** Subdomains resources **/
	if json.Subdomains == nil {
		json.Subdomains = make([]string, 0)
	}

	subdomains := util.Helpers.RemoveDuplicates(json.Subdomains)
	for _, subdomain := range subdomains {
		record := &structs.DomainZoneRecord{
			Zone:      json.Domain.Zone,
			Subdomain: fmt.Sprintf("%s.%s", subdomain, json.Domain.Subdomain),
			Fieldtype: json.Domain.Fieldtype,
			Ttl:       json.Domain.Ttl,
			Target:    json.Domain.Target,
		}

		domains = append(domains, domainStruct{
			ref:    fmt.Sprintf("%s-subdomain-%s", subdomain, json.Ref),
			record: record,
		})

		// Add subdomain to platform metadata
		if json.Platform != nil {
			metadataKey := fmt.Sprintf("%s_subdomain", subdomain)
			json.Platform.Metadata[metadataKey] = util.Helpers.ConcatenateSubdomain(
				record.Subdomain,
				record.Zone,
			)
		}
	}

	/*
	 * Validate platform metadata
	 */
	if !validators.ValidatePlatformMetadata(ctx, *json.Platform) {
		return nil
	}

	// Skip apply when both resources exist
	var shouldSkipApply = false
	eVm := terraform.Resources.GetProxmoxVmQemuResource(json.Ref)
	eDomain := terraform.Resources.GetOvhDomainZoneResource(json.Ref)

	if eVm != nil && eDomain != nil {
		shouldSkipApply = true
	}

	// Failure when resources exists on POST request
	if ctx.Request.Method == "POST" {
		if eVm != nil || eDomain != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": ResourceExistsError,
				"resource": map[string]interface{}{
					"vm":     eVm,
					"domain": eDomain,
				},
			})

			return nil
		}
	}

	/** Check if VM Id doesn't exist
	if json.Vm.Vmid != 0 {
		if exists := proxmox.VmQemuIDExists(json.Vm.Vmid); exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "VM ID already exists !",
			})
			return nil
		}
	}
	**/

	/*
	 * If Target Node is set to auto,
	 * then selected automatic node based on resource Availability
	 */
	target_node := json.Vm.TargetNode
	if target_node == "auto" {
		nodeStatus, err := proxmox.SelectNodeWithMostResources()
		if err != nil || nodeStatus == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "No enough proxmox resources",
				"message": err.Error(),
			})
			return nil
		}

		json.Vm.TargetNode = nodeStatus.Node
	}

	json.Vm.Description = fmt.Sprintf("https://%s", domainStr)

	/** Create job task **/
	task := jobs.ResourcesJob{
		Ref:           json.Ref,
		PostBody:      json,
		ResourceState: true,
		Description:   "Resources creation",
		Handler:       ctx.Request.URL.String(),
		Method:        ctx.Request.Method,
		Task: func(ctx context.Context, job *entities.Job) error {
			if shouldSkipApply {
				return nil
			}

			// Reselect Targe node
			if target_node == "auto" {
				if nodeStatus, err := proxmox.SelectNodeWithMostResources(); nodeStatus != nil && err == nil {
					json.Vm.TargetNode = nodeStatus.Node
					db.Jobs.JobUpdatePostBody(job, json)
				}
			}

			// Reset Immutable vm fields
			structs.ResetImmutableProxmoxVmQemu(&structs.ResetProxmoxVmQemuFields{
				Vm:       json.Vm,
				Platform: *json.Platform,
				Ref:      json.Ref,
				JobID:    job.ID,
			})

			/** Domain resources **/
			for _, v := range domains {
				terraform.Resources.WriteOvhDomainZoneResource(v.ref, v.record)
			}

			/** Proxmox resources **/
			terraform.Resources.WriteProxmoxVmQemuResource(json.Ref, json.Vm)

			// Terraform Apply changes
			return terraform.Exec.Apply(true)
		},
	}

	return jobs.ResourcesJobTask(task)

}

func autoComposeMxDomain(resourceDomain string, json *resourcesBody) *structs.DomainZoneRecord {
	if json.MxDomain != nil && len(*json.MxDomain) > 1 {
		mx_value := *json.MxDomain
		if mx_value == "auto" || mx_value == resourceDomain {
			subdomain, rootDomain := util.Helpers.ExtractSubdomainAndRootDomain(resourceDomain)
			return &structs.DomainZoneRecord{
				Zone:      rootDomain,
				Subdomain: util.Helpers.RemoveFirstSegment(subdomain),
				Fieldtype: "MX",
				Ttl:       3600,
				Target:    fmt.Sprintf("1 %s.", resourceDomain),
			}
		} else {
			subdomain, rootDomain := util.Helpers.ExtractSubdomainAndRootDomain(mx_value)
			return &structs.DomainZoneRecord{
				Zone:      rootDomain,
				Subdomain: subdomain,
				Fieldtype: "MX",
				Ttl:       3600,
				Target:    fmt.Sprintf("1 %s.", resourceDomain),
			}
		}
	}

	return nil
}
