package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/resources/provisioning"
	"smatflow/platform-installer/lib/structs"

	"github.com/gin-gonic/gin"
)

func CreateProvisioning(c *gin.Context) {
	var json structs.Provisioning

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validatePlatform(c, *json.Platform) {
		return
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		provisioning.CreateProvisioning(json)
		return nil
	})

	c.JSON(http.StatusOK, json)
}
