package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/resources"
	"smatflow/platform-installer/lib/structs"
	"smatflow/platform-installer/lib/terraform"

	"github.com/gin-gonic/gin"
)

type Vm struct {
	Ref      string                 `json:"ref" binding:"required,ascii,lowercase"`
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
	if len(json.Platform.Name) > 0 {
		if !files.ExistsProvisionerPlaformReadDir(json.Platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return
		}
	}

	go func() {
		if err := lib.Queue.QueueTask(func(ctx context.Context) error {
			// Reset unmutable vm fields
			structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform)
			// Create or update resources
			resources.CreateOrWriteProxmoxResource(json.Ref, json.Vm)
			// Terraform Apply changes
			defer terraform.Tf.Apply()
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	c.JSON(http.StatusOK, json)
}

func DeleteVm(c *gin.Context) {
	var data ResourcesRef

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	go func() {
		if err := lib.Queue.QueueTask(func(ctx context.Context) error {
			// Remove resources
			resources.DeleteProxmoxResource(data.Ref)
			// Terraform Apply changes
			defer terraform.Tf.Apply()
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	c.JSON(http.StatusOK, data)
}
