package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"smatflow/platform-installer/pkg/filesystem"
	"smatflow/platform-installer/pkg/http/validators"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/resources/jobs"
	"smatflow/platform-installer/pkg/resources/terraform"
	"smatflow/platform-installer/pkg/structs"
)

// Store resources and apply
type Resources struct {
	Ref      string                    `json:"ref" binding:"required,resourceref"`
	Domain   *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
	Vm       *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
	Platform *structs.Platform         `json:"platform"`
}

// Delete resources and apply
type ResourcesRef struct {
	Ref string `uri:"ref" binding:"required,resourceref"`
}

const ResourceExistsError = `The resource reference already exists. If you plan to create a new resource, 
please use a different resource reference name or use PUT method to update resource.
`

func CreateResources(c *gin.Context) {
	json := Resources{
		Vm:       structs.NewProxmoxVmQemu(""),
		Platform: &structs.Platform{Metadata: &map[string]interface{}{}},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	task := jobs.ResourcesJob{
		Ref:           json.Ref,
		PostBody:      json,
		ResourceState: true,
		Description:   "Resources creation",
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context) error {
			// Reset unmutable vm fields
			structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform, json.Ref)
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

func DeleteResources(c *gin.Context) {
	var uri ResourcesRef
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
		Task: func(ctx context.Context) error {
			// Remove resources
			terraform.Resources.DeleteOvhDomainZoneResource(uri.Ref)
			terraform.Resources.DeleteProxmoxVmQemuResource(uri.Ref)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)
			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: "Ref: " + uri.Ref,
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
func GetResourcesState(c *gin.Context) {
	state := terraform.Exec.Show()

	if state == nil {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.JSON(http.StatusOK, state)
}

func GetResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Proxmox": terraform.Resources.GetProxmoxResource(),
		"Ovh":     terraform.Resources.GetOvhResource(),
	})
}

func GetResourcesByReference(c *gin.Context) {
	var uri ResourcesRef
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

func GetPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, filesystem.ReadProvisionerPlaforms())
}
