package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ============================================================================
// CORS 설정 타입
// ============================================================================

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins     []string      `json:"allow_origins"`
	AllowMethods     []string      `json:"allow_methods"`
	AllowHeaders     []string      `json:"allow_headers"`
	ExposeHeaders    []string      `json:"expose_headers"`
	AllowCredentials bool          `json:"allow_credentials"`
	AllowWildcard    bool          `json:"allow_wildcard"`
	MaxAge           time.Duration `json:"max_age"`
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// DevelopmentCORSConfig returns development CORS configuration
func DevelopmentCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		AllowWildcard:   true,
		MaxAge:          24 * time.Hour,
	}
}

// ProductionCORSConfig returns production CORS configuration
func ProductionCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{
			"https://app.example.com",
			"https://www.example.com",
			"https://admin.example.com",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Link",
			"X-Total-Count",
			"X-Request-ID",
		},
		AllowCredentials: true,
		AllowWildcard:   false,
		MaxAge:          3600 * time.Second,
	}
}

// ============================================================================
// Custom CORS Middleware
// ============================================================================

// CustomCORS creates a custom CORS middleware
func CustomCORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, config.AllowOrigins) {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", origin)

			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight request
			if c.Request.Method == "OPTIONS" {
				c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
				c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))

				if len(config.ExposeHeaders) > 0 {
					c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
				}

				c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", int(config.MaxAge.Seconds())))
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			// Set expose headers for actual request
			if len(config.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
			}
		} else if !config.AllowWildcard {
			// Origin not allowed
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if origin is in allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// Support wildcard subdomains
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}

// ============================================================================
// Dynamic CORS Configuration
// ============================================================================

// DynamicCORS creates a CORS middleware with dynamic configuration
func DynamicCORS() gin.HandlerFunc {
	// This could load from database or config file
	allowedOrigins := make(map[string]bool)
	allowedOrigins["http://localhost:3000"] = true
	allowedOrigins["https://app.example.com"] = true

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Dynamic check (could query database)
		if allowedOrigins[origin] || isDynamicallyAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")

			if c.Request.Method == "OPTIONS" {
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
				c.Header("Access-Control-Max-Age", "86400")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

// isDynamicallyAllowed could check database or external service
func isDynamicallyAllowed(origin string) bool {
	// Example: Check if origin is from a trusted partner
	// This could query a database or cache
	trustedPartners := []string{
		"https://partner1.com",
		"https://partner2.com",
	}

	for _, partner := range trustedPartners {
		if origin == partner {
			return true
		}
	}

	return false
}

// ============================================================================
// Per-Route CORS Configuration
// ============================================================================

// PerRouteCORS allows different CORS settings per route
func PerRouteCORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if allowed == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				break
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ============================================================================
// CORS with Security Headers
// ============================================================================

// SecureCORS adds CORS with additional security headers
func SecureCORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// CORS headers
		if isOriginAllowed(origin, config.AllowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)

			if config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Handle preflight
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", int(config.MaxAge.Seconds())))
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ============================================================================
// API Handlers
// ============================================================================

// Public API (allows all origins)
func publicHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a public API endpoint",
		"cors":    "Allows all origins",
		"time":    time.Now().Unix(),
	})
}

// Private API (restricted origins)
func privateHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a private API endpoint",
		"cors":    "Restricted to specific origins",
		"user":    "authenticated_user",
		"time":    time.Now().Unix(),
	})
}

// Upload handler (demonstrates CORS with file upload)
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file type
	if !strings.HasSuffix(file.Filename, ".jpg") && !strings.HasSuffix(file.Filename, ".png") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG and PNG files allowed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// WebSocket endpoint (special CORS handling)
func wsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "WebSocket endpoint",
		"note":    "WebSocket connections have special CORS requirements",
		"url":     "ws://localhost:8080/ws",
	})
}

// ============================================================================
// Router Setup
// ============================================================================

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Get environment
	env := getEnv("APP_ENV", "development")

	// Configure CORS based on environment
	var corsConfig CORSConfig
	switch env {
	case "production":
		corsConfig = ProductionCORSConfig()
		log.Println("🔒 Using production CORS configuration")
	case "development":
		corsConfig = DevelopmentCORSConfig()
		log.Println("🔧 Using development CORS configuration")
	default:
		corsConfig = DefaultCORSConfig()
		log.Println("📦 Using default CORS configuration")
	}

	// Method 1: Using custom CORS middleware
	router.Use(CustomCORS(corsConfig))

	// Method 2: Using gin-contrib/cors package (alternative)
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     corsConfig.AllowOrigins,
	// 	AllowMethods:     corsConfig.AllowMethods,
	// 	AllowHeaders:     corsConfig.AllowHeaders,
	// 	ExposeHeaders:    corsConfig.ExposeHeaders,
	// 	AllowCredentials: corsConfig.AllowCredentials,
	// 	MaxAge:           corsConfig.MaxAge,
	// }))

	// Public routes (allow all origins)
	public := router.Group("/public")
	{
		public.GET("/info", publicHandler)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
	}

	// API routes with default CORS
	api := router.Group("/api")
	{
		api.GET("/data", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []string{"item1", "item2", "item3"},
			})
		})

		api.POST("/echo", func(c *gin.Context) {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, body)
		})
	}

	// Private routes with strict CORS
	private := router.Group("/private")
	private.Use(PerRouteCORS([]string{
		"https://app.example.com",
		"http://localhost:3000",
	}))
	{
		private.GET("/user", privateHandler)
		private.POST("/upload", uploadHandler)
	}

	// Admin routes with very strict CORS
	admin := router.Group("/admin")
	admin.Use(SecureCORS(CORSConfig{
		AllowOrigins:     []string{"https://admin.example.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:          time.Hour,
	}))
	{
		admin.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"users": []string{"admin1", "admin2"}})
		})

		admin.GET("/settings", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"settings": "admin settings"})
		})
	}

	// WebSocket endpoint (special handling)
	router.GET("/ws", wsHandler)

	// CORS configuration endpoint
	router.GET("/cors/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"environment":       env,
			"allowed_origins":   corsConfig.AllowOrigins,
			"allowed_methods":   corsConfig.AllowMethods,
			"allowed_headers":   corsConfig.AllowHeaders,
			"expose_headers":    corsConfig.ExposeHeaders,
			"allow_credentials": corsConfig.AllowCredentials,
			"max_age":          corsConfig.MaxAge.String(),
		})
	})

	// Test endpoints for different methods
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "GET"})
	})

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "POST"})
	})

	router.PUT("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "PUT"})
	})

	router.DELETE("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "DELETE"})
	})

	router.PATCH("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"method": "PATCH"})
	})

	return router
}

// ============================================================================
// Utility Functions
// ============================================================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ============================================================================
// Main
// ============================================================================

func main() {
	router := setupRouter()

	log.Println("🚀 Server starting on :8080")
	log.Println("🌐 CORS configuration loaded")
	log.Println("📝 Environment:", getEnv("APP_ENV", "development"))
	log.Println("")
	log.Println("Test CORS with:")
	log.Println("  curl -H 'Origin: http://localhost:3000' http://localhost:8080/api/data")
	log.Println("  curl -X OPTIONS -H 'Origin: http://localhost:3000' http://localhost:8080/api/data")
	log.Println("")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

