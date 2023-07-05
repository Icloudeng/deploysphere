package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/jobs"
	"smatflow/platform-installer/lib"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	r := gin.Default()

	port := strconv.Itoa(*flag.Int("port", 8088, "Server port"))
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("resourceref", lib.ResourcesRefValidation)
	}

	api := r.Group("/local", basicAuth)
	apiJobs := r.Group("/jobs", basicAuth)

	// Routes
	BindLocalJobsRoutes(api)
	jobs.BindDatabaseJobsRoutes(apiJobs)

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
