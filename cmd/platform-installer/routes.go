package main

import (
	"net/http"
	"smatflow/platform-installer/pkg/env"
	frontproxy "smatflow/platform-installer/pkg/http/front_proxy"
	"smatflow/platform-installer/pkg/http/handlers"

	"github.com/gin-gonic/gin"
)

func BindLocalJobsRoutes(api *gin.RouterGroup) {
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

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
	api.POST("/provisioning", handlers.CreatePlatformProvisioning)

	api.POST("/provisioning/configuration", handlers.CreateConfigurationProvisioning)

	// Front Proxy
	if env.EnvConfig.FRONT_PROXY {
		api.Any("/ui/*all", frontproxy.Proxy)
	}
}
