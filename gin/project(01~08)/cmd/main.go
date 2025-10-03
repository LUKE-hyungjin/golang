package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gin-project/internal/handlers"
	"gin-project/internal/middleware"
)

// Template functions
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatCurrency(price float64) string {
	return fmt.Sprintf("₩%.0f", price)
}

func timeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	return formatDate(t)
}

func main() {
	// Create Gin engine without default middleware
	r := gin.New()

	// Register custom template functions
	r.SetFuncMap(template.FuncMap{
		"formatDate":     formatDate,
		"formatDateTime": formatDateTime,
		"formatCurrency": formatCurrency,
		"timeAgo":        timeAgo,
	})

	// Load templates
	r.LoadHTMLGlob("web/templates/*")

	// ========================================
	// Global Middleware (Lesson 05)
	// ========================================
	r.Use(gin.Recovery())                   // Panic recovery
	r.Use(middleware.LoggerMiddleware())    // Custom logger
	r.Use(middleware.CORSMiddleware())      // CORS
	r.Use(middleware.RequestIDMiddleware()) // Request ID

	// ========================================
	// Static Files (Lesson 07)
	// ========================================
	r.Static("/static", "./web/static")
	r.Static("/uploads", "./uploads")
	r.StaticFile("/favicon.ico", "./web/static/favicon.ico")

	// ========================================
	// Public Routes (Lesson 08 - Template Rendering)
	// ========================================
	r.GET("/", handlers.RenderHome)
	r.GET("/posts", handlers.RenderPosts)
	r.GET("/posts/:id", handlers.RenderPost)
	r.GET("/login", handlers.RenderLogin)
	r.GET("/register", handlers.RenderRegister)

	// ========================================
	// API Version 1 (Lessons 01-06)
	// ========================================
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RateLimitMiddleware(100)) // Rate limiting
	{
		// Health check (Lesson 01)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":     "healthy",
				"version":    "1.0",
				"time":       time.Now(),
				"request_id": c.GetString("RequestID"),
			})
		})

		// Auth routes (Lesson 02, 03)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Posts routes (Lesson 02, 03, 04)
		posts := v1.Group("/posts")
		{
			posts.GET("", handlers.GetPosts)           // List with pagination
			posts.GET("/:id", handlers.GetPost)        // Get single post

			// Protected routes - require authentication
			protected := posts.Group("")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("", handlers.CreatePost)
				protected.PUT("/:id", handlers.UpdatePost)
				protected.DELETE("/:id", handlers.DeletePost)

				// Comments
				protected.GET("/:id/comments", handlers.GetComments)
				protected.POST("/comments", handlers.CreateComment)
				protected.DELETE("/comments/:id", handlers.DeleteComment)
			}
		}

		// Upload routes (Lesson 03, 07)
		upload := v1.Group("/upload")
		upload.Use(middleware.AuthMiddleware())
		{
			upload.POST("/image", handlers.UploadImage)
		}

		// Protected user routes
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", handlers.GetProfile)
		}
	}

	// ========================================
	// API Version 2 (Lesson 06 - API Versioning)
	// ========================================
	v2 := r.Group("/api/v2")
	v2.Use(middleware.RateLimitMiddleware(150)) // Higher rate limit for v2
	v2.Use(func(c *gin.Context) {
		c.Header("X-API-Version", "2.0")
		c.Next()
	})
	{
		v2.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"version": "2.0",
				"features": []string{
					"enhanced-pagination",
					"advanced-filtering",
					"real-time-updates",
				},
				"time": time.Now(),
			})
		})

		// Enhanced posts API with better response structure
		v2Posts := v2.Group("/posts")
		{
			v2Posts.GET("", func(c *gin.Context) {
				handlers.GetPosts(c)
				// V2 wraps response in data field
			})
		}
	}

	// ========================================
	// Admin Routes (Lesson 05, 06)
	// ========================================
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", handlers.RenderAdminDashboard)

		admin.GET("/api/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"total_users":    2,
				"total_posts":    2,
				"total_comments": 2,
				"uptime":         "100%",
			})
		})
	}

	// ========================================
	// 404 Handler (Lesson 08)
	// ========================================
	r.NoRoute(func(c *gin.Context) {
		// Check if it's an API request
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:5] == "/api/" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "API endpoint not found",
				"path":  c.Request.URL.Path,
			})
			return
		}

		// Render 404 page for web requests
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Page Not Found",
			"path":  c.Request.URL.Path,
		})
	})

	// ========================================
	// Server Startup
	// ========================================
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║   Blog Community Platform - Gin Framework         ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Server Configuration:")
	fmt.Println("  → Port: 8080")
	fmt.Println("  → Environment: Development")
	fmt.Println()
	fmt.Println("Available Endpoints:")
	fmt.Println("  🌐 Web Interface:")
	fmt.Println("     → Home:         http://localhost:8080/")
	fmt.Println("     → Posts:        http://localhost:8080/posts")
	fmt.Println("     → Login:        http://localhost:8080/login")
	fmt.Println("     → Register:     http://localhost:8080/register")
	fmt.Println()
	fmt.Println("  📡 API v1:         http://localhost:8080/api/v1")
	fmt.Println("     → Health:       GET  /api/v1/health")
	fmt.Println("     → Register:     POST /api/v1/auth/register")
	fmt.Println("     → Login:        POST /api/v1/auth/login")
	fmt.Println("     → Posts:        GET  /api/v1/posts")
	fmt.Println("     → Create Post:  POST /api/v1/posts (auth required)")
	fmt.Println()
	fmt.Println("  📡 API v2:         http://localhost:8080/api/v2")
	fmt.Println("     → Health:       GET  /api/v2/health")
	fmt.Println()
	fmt.Println("  🔧 Admin Panel:    http://localhost:8080/admin/dashboard")
	fmt.Println()
	fmt.Println("Test Credentials:")
	fmt.Println("  → Admin:  admin / password123")
	fmt.Println("  → User:   user / password123")
	fmt.Println()
	fmt.Println("API Authentication:")
	fmt.Println("  → Admin Token:  admin-token-456")
	fmt.Println("  → User Token:   valid-token-123")
	fmt.Println()
	fmt.Println("Server starting...")
	fmt.Println("════════════════════════════════════════════════════")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
