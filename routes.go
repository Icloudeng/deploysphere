package main

import (
	"smatflow/platform-installer/lib/handlers"

	"github.com/gin-gonic/gin"
)

func BindLocalJobsRoutes(api *gin.RouterGroup) {
	// Resources
	api.POST("/resources", handlers.CreateResources)

	api.PUT("/resources", handlers.CreateResources)

	api.DELETE("/resources/:ref", handlers.DeleteResources)

	api.GET("/resources", handlers.GetResources)

	api.GET("/resources/state", handlers.GetResourcesState)

	// Domain
	api.POST("/resource/domain", handlers.CreateDomain)

	api.PUT("/resource/domain", handlers.CreateDomain)

	api.DELETE("/resource/domain/:ref", handlers.DeleteDomain)

	//VM
	api.POST("/resource/vm", handlers.CreateVm)

	api.PUT("/resource/vm", handlers.CreateVm)

	api.DELETE("/resource/vm/:ref", handlers.DeleteVm)

	// Platforms
	api.GET("/platforms", handlers.GetPlatforms)

	// Proxy host
	api.POST("/proxy-host", handlers.CreateProxyHost)

	api.DELETE("/proxy-host", handlers.DeleteProxyHost)

	// The Provisioning
	api.POST("/provisioning", handlers.CreateProvisioning)
}
