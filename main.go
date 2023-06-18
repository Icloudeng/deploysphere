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

	port := flag.Int("port", 8088, "Server port")
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	api := r.Group("/", basicAuth)
	api.POST("/provision", Handlers.provision)
	api.DELETE("/provision/:ref", Handlers.deleteProvision)

	log.Println("Server running on PORT: ", port)
	log.Fatal(r.Run(":" + strconv.Itoa(*port)))
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
