package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/banking-system/internal/config"
	"github.com/example/banking-system/internal/handlers"
	"github.com/example/banking-system/internal/middleware"
	"github.com/example/banking-system/internal/models"
	"github.com/example/banking-system/internal/services"
	"github.com/example/banking-system/pkg/database"
	"github.com/example/banking-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// 2. Initialize logger
	log := logger.NewLogger(cfg.Logging)
	log.Info("🚀 Starting Banking System")
	log.Info(fmt.Sprintf("Environment: %s", cfg.App.Environment))
	log.Info(fmt.Sprintf("Mode: %s", cfg.App.Mode))

	// 3. Set Gin mode based on configuration
	switch cfg.App.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Initialize database with migration
	db, err := database.InitDB(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{},
		&models.AuditLog{},
		&models.Migration{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 5. Initialize services with dependency injection
	container := &services.Container{
		DB:     db,
		Config: cfg,
		Logger: log,
	}

	userService := services.NewUserService(container)
	accountService := services.NewAccountService(container)
	transactionService := services.NewTransactionService(container)
	auditService := services.NewAuditService(container)
	migrationService := services.NewMigrationService(container)

	// 6. Run migrations
	if err := migrationService.RunMigrations(); err != nil {
		log.Error("Failed to run migrations:", err)
	}

	// 7. Seed data in development mode
	if cfg.App.Environment == "development" {
		seeder := services.NewSeederService(container)
		if count, _ := seeder.GetUserCount(); count == 0 {
			log.Info("Seeding initial data...")
			if err := seeder.SeedAll(); err != nil {
				log.Error("Failed to seed data:", err)
			}
		}
	}

	// 8. Initialize router with middlewares
	router := gin.New()

	// Global middlewares (Lesson 10, 11)
	router.Use(middleware.RequestID())
	router.Use(middleware.LoggingMiddleware(log))
	router.Use(middleware.RecoveryMiddleware(log))
	router.Use(middleware.ErrorHandler())

	// Mode-specific middlewares (Lesson 14)
	if cfg.App.Mode == "debug" {
		router.Use(middleware.DebugMiddleware())
	} else if cfg.App.Mode == "release" {
		router.Use(middleware.SecurityHeaders())
		router.Use(middleware.RateLimiter(cfg.Security.RateLimit))
	}

	// 9. Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	userHandler := handlers.NewUserHandler(userService, auditService)
	accountHandler := handlers.NewAccountHandler(accountService, auditService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, auditService)
	migrationHandler := handlers.NewMigrationHandler(migrationService)
	adminHandler := handlers.NewAdminHandler(container)

	// 10. Setup routes
	setupRoutes(router, healthHandler, userHandler, accountHandler, transactionHandler, migrationHandler, adminHandler, cfg)

	// 11. Setup graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info(fmt.Sprintf("Server listening on port %d", cfg.Server.Port))
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exited")
}

func setupRoutes(
	router *gin.Engine,
	health *handlers.HealthHandler,
	user *handlers.UserHandler,
	account *handlers.AccountHandler,
	transaction *handlers.TransactionHandler,
	migration *handlers.MigrationHandler,
	admin *handlers.AdminHandler,
	cfg *config.Config,
) {
	// Health check endpoints
	router.GET("/health", health.Check)
	router.GET("/health/db", health.DBCheck)
	router.GET("/ready", health.Ready)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public endpoints
		public := v1.Group("")
		{
			public.POST("/register", user.Register)
			public.POST("/login", user.Login)
		}

		// Protected endpoints (would add JWT middleware in production)
		protected := v1.Group("")
		// protected.Use(middleware.AuthMiddleware())
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("", user.List)
				users.GET("/:id", user.Get)
				users.PUT("/:id", user.Update)
				users.DELETE("/:id", user.Delete)
				users.GET("/:id/accounts", user.GetAccounts)
			}

			// Account management
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", account.Create)
				accounts.GET("", account.List)
				accounts.GET("/:id", account.Get)
				accounts.PUT("/:id", account.Update)
				accounts.DELETE("/:id", account.Delete)
				accounts.POST("/:id/deposit", account.Deposit)
				accounts.POST("/:id/withdraw", account.Withdraw)
				accounts.GET("/:id/transactions", account.GetTransactions)
				accounts.GET("/:id/balance", account.GetBalance)
			}

			// Transaction management
			transactions := protected.Group("/transactions")
			{
				transactions.POST("/transfer", transaction.Transfer)
				transactions.GET("", transaction.List)
				transactions.GET("/:id", transaction.Get)
				transactions.GET("/report", transaction.Report)
			}
		}
	}

	// Admin endpoints
	if cfg.App.Environment != "production" {
		adminGroup := router.Group("/admin")
		{
			// Migration management (Lesson 16)
			adminGroup.GET("/migrations", migration.Status)
			adminGroup.POST("/migrations/run", migration.Run)
			adminGroup.POST("/migrations/rollback", migration.Rollback)
			adminGroup.POST("/migrations/reset", migration.Reset)

			// Seed data management
			adminGroup.POST("/seed", admin.Seed)
			adminGroup.POST("/seed/clean", admin.Clean)
			adminGroup.POST("/seed/export", admin.Export)
			adminGroup.POST("/seed/import", admin.Import)

			// System info
			adminGroup.GET("/info", admin.SystemInfo)
			adminGroup.GET("/config", admin.GetConfig)
			adminGroup.GET("/metrics", admin.Metrics)

			// Testing endpoints
			adminGroup.POST("/test/concurrent", admin.TestConcurrency)
			adminGroup.POST("/test/deadlock", admin.TestDeadlock)
			adminGroup.POST("/test/timeout", admin.TestTimeout)
		}
	}

	// Debug endpoints (Lesson 14)
	if cfg.App.Mode == "debug" {
		debug := router.Group("/debug")
		{
			debug.GET("/vars", admin.DebugVars)
			debug.GET("/pprof/*action", admin.PProf)
			debug.GET("/routes", admin.Routes)
		}
	}

	// Static files and templates (if needed)
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/**/*")
}