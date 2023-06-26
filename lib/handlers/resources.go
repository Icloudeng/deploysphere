package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/resources"
	"smatflow/platform-installer/lib/resources/ovh"
	"smatflow/platform-installer/lib/resources/proxmox"
	"smatflow/platform-installer/lib/structs"
	"smatflow/platform-installer/lib/terraform"
)

// Store resources and apply
type Resources struct {
	Ref      string                    `json:"ref" binding:"required,ascii"`
	Domain   *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
	Vm       *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
	Platform *structs.Platform         `json:"platform"`
}

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

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Reset unmutable vm fields
		structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform)
		// Create or update resources
		resources.CreateOrWriteOvhResource(json.Ref, json.Domain)
		resources.CreateOrWriteProxmoxResource(json.Ref, json.Vm)

		// Terraform Apply changes
		defer terraform.Tf.Apply()
		return nil
	})

	c.JSON(http.StatusOK, json)
}

// Delete resources and apply
type ResourcesRef struct {
	Ref string `uri:"ref" binding:"required,ascii"`
}

func DeleteResources(c *gin.Context) {
	var data ResourcesRef
	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Remove resources
		resources.DeleteOvhResource(data.Ref)
		resources.DeleteProxmoxResource(data.Ref)

		// Clean up resource event publish
		// lib.BusEvent.Publish(lib.RESOURCES_CLEANUP_EVENT)

		// Terraform Apply changes
		defer terraform.Tf.Apply()

		return nil
	})

	c.JSON(http.StatusOK, data)
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
	res := struct {
		Proxmox *proxmox.Resource
		Domain  *ovh.Resource
	}{
		Proxmox: resources.GetProxmoxResource(),
		Domain:  resources.GetOvhResource(),
	}

	c.JSON(http.StatusOK, res)
}

func GetPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, files.ReadProvisionerPlaforms())
}
