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

	api.GET("/resources/:ref", handlers.GetResourcesByReference)

	api.GET("/resources", handlers.GetResources)

	api.GET("/resources/state", handlers.GetResourcesState)

	// Domain
	api.POST("/resources/domain", handlers.CreateDomain)

	api.PUT("/resources/domain", handlers.CreateDomain)

	api.DELETE("/resources/domain/:ref", handlers.DeleteDomain)

	//VM
	api.POST("/resources/vm", handlers.CreateVm)

	api.PUT("/resources/vm", handlers.CreateVm)

	api.DELETE("/resources/vm/:ref", handlers.DeleteVm)

	// Platforms
	api.GET("/platforms", handlers.GetPlatforms)

	// Proxy host
	api.POST("/proxy-host", handlers.CreateProxyHost)

	api.DELETE("/proxy-host", handlers.DeleteProxyHost)

	// The Provisioning
	api.POST("/provisioning", handlers.CreateProvisioning)
}
