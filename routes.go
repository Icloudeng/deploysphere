package main

import (
	"smatflow/platform-installer/lib/handlers"

	"github.com/gin-gonic/gin"
)

func BindLocalJobsRoutes(api *gin.RouterGroup) {
	// Resources
	api.GET("/resources", handlers.GetResources)

	api.POST("/resources", handlers.CreateResources)

	api.DELETE("/resources/:ref", handlers.DeleteResources)

	api.GET("/resources/state", handlers.GetResourcesState)

	// Platforms
	api.GET("/platforms", handlers.GetPlatforms)

	// Domain
	api.POST("/domain", handlers.CreateDomain)

	api.DELETE("/domain/:ref", handlers.DeleteDomain)

	//VM
	api.POST("/vm", handlers.CreateVm)

	api.DELETE("/vm/:ref", handlers.DeleteVm)

	// Proxy host
	api.POST("/proxy-host", handlers.CreateProxyHost)

	api.DELETE("/proxy-host", handlers.DeleteProxyHost)

	// The Provisioning
	api.POST("/provisioning", handlers.CreateProvisioning)
}
