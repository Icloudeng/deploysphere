package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"smatflow/platform-installer/internal/database/entities"
	"smatflow/platform-installer/internal/filesystem"
	"smatflow/platform-installer/internal/http/validators"
	"smatflow/platform-installer/internal/pubsub"
	"smatflow/platform-installer/internal/resources/jobs"
	"smatflow/platform-installer/internal/resources/proxmox"
	"smatflow/platform-installer/internal/resources/terraform"
	"smatflow/platform-installer/internal/structs"
)

// Store resources and apply
type (
	resourcesBody struct {
		Ref      string                    `json:"ref" binding:"required,resourceref"`
		Domain   *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
		Vm       *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
		Platform *structs.Platform         `json:"platform"`
	}

	resourcesRefUri struct {
		Ref string `uri:"ref" binding:"required,resourceref"`
	}

	resourcesHandler struct{}
)

var Resources resourcesHandler

const ResourceExistsError = `The resource reference already exists. If you plan to create a new resource, 
please use a different resource reference name or use PUT method to update resource.
`

func (resourcesHandler) CreateResources(c *gin.Context) {
	json := resourcesBody{
		Vm:       structs.NewProxmoxVmQemu(""),
		Platform: &structs.Platform{Metadata: map[string]interface{}{}},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If domain key doesn't exist in metadata platform
	// then auto fill with the passed domain resource
	metadata := json.Platform.Metadata
	_, domain_exists := metadata["domain"]
	if !domain_exists {
		domain := fmt.Sprintf("%s.%s", json.Domain.Subdomain, json.Domain.Zone)
		json.Platform.Metadata["domain"] = domain
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidatePlatformMetadata(c, *json.Platform) {
		return
	}

	// Check if Resource when post request
	if c.Request.Method == "POST" {
		_vm := terraform.Resources.GetProxmoxVmQemuResource(json.Ref)
		_domain := terraform.Resources.GetOvhDomainZoneResource(json.Ref)

		if _vm != nil || _domain != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": ResourceExistsError,
				"resource": map[string]interface{}{
					"vm":     _vm,
					"domain": _domain,
				},
			})
			return
		}
	}

	// Check if VM Id doesn't exist
	// if json.Vm.Vmid != 0 {
	// 	if exists := proxmox.VmQemuIDExists(json.Vm.Vmid); exists {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 			"error": "VM ID already exists !",
	// 		})

	// 		return
	// 	}
	// }

	// If Target Node is set to auto,
	// then selected automatic node based on resourse Availability
	if json.Vm.TargetNode == "auto" {
		nodeStatus, err := proxmox.SelectNodeWithMostResources()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "No enough proxmox resources",
			})
			return
		}

		json.Vm.TargetNode = nodeStatus.Node
	}

	task := jobs.ResourcesJob{
		Ref:           json.Ref,
		PostBody:      json,
		ResourceState: true,
		Description:   "Resources creation",
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job entities.Job) error {
			// Reset unmutable vm fields
			structs.ResetUnmutableProxmoxVmQemu(structs.ResetProxmoxVmQemuFields{
				Vm:       json.Vm,
				Platform: *json.Platform,
				Ref:      json.Ref,
				JobID:    job.ID,
			})
			// Create or update resources
			terraform.Resources.WriteOvhDomainZoneResource(json.Ref, json.Domain)
			terraform.Resources.WriteProxmoxVmQemuResource(json.Ref, json.Vm)

			// Terraform Apply changes
			return terraform.Exec.Apply(true)
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func (resourcesHandler) DeleteResources(c *gin.Context) {
	var uri resourcesRefUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           uri.Ref,
		PostBody:      uri,
		Description:   "Resources deletion",
		ResourceState: false, // Disable on resource deletion
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job entities.Job) error {
			// Remove resources
			terraform.Resources.DeleteOvhDomainZoneResource(uri.Ref)
			terraform.Resources.DeleteProxmoxVmQemuResource(uri.Ref)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)

			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: fmt.Sprintf("Job ID: %d\nRef: %s", job.ID, uri.Ref),
					Logs:    "Resources deleted",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": uri, "job": job})
}

// Get resources state from terraform
func (resourcesHandler) GetTerraformResourcesState(c *gin.Context) {
	state := terraform.Exec.Show()

	if state == nil {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (resourcesHandler) GetResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Proxmox": terraform.Resources.GetProxmoxResource(),
		"Ovh":     terraform.Resources.GetOvhResource(),
	})
}

func (resourcesHandler) GetResourcesByReference(c *gin.Context) {
	var uri resourcesRefUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Vm":     terraform.Resources.GetProxmoxVmQemuResource(uri.Ref),
		"Domain": terraform.Resources.GetOvhDomainZoneResource(uri.Ref),
		"Ref":    uri.Ref,
	})
}

func (resourcesHandler) GetPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, filesystem.ReadProvisionerPlaforms())
}
