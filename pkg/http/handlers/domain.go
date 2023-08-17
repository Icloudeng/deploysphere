package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/resources"
	"smatflow/platform-installer/pkg/resources/jobs"
	"smatflow/platform-installer/pkg/structs"
	"smatflow/platform-installer/pkg/terraform"

	"github.com/gin-gonic/gin"
)

type Domain struct {
	Ref    string                    `json:"ref" binding:"required,resourceref"`
	Domain *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
}

func CreateDomain(c *gin.Context) {
	var json Domain

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if Resource when post request
	if c.Request.Method == "POST" {
		_ovh := resources.GetOvhDomainZoneResource(json.Ref)

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
		Task: func(ctx context.Context) error {
			// Create or update resources
			resources.WriteOvhDomainZoneResource(json.Ref, json.Domain)

			// Terraform Apply changes
			err := terraform.Tf.Apply(true)
			if err == nil {
				events.BusEvent.Publish(events.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "succeeded",
					Details: "Ref: " + json.Ref,
					Logs:    "Domain Resource created",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func DeleteDomain(c *gin.Context) {
	var data ResourcesRef

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           data.Ref,
		PostBody:      data,
		Description:   "Domain Resource deletion",
		Handler:       c.Request.URL.String(),
		ResourceState: false, // Disable on resource deletion
		Task: func(ctx context.Context) error {
			// Remove resources
			resources.DeleteOvhDomainZoneResource(data.Ref)

			// Terraform Apply changes
			err := terraform.Tf.Apply(true)
			if err == nil {
				events.BusEvent.Publish(events.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: "Ref: " + data.Ref,
					Logs:    "Domain Resource deleted",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": data, "job": job})
}
