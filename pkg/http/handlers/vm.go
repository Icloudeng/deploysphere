package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/resources"
	"smatflow/platform-installer/pkg/resources/jobs"
	"smatflow/platform-installer/pkg/structs"
	"smatflow/platform-installer/pkg/terraform"
	"smatflow/platform-installer/pkg/validators"

	"github.com/gin-gonic/gin"
)

type Vm struct {
	Ref      string                 `json:"ref" binding:"required,resourceref"`
	Vm       *structs.ProxmoxVmQemu `json:"vm" binding:"required,json"`
	Platform *structs.Platform      `json:"platform"`
}

func CreateVm(c *gin.Context) {
	json := Vm{
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

	task := jobs.ResourcesJob{
		Ref:           json.Ref,
		PostBody:      json,
		Description:   "VM Resource creation",
		ResourceState: true,
		Task: func(ctx context.Context) error {
			// Reset unmutable vm fields
			structs.ResetUnmutableProxmoxVmQemu(json.Vm, *json.Platform, json.Ref)
			// Create or update resources
			resources.WriteProxmoxVmQemuResource(json.Ref, json.Vm)
			// Terraform Apply changes
			return terraform.Tf.Apply(true)
		},
	}

	jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, json)
}

func DeleteVm(c *gin.Context) {
	var data ResourcesRef

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           data.Ref,
		PostBody:      data,
		Description:   "VM Resource deletion",
		ResourceState: false, // Disable on resource deletion
		Task: func(ctx context.Context) error {
			// Remove resources
			resources.DeleteProxmoxVmQemuResource(data.Ref)

			// Terraform Apply changes
			err := terraform.Tf.Apply(true)
			if err == nil {
				events.BusEvent.Publish(events.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: "Ref: " + data.Ref,
					Logs:    "VM Resource deleted",
				})
			}

			return err
		},
	}

	jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, data)
}
