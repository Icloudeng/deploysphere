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
	Resources []*struct {
		Ref string `json:"ref" binding:"required,resourceref"`
	} `json:"resources" binding:"required"`
}

// DELETE resources bulk
func (resourcesHandler) DeleteResourcesBulk(c *gin.Context) {
	var bulk ResourcesRefBulk

	if err := c.ShouldBindUri(&bulk); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           "resources-bulk-deletion",
		PostBody:      bulk,
		Description:   "Resources bulk deletion",
		ResourceState: false, // Disable on resource deletion
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job entities.Job) error {
			for _, res := range bulk.Resources {
				// Remove resources
				terraform.Resources.DeleteOvhDomainZoneResource(res.Ref)
				terraform.Resources.DeleteProxmoxVmQemuResource(res.Ref)
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

	c.JSON(http.StatusOK, gin.H{"data": bulk, "job": job})
}
