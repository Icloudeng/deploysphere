package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
)

type (
	autoProvisioning struct {
		Type              string `json:"type" binding:"required"`
		PlatformUrl       string `json:"platform_url" binding:"required,fqdn"`
		PlatformConfigUrl string `json:"platform_config_url" binding:"required,fqdn"`
	}
)

func (provisioningHandler) CreateAutoConfigurationProvisioning(c *gin.Context) {
	body := &autoProvisioning{}

	if err := c.ShouldBindJSON(body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get platform reference from url
	platformRef, _, err := terraform.Resources.GetOvhDomainZoneFromUrl(body.PlatformUrl)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"url":   body.PlatformUrl,
		})
		return
	}

	// Get platform Config reference from url
	platformConfigRef, _, err := terraform.Resources.GetOvhDomainZoneFromUrl(body.PlatformConfigUrl)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"url":   body.PlatformUrl,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": body, "job": []string{platformRef, platformConfigRef}})
}
