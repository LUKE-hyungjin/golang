package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project-security/internal/auth"
	"project-security/internal/handlers"
	"project-security/internal/middleware"
	"project-security/internal/repository"
	"project-security/internal/services"
	"project-security/internal/validator"
	"project-security/pkg/config"
	"project-security/pkg/database"
	"project-security/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	log := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtManager)
	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)

	// Initialize validator
	validator.Init()

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(log))
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())

	// CORS configuration
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}
	router.Use(middleware.CORS(corsConfig))

	// Rate limiting
	router.Use(middleware.RateLimit(cfg.RateLimit.RequestsPerMinute))

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// Routes
	setupRoutes(router, healthHandler, authHandler, userHandler, postHandler, jwtManager)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info("Starting server on port " + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}

func setupRoutes(
	router *gin.Engine,
	health *handlers.HealthHandler,
	auth *handlers.AuthHandler,
	user *handlers.UserHandler,
	post *handlers.PostHandler,
	jwtManager *auth.JWTManager,
) {
	// Health check
	router.GET("/health", health.Check)
	router.GET("/ready", health.Ready)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("")
		{
			// Auth endpoints
			public.POST("/auth/register", auth.Register)
			public.POST("/auth/login", auth.Login)
			public.POST("/auth/refresh", auth.RefreshToken)

			// Public post endpoints
			public.GET("/posts", post.List)
			public.GET("/posts/:id", post.Get)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtManager))
		{
			// User endpoints
			protected.GET("/users/profile", user.GetProfile)
			protected.PUT("/users/profile", user.UpdateProfile)
			protected.POST("/users/change-password", user.ChangePassword)

			// Post endpoints (authenticated)
			protected.POST("/posts", post.Create)
			protected.PUT("/posts/:id", middleware.ResourceOwner(), post.Update)
			protected.DELETE("/posts/:id", middleware.ResourceOwner(), post.Delete)

			// Admin only routes
			admin := protected.Group("")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/users", user.List)
				admin.GET("/users/:id", user.Get)
				admin.DELETE("/users/:id", user.Delete)
				admin.GET("/stats", user.GetStats)
			}
		}
	}
}