package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/resources"
	"smatflow/platform-installer/lib/structs"
	"smatflow/platform-installer/lib/terraform"

	"github.com/gin-gonic/gin"
)

type Vm struct {
	Ref      string                 `json:"ref" binding:"required,resourceref"`
	Vm       *structs.ProxmoxVmQemu `json:"vm" binding:"required,json"`
	Platform *structs.Platform      `json:"platform"`
}

func CreateVm(c *gin.Context) {
	json := Vm{
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

		if _vm != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": ResourceExistsError,
				"resource": map[string]interface{}{
					"vm": _vm,
				},
			})
			return
		}
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Reset unmutable vm fields
		structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform)
		// Create or update resources
		resources.CreateOrWriteProxmoxResource(json.Ref, json.Vm)
		// Terraform Apply changes
		defer terraform.Tf.Apply(true)
		return nil
	})

	c.JSON(http.StatusOK, json)
}

func DeleteVm(c *gin.Context) {
	var data ResourcesRef

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Remove resources
		resources.DeleteProxmoxVmQemuResource(data.Ref)
		// Terraform Apply changes
		defer terraform.Tf.Apply(true)
		return nil
	})

	c.JSON(http.StatusOK, data)
}

func validatePlatform(c *gin.Context, platform structs.Platform) bool {
	if len(platform.Name) > 0 {
		// Check if plaform has provisionner script
		if !files.ExistsProvisionerPlaformReadDir(platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return false
		}

		// Check if platform meta fields exist
		requriedFields := files.ReadPlatformMetadataFields()
		meta := structs.PlatformMetadataFields{}
		json.Unmarshal(requriedFields, &meta)

		metadata := *platform.Metadata

		if values, found := meta[platform.Name]; found {
			for _, val := range values {
				if _, exists := metadata[val]; !exists {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": fmt.Sprintf("platform (%s), Metadata field (%s) required", platform.Name, val)},
					)
					return false
				}
			}
		}
	}

	return true
}
