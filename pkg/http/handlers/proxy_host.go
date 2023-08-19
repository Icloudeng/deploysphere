package handlers

import (
	"context"
	"fmt"
	"net/http"
	"smatflow/platform-installer/pkg/database"
	"smatflow/platform-installer/pkg/filesystem"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/resources/jobs"
	proxyhost "smatflow/platform-installer/pkg/resources/proxy_host"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

type proxyHostHandler struct{}

var ProxyHost proxyHostHandler

func (proxyHostHandler) CreateProxyHost(c *gin.Context) {
	var json structs.ProxyHost

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if plaform exists
	if !filesystem.ExistsProvisionerPlaformReadDir(json.Platform) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           json.Domain,
		PostBody:      json,
		Description:   "Proxy Host Creation",
		ResourceState: false,
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job database.Job) error {
			proxyhost.CreateProxyHost(json)
			return nil
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

func (proxyHostHandler) DeleteProxyHost(c *gin.Context) {
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
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job database.Job) error {
			proxyhost.DeleteProxyHost(json.Domain)

			pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
				Status:  "info",
				Details: fmt.Sprintf("Job Id: %d \n Domain: %s", job.ID, json.Domain),
				Logs:    "Proxy Host deleted",
			})
			return nil
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}
