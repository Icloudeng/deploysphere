package handlers

import (
	"net/http"
	"smatflow/platform-installer/pkg/database/entities"

	"github.com/gin-gonic/gin"
)

type (
	resourcesStateHandler struct{}
)

var ResourceState resourcesStateHandler

// Get resources state from terraform
func (resourcesStateHandler) GetByRef(c *gin.Context) {
	var uri resourcesRefUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	// Resource state from db
	repository := entities.ResourcesStateRepository{}
	resources := repository.GetByRef(uri.Ref)

	c.JSON(http.StatusOK, gin.H{
		"data": resources,
	})
}
