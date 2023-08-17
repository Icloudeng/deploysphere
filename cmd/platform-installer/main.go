package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/pkg/env"
	frontproxy "smatflow/platform-installer/pkg/http/front_proxy"
	"smatflow/platform-installer/pkg/http/validators"
	"smatflow/platform-installer/pkg/http/ws"
	"smatflow/platform-installer/pkg/ldap"
	"smatflow/platform-installer/pkg/pubsub/subscribers"

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

	// UI (Front Proxy)
	if env.Config.FRONT_PROXY {
		r.Group("/ui").Any("/*proxyPath", frontproxy.Proxy)
	}

	// Websocket bind
	api.Any("/ws", ws.ServeWs)

	// Start server
	log.Println("Server running on PORT: ", port)
	log.Fatalln(r.Run(":" + port))
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
