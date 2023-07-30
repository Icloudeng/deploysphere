package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/events"
	"smatflow/platform-installer/pkg/files"
	"smatflow/platform-installer/pkg/queue"
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

	queue.Queue.QueueTask(func(ctx context.Context) error {
		proxyhost.CreateProxyHost(json)
		return nil
	})

	c.JSON(http.StatusOK, json)
}

func DeleteProxyHost(c *gin.Context) {
	var json structs.ProxyHostDomain

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queue.Queue.QueueTask(func(ctx context.Context) error {
		proxyhost.DeleteProxyHost(json.Domain)

		defer events.BusEvent.Publish(events.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
			Status:  "info",
			Details: "Domain: " + json.Domain,
			Logs:    "Proxy Host deleted",
		})
		return nil
	})

	c.JSON(http.StatusOK, json)
}
