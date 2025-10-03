package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
}

// CORS middleware handles CORS headers
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, config.AllowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)

			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			if len(config.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
			}

			// Handle preflight request
			if c.Request.Method == "OPTIONS" {
				c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
				c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))

				if config.MaxAge > 0 {
					c.Header("Access-Control-Max-Age", strconv.Itoa(int(config.MaxAge.Seconds())))
				}

				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

// isOriginAllowed checks if origin is in allowed list
func isOriginAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
		// Support wildcard subdomains
		if strings.HasPrefix(a, "*.") {
			domain := strings.TrimPrefix(a, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}