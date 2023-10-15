package main

import (
	"net/http"

	"github.com/icloudeng/platform-installer/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func BindLocalJobsRoutes(api *gin.RouterGroup) {
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Resources
	api.POST("/resources", handlers.Resources.CreateResources)

	api.POST("/resources/templates", handlers.Resources.CreateResourcesFromTemplate)

	api.PUT("/resources", handlers.Resources.CreateResources)

	api.DELETE("/resources/:ref", handlers.Resources.DeleteResources)

	api.GET("/resources/:ref", handlers.Resources.GetResourcesByReference)

	api.GET("/resources", handlers.Resources.GetResources)

	api.GET("/resources/state", handlers.Resources.GetTerraformResourcesState)

	// Domain
	api.POST("/resources/domain", handlers.Domain.CreateDomain)

	api.PUT("/resources/domain", handlers.Domain.CreateDomain)

	api.DELETE("/resources/domain/:ref", handlers.Domain.DeleteDomain)

	//VM
	api.POST("/resources/vm", handlers.Vm.CreateVm)

	api.PUT("/resources/vm", handlers.Vm.CreateVm)

	api.DELETE("/resources/vm/:ref", handlers.Vm.DeleteVm)

	// Platforms
	api.GET("/platforms", handlers.Resources.GetPlatforms)

	// Proxy host
	api.POST("/proxy-host", handlers.ProxyHost.CreateProxyHost)

	api.DELETE("/proxy-host", handlers.ProxyHost.DeleteProxyHost)

	// The Provisioning
	api.POST("/provisioning", handlers.Provisioning.CreatePlatformProvisioning)

	api.POST("/provisioning/configuration", handlers.Provisioning.CreateConfigurationProvisioning)

	api.POST("/provisioning/auto-configuration", handlers.Provisioning.CreateAutoConfigurationProvisioning)

	// Jobs
	api.GET("/jobs/:id", handlers.Jobs.GetJobsByID)

	// Resource State
	api.GET("/resources-state/:ref", handlers.ResourceState.GetByRef)

	// Client
	api.POST("/clients", handlers.Client.CreateClient)

	// Platform Template
	api.POST("/platforms/templates", handlers.PlatformsTemplates.CreateResourcesTemplate)

	api.GET("/platforms/templates/:platform", handlers.PlatformsTemplates.GetResourcesTemplate)
}
