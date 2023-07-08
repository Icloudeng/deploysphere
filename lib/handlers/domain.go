package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/resources"
	"smatflow/platform-installer/lib/structs"
	"smatflow/platform-installer/lib/terraform"

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

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Create or update resources
		resources.CreateOrWriteOvhResource(json.Ref, json.Domain)
		// Terraform Apply changes
		defer terraform.Tf.Apply(true)
		return nil
	})

	c.JSON(http.StatusOK, json)
}

func DeleteDomain(c *gin.Context) {
	var data ResourcesRef

	if err := c.ShouldBindUri(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err})
		return
	}

	lib.Queue.QueueTask(func(ctx context.Context) error {
		// Remove resources
		resources.DeleteOvhDomainZoneResource(data.Ref)
		// Terraform Apply changes
		defer terraform.Tf.Apply(true)
		return nil
	})

	c.JSON(http.StatusOK, data)
}
