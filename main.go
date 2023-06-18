package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

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

	api.POST("/resources", Handlers.createResources)

	api.DELETE("/resources/:ref", Handlers.deleteResources)

	api.GET("/resources/state", Handlers.getResourcesState)

	// Start server
	log.Println("Server running on PORT: ", port)
	log.Fatal(r.Run(":" + port))
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	username, password, hasAuth := c.Request.BasicAuth()

	if hasAuth && LDAPExistBindUser(username, password) {
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
		return
	}
}
