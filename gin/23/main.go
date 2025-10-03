package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== Well-Formatted Code Example ====================

// User represents a user in the system
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

// UserService handles business logic for users
type UserService struct {
	repo UserRepository
	mu   sync.RWMutex
}

// NewUserService creates a new user service instance
func NewUserService(repo UserRepository) *UserService {
	if repo == nil {
		panic("repository cannot be nil")
	}

	return &UserService{
		repo: repo,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uint) (*User, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, username, email string) (*User, error) {
	// Validate input
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrUserExists
	}

	user := &User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// ==================== Custom Errors ====================

var (
	// ErrInvalidID indicates an invalid ID was provided
	ErrInvalidID = errors.New("invalid ID")

	// ErrUserNotFound indicates the user was not found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserExists indicates a user already exists
	ErrUserExists = errors.New("user already exists")

	// ErrInvalidUsername indicates an invalid username
	ErrInvalidUsername = errors.New("invalid username")

	// ErrInvalidEmail indicates an invalid email
	ErrInvalidEmail = errors.New("invalid email")
)

// ==================== Validation Functions ====================

// validateUsername checks if a username is valid
func validateUsername(username string) error {
	username = strings.TrimSpace(username)

	switch {
	case len(username) < 3:
		return fmt.Errorf("%w: too short", ErrInvalidUsername)
	case len(username) > 20:
		return fmt.Errorf("%w: too long", ErrInvalidUsername)
	case !isAlphanumeric(username):
		return fmt.Errorf("%w: must be alphanumeric", ErrInvalidUsername)
	default:
		return nil
	}
}

// validateEmail checks if an email is valid
func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	if !strings.Contains(email, "@") {
		return fmt.Errorf("%w: missing @", ErrInvalidEmail)
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("%w: invalid format", ErrInvalidEmail)
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("%w: empty local or domain part", ErrInvalidEmail)
	}

	return nil
}

// isAlphanumeric checks if a string contains only letters and numbers
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// ==================== HTTP Handlers ====================

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *UserService) *UserHandler {
	if service == nil {
		panic("service cannot be nil")
	}

	return &UserHandler{
		service: service,
	}
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.GetUser(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.CreateUser(ctx, req.Username, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
		case errors.Is(err, ErrInvalidUsername), errors.Is(err, ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ==================== Middleware ====================

// LoggerMiddleware logs incoming requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %v",
			c.Request.Method,
			c.Request.RequestURI,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ==================== Repository Implementation ====================

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[uint]*User
	idSeq uint
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[uint]*User),
		idSeq: 1,
	}
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(_ context.Context, id uint) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Return a copy to avoid mutation
	userCopy := *user
	return &userCopy, nil
}

// FindByEmail finds a user by email
func (r *InMemoryUserRepository) FindByEmail(_ context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, ErrUserNotFound
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(_ context.Context, user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.idSeq
	r.idSeq++
	r.users[user.ID] = user

	return nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(_ context.Context, user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// ==================== Router Setup ====================

// SetupRouter configures and returns the Gin router
func SetupRouter(handler *UserHandler) *gin.Engine {
	router := gin.New()

	// Add middleware
	router.Use(ErrorHandlerMiddleware())
	router.Use(LoggerMiddleware())
	router.Use(gin.Recovery())

	// Define routes
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/:id", handler.GetUser)
			users.POST("", handler.CreateUser)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return router
}

// ==================== Configuration ====================

// Config holds application configuration
type Config struct {
	Port        string
	Environment string
	LogLevel    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	return config
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ==================== Main Function ====================

func main() {
	// Load configuration
	config := LoadConfig()

	// Set Gin mode based on environment
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create dependencies
	repo := NewInMemoryUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	// Setup router
	router := SetupRouter(handler)

	// Seed some test data
	seedTestData(repo)

	// Start server
	addr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server starting on %s in %s mode", addr, config.Environment)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedTestData adds some initial test data
func seedTestData(repo *InMemoryUserRepository) {
	ctx := context.Background()

	users := []*User{
		{
			Username:  "admin",
			Email:     "admin@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			log.Printf("Failed to seed user: %v", err)
		}
	}

	log.Printf("Seeded %d test users", len(users))
}