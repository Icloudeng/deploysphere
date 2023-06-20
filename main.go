package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/lib"
	"smatflow/platform-installer/lib/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	port := strconv.Itoa(*flag.Int("port", 8088, "Server port"))
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	api := r.Group("/", basicAuth)

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

	// Start server
	log.Println("Server running on PORT: ", port)
	log.Fatal(r.Run(":" + port))
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	username, password, hasAuth := c.Request.BasicAuth()

	if hasAuth && lib.LDAPExistBindUser(username, password) {
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
		return
	}
}
