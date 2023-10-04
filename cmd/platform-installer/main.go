package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/internal/env"
	frontproxy "smatflow/platform-installer/internal/http/front_proxy"

	"smatflow/platform-installer/internal/http/validators"
	"smatflow/platform-installer/internal/http/ws"
	"smatflow/platform-installer/internal/ldap"
	"smatflow/platform-installer/internal/pubsub/subscribers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	subscribers.EventSubscribers()
}

func main() {
	app := gin.Default()

	// Get port if passed arg
	port := strconv.Itoa(*flag.Int("port", 8088, "Server port"))
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	// Validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("resourceref", validators.ResourcesRefValidation)
	}

	// Routes
	api := app.Group("/", basicAuth)
	BindLocalJobsRoutes(api)

	// UI (Front Proxy)
	if env.Config.FRONT_PROXY {
		app.Group("/ui").Any("/*proxyPath", frontproxy.Proxy)
	}

	// Websocket bind
	api.Any("/ws", ws.ServeWs)

	// Start server
	log.Println("Server running on PORT: ", port)
	log.Fatalln(app.Run(":" + port))
}

func basicAuth(c *gin.Context) {
	// If LDAP_AUTH disable then authorize all request
	if !env.Config.LDAP_AUTH {
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
