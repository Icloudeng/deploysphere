package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"smatflow/platform-installer/pkg/env"
	frontproxy "smatflow/platform-installer/pkg/http/front_proxy"

	sentrygin "github.com/getsentry/sentry-go/gin"

	"smatflow/platform-installer/pkg/http/validators"
	"smatflow/platform-installer/pkg/http/ws"
	"smatflow/platform-installer/pkg/ldap"
	"smatflow/platform-installer/pkg/pubsub/subscribers"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	subscribers.EventSubscribers()
}

func main() {
	// Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.Config.SENTRY_DSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            true,
		AttachStacktrace: true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	app := gin.Default()

	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

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
