package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/pubsub"
	"github.com/icloudeng/platform-installer/internal/resources/jobs"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
	"github.com/icloudeng/platform-installer/internal/structs"
)

type ResourcesRefBulk struct {
	Refs []string `json:"refs" binding:"required"`
}

// DELETE resources bulk
func (resourcesHandler) DeleteResourcesBulk(c *gin.Context) {
	var reqBody ResourcesRefBulk

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           "resources-bulk-deletion",
		PostBody:      reqBody,
		Description:   "Resources bulk deletion",
		ResourceState: false, // Disable on resource deletion
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job entities.Job) error {
			for _, ref := range reqBody.Refs {
				// Remove resources
				terraform.Resources.DeleteOvhDomainZoneResource(ref)
				terraform.Resources.DeleteOvhDomainZoneResource(fmt.Sprintf("mx-%s", ref))
				terraform.Resources.DeleteProxmoxVmQemuResource(ref)
			}

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)

			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: fmt.Sprintf("Job ID: %d", job.ID),
					Logs:    "Resources bulk deletion",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": reqBody, "job": job})
}
