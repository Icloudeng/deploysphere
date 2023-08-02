package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/jobs"
	"smatflow/platform-installer/pkg/env"
	"smatflow/platform-installer/pkg/events/subscribers"
	"smatflow/platform-installer/pkg/ldap"
	"smatflow/platform-installer/pkg/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	subscribers.EventSubscribers()
}

func main() {
	r := gin.Default()

	port := strconv.Itoa(*flag.Int("port", 8088, "Server port"))
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("resourceref", validators.ResourcesRefValidation)
	}

	// Routes
	api := r.Group("/", basicAuth)
	BindLocalJobsRoutes(api)

	apiJobs := r.Group("/jobs", basicAuth)
	jobs.BindDatabaseJobsRoutes(apiJobs)

	// Start server
	log.Println("Server running on PORT: ", port)
	log.Fatal(r.Run(":" + port))
}

func basicAuth(c *gin.Context) {
	// If LDAP_AUTH disable then authorize all request
	if !env.EnvConfig.LDAP_AUTH {
		c.Next()
		return
	}
	// Get the Basic Authentication credentials
	username, password, hasAuth := c.Request.BasicAuth()

	if hasAuth && ldap.LDAPExistBindUser(username, password) {
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
		return
	}
}
