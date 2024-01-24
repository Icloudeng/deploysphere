package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/http/validators"
	"github.com/icloudeng/platform-installer/internal/pubsub"
	"github.com/icloudeng/platform-installer/internal/resources/db"
	"github.com/icloudeng/platform-installer/internal/resources/jobs"
	"github.com/icloudeng/platform-installer/internal/resources/proxmox"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
	"github.com/icloudeng/platform-installer/internal/structs"

	"github.com/gin-gonic/gin"
)

type (
	vmBody struct {
		Ref      string                 `json:"ref" binding:"required,resourceref"`
		Vm       *structs.ProxmoxVmQemu `json:"vm" binding:"required,json"`
		Platform *structs.Platform      `json:"platform"`
	}

	vmHandler struct{}
)

var Vm vmHandler

func (vmHandler) CreateVm(c *gin.Context) {
	json := vmBody{
		Vm:       structs.NewProxmoxVmQemu(""),
		Platform: &structs.Platform{Metadata: map[string]interface{}{}},
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

	// Check if VM Id doesn't exist
	if json.Vm.Vmid != 0 {
		if exists := proxmox.VmQemuIDExists(json.Vm.Vmid); exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "VM ID already exists !",
			})
			return
		}
	}

	// If Target Node is set to auto,
	// then selected automatic node based on resourse Availability
	target_node := json.Vm.TargetNode
	if target_node == "auto" {
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
		Description:   "VM Resource creation",
		ResourceState: true,
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job *entities.Job) error {
			// Reselect Targe node
			if target_node == "auto" {
				if nodeStatus, err := proxmox.SelectNodeWithMostResources(); nodeStatus != nil && err == nil {
					json.Vm.TargetNode = nodeStatus.Node
					db.Jobs.JobUpdatePostBody(job, json)
				}
			}

			// Reset unmutable vm fields
			structs.ResetImmutableProxmoxVmQemu(&structs.ResetProxmoxVmQemuFields{
				Vm:       json.Vm,
				Platform: *json.Platform,
				Ref:      json.Ref,
				JobID:    job.ID,
			})
			// Create or update resources
			terraform.Resources.WriteProxmoxVmQemuResource(json.Ref, json.Vm)
			// Terraform Apply changes
			return terraform.Exec.Apply(true)
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func (vmHandler) DeleteVm(c *gin.Context) {
	var data resourcesRefUri

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           data.Ref,
		PostBody:      data,
		Description:   "VM Resource deletion",
		ResourceState: false, // Disable on resource deletion
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job *entities.Job) error {
			// Remove resources
			terraform.Resources.DeleteProxmoxVmQemuResource(data.Ref)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)
			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: fmt.Sprintf("Job ID: %d\nRef: %s", job.ID, data.Ref),
					Logs:    "VM Resource deleted",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": data, "job": job})
}
