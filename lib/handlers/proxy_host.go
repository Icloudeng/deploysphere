package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/files"
	proxyhost "smatflow/platform-installer/lib/resources/proxy_host"
	"smatflow/platform-installer/lib/structs"

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

	lib.Queue.QueueTask(func(ctx context.Context) error {
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

	lib.Queue.QueueTask(func(ctx context.Context) error {
		proxyhost.DeleteProxyHost(json.Domain)
		return nil
	})

	c.JSON(http.StatusOK, json)
}
