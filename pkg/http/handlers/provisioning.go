package handlers

import (
	"context"
	"net/http"
	"smatflow/platform-installer/pkg/resources/jobs"
	"smatflow/platform-installer/pkg/resources/provisioning"
	"smatflow/platform-installer/pkg/structs"
	"smatflow/platform-installer/pkg/validators"

	"github.com/gin-gonic/gin"
)

func CreatePlatformProvisioning(c *gin.Context) {
	body := &structs.Provisioning{
		Platform: &structs.Platform{
			Metadata: &map[string]interface{}{},
		},
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidatePlatformMetadata(c, *body.Platform) {
		return
	}

	// Validate ref and bind resource state platform values
	if !validators.ValidatePlatformProvisionAndBindResourceState(body) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cannot find resource linked to passed reference name!",
		})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           body.MachineIp,
		PostBody:      body,
		Description:   "Platform Provisioning",
		ResourceState: false,
		Task: func(ctx context.Context) error {
			provisioning.CreatePlatformProvisioning(*body)
			return nil
		},
	}

	jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, body)
}

func CreateConfigurationProvisioning(c *gin.Context) {
	body := &structs.Provisioning{
		Platform: &structs.Platform{
			Metadata: &map[string]interface{}{},
		},
	}

	if err := c.ShouldBindJSON(body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidateConfigurationMetadata(c, *body.Platform) {
		return
	}

	// Validate ref and bind resource state platform values
	if !validators.ValidatePlatformProvisionAndBindResourceState(body) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cannot find resource linked to passed reference name!",
		})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           body.MachineIp,
		PostBody:      body,
		Description:   "Platform Configuration Provisioning",
		ResourceState: false,
		Task: func(ctx context.Context) error {
			provisioning.CreateConfigurationProvisioning(*body)
			return nil
		},
	}

	jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, body)
}
