package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/resources"
	"smatflow/platform-installer/lib/structs"
	"smatflow/platform-installer/lib/terraform"
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
		Vm:       structs.NewProxmoxVmQemu(),
		Platform: &structs.Platform{Metadata: &map[string]interface{}{}},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validatePlatform(c, *json.Platform) {
		return
	}

	// Check if Resource when post request
	if c.Request.Method == "POST" {
		_vm := resources.GetProxmoxVmQemuResource(json.Ref)
		_domain := resources.GetOvhDomainZoneResource(json.Ref)

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

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Reset unmutable vm fields
		structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform)
		// Create or update resources
		resources.CreateOrWriteOvhResource(json.Ref, json.Domain)
		resources.CreateOrWriteProxmoxResource(json.Ref, json.Vm)

		// Terraform Apply changes
		defer terraform.Tf.Apply(true)
		return nil
	})

	c.JSON(http.StatusOK, json)
}

func DeleteResources(c *gin.Context) {
	var uri ResourcesRef
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Remove resources
		resources.DeleteOvhDomainZoneResource(uri.Ref)
		resources.DeleteProxmoxVmQemuResource(uri.Ref)

		// Clean up resource event publish
		// lib.BusEvent.Publish(lib.RESOURCES_CLEANUP_EVENT, domain)

		// Terraform Apply changes
		defer lib.BusEvent.Publish(lib.NOTIFIER_RESOURCES_EVENT, structs.Notifier{
			Status:  "info",
			Details: "Ref: " + uri.Ref,
			Logs:    "Resources deleted",
		})
		defer terraform.Tf.Apply(true)

		return nil
	})

	c.JSON(http.StatusOK, uri)
}

// Get resources state from terraform
func GetResourcesState(c *gin.Context) {
	state := terraform.Tf.Show()

	if state == nil {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.JSON(http.StatusOK, state)
}

func GetResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Proxmox": resources.GetProxmoxResource(),
		"Ovh":     resources.GetOvhResource(),
	})
}

func GetResourcesByReference(c *gin.Context) {
	var uri ResourcesRef
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Vm":     resources.GetProxmoxVmQemuResource(uri.Ref),
		"Domain": resources.GetOvhDomainZoneResource(uri.Ref),
		"Ref":    uri.Ref,
	})
}

func GetPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, files.ReadProvisionerPlaforms())
}
