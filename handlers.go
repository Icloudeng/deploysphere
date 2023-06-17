package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"smatflow/platform-installer/files"
	"smatflow/platform-installer/resources"
	"smatflow/platform-installer/structs"
)

// Handler top lever
type Handler struct{}

var Handlers = Handler{}

// Store provistion and apply

type Provision struct {
	Ref      string                    `json:"ref" binding:"required,ascii,lowercase"`
	Domain   *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
	Vm       *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
	Platform *structs.Platform         `json:"platform"`
}

func (s *Handler) provision(c *gin.Context) {
	json := Provision{
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
		if err := Queue.QueueTask(func(ctx context.Context) error {
			// Reset unmutable vm fields
			structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform)
			// Create or update resources
			resources.CreateOrWriteOvhResource(json.Ref, json.Domain)
			resources.CreateOrWriteProxmoxResource(json.Ref, json.Vm)

			// Terraform Apply changes
			defer Tf.Apply()
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	c.AsciiJSON(http.StatusOK, json)
}

// Delete provistion and apply
type ProvisionRef struct {
	Ref string `uri:"ref" binding:"required,alpha,lowercase"`
}

func (s *Handler) deleteProvision(c *gin.Context) {
	var data ProvisionRef
	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	go func() {
		if err := Queue.QueueTask(func(ctx context.Context) error {
			// Remove resources
			resources.DeleteOvhResource(data.Ref)
			resources.DeleteProxmoxResource(data.Ref)

			// Terraform Apply changes
			defer Tf.Apply()
			return nil
		}); err != nil {
			panic(err)
		}
	}()

	c.AsciiJSON(http.StatusOK, data)
}
