package frontproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/icloudeng/platform-installer/internal/env"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	remote, err := url.Parse(env.Config.FRONT_URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	//Define the director func
	//This is a good place to log, for example
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = "/ui" + c.Param("proxyPath")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
