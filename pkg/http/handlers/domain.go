package handlers

import (
	"context"
	"fmt"
	"net/http"
	"smatflow/platform-installer/pkg/database/entities"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/resources/jobs"
	"smatflow/platform-installer/pkg/resources/terraform"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

type (
	domainBody struct {
		Ref    string                    `json:"ref" binding:"required,resourceref"`
		Domain *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
	}

	domainHandler struct{}
)

var Domain domainHandler

func (domainHandler) CreateDomain(c *gin.Context) {
	var json domainBody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if Resource when post request
	if c.Request.Method == "POST" {
		_ovh := terraform.Resources.GetOvhDomainZoneResource(json.Ref)

		if _ovh != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": ResourceExistsError,
				"resource": map[string]interface{}{
					"domain": _ovh,
				},
			})
			return
		}
	}

	task := jobs.ResourcesJob{
		Ref:           json.Ref,
		PostBody:      json,
		Description:   "Domain Resource Creation",
		ResourceState: true,
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job entities.Job) error {
			// Create or update resources
			terraform.Resources.WriteOvhDomainZoneResource(json.Ref, json.Domain)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)
			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "succeeded",
					Details: fmt.Sprintf("Job ID: %d\nRef: %s", job.ID, json.Ref),
					Logs:    "Domain Resource created",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func (domainHandler) DeleteDomain(c *gin.Context) {
	var data resourcesRefUri

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           data.Ref,
		PostBody:      data,
		Description:   "Domain Resource deletion",
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		ResourceState: false, // Disable on resource deletion
		Task: func(ctx context.Context, job entities.Job) error {
			// Remove resources
			terraform.Resources.DeleteOvhDomainZoneResource(data.Ref)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)

			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: fmt.Sprintf("Job ID: %d\nRef: %s", job.ID, data.Ref),
					Logs:    "Domain Resource deleted",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": data, "job": job})
}
