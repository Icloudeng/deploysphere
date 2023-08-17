package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/files"
	"smatflow/platform-installer/pkg/resources/jobs"
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

func CreateProxyHost(c *gin.Context) {
	var json structs.ProxyHost

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if plaform exists
	if !files.ExistsProvisionerPlaformReadDir(json.Platform) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           json.Domain,
		PostBody:      json,
		Description:   "Proxy Host Creation",
		ResourceState: false,
		Handler:       c.Request.URL.String(),
		Task: func(ctx context.Context) error {
			proxyhost.CreateProxyHost(json)
			return nil
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func DeleteProxyHost(c *gin.Context) {
	var json structs.ProxyHostDomain

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           json.Domain,
		PostBody:      json,
		Description:   "Proxy Host Deletion",
		ResourceState: false,
		Handler:       c.Request.URL.String(),
		Task: func(ctx context.Context) error {
			proxyhost.DeleteProxyHost(json.Domain)

			events.BusEvent.Publish(events.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
				Status:  "info",
				Details: "Domain: " + json.Domain,
				Logs:    "Proxy Host deleted",
			})
			return nil
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}
