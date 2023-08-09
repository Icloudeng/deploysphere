package frontproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"smatflow/platform-installer/pkg/env"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	remote, err := url.Parse(env.EnvConfig.FRONT_URL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	//Define the director func
	//This is a good place to log, for example
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
