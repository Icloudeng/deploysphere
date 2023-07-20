package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/queue"
	"smatflow/platform-installer/pkg/resources/provisioning"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

func CreateProvisioning(c *gin.Context) {
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
	if !validatePlatform(c, *json.Platform) {
		return
	}

	queue.Queue.QueueTask(func(ctx context.Context) error {
		provisioning.CreateProvisioning(json)
		return nil
	})

	c.JSON(http.StatusOK, json)
}
