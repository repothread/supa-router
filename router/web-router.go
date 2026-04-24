package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetWebRouter(router *gin.Engine, indexPage []byte) {
	distDir := common.ResolveFrontendDistDir()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	if distDir != "" {
		router.Use(static.Serve("/", static.LocalFile(distDir, false)))
	}
	router.NoRoute(func(c *gin.Context) {
		c.Set(middleware.RouteTagKey, "web")
		if strings.HasPrefix(c.Request.RequestURI, "/v1") || strings.HasPrefix(c.Request.RequestURI, "/api") || strings.HasPrefix(c.Request.RequestURI, "/assets") {
			controller.RelayNotFound(c)
			return
		}
		if len(indexPage) > 0 {
			c.Header("Cache-Control", "no-cache")
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
			return
		}
		if distDir == "" {
			c.Header("Cache-Control", "no-store")
			c.Data(http.StatusServiceUnavailable, "text/html; charset=utf-8", []byte(common.FrontendUnavailablePageHTML))
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.File(filepath.Join(distDir, "index.html"))
	})
}
