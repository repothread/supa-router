package middleware

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Cache-Control", "New-Api-User", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type", "Auth-Version", "X-New-Api-Version"}

	allowedOrigins := getCORSAllowedOrigins()
	if len(allowedOrigins) == 0 {
		config.AllowAllOrigins = true
		config.AllowCredentials = false
	} else {
		config.AllowOrigins = allowedOrigins
		config.AllowCredentials = true
	}
	return cors.New(config)
}

func getCORSAllowedOrigins() []string {
	raw := strings.TrimSpace(common.GetEnvOrDefaultString("CORS_ALLOW_ORIGINS", ""))
	if raw == "" {
		if frontendBaseURL := system_setting.GetFrontendBaseURL(); frontendBaseURL != "" {
			return []string{frontendBaseURL}
		}
		return nil
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		origins = append(origins, part)
	}
	return origins
}

func PoweredBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-New-Api-Version", common.Version)
		c.Next()
	}
}
