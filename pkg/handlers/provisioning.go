package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/queue"
	"smatflow/platform-installer/pkg/resources/provisioning"
	"smatflow/platform-installer/pkg/structs"
	"smatflow/platform-installer/pkg/validators"

	"github.com/gin-gonic/gin"
)

func CreatePlatformProvisioning(c *gin.Context) {
	json := structs.Provisioning{
		Platform: &structs.Platform{
			Metadata: &map[string]interface{}{},
		},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidatePlatformMetadata(c, *json.Platform) {
		return
	}

	queue.Queue.QueueTask(func(ctx context.Context) error {
		provisioning.CreatePlatformProvisioning(json)
		return nil
	})

	c.JSON(http.StatusOK, json)
}

func CreateConfigurationProvisioning(c *gin.Context) {
	json := structs.Provisioning{
		Platform: &structs.Platform{
			Metadata: &map[string]interface{}{},
		},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidateConfigurationMetadata(c, *json.Platform) {
		return
	}

	queue.Queue.QueueTask(func(ctx context.Context) error {
		provisioning.CreateConfigurationProvisioning(json)
		return nil
	})

	c.JSON(http.StatusOK, json)
}
