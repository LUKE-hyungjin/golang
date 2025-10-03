package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// ============================================================================
// 도메인 모델
// ============================================================================

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Products   []OrderItem `json:"products"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// ============================================================================
// 인터페이스 정의 (Port)
// ============================================================================

// Repository 인터페이스
type UserRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

type ProductRepository interface {
	FindByID(ctx context.Context, id int) (*Product, error)
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	UpdateStock(ctx context.Context, id int, quantity int) error
	List(ctx context.Context, limit, offset int) ([]*Product, error)
}

type OrderRepository interface {
	FindByID(ctx context.Context, id int) (*Order, error)
	FindByUserID(ctx context.Context, userID int) ([]*Order, error)
	Create(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id int, status string) error
}

// Service 인터페이스
type UserService interface {
	GetUser(ctx context.Context, id int) (*User, error)
	CreateUser(ctx context.Context, email, name, role string) (*User, error)
	UpdateUser(ctx context.Context, id int, name string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*User, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, id int) (*Product, error)
	CreateProduct(ctx context.Context, name, description string, price float64, stock int) (*Product, error)
	UpdateStock(ctx context.Context, id int, quantity int) error
	ListProducts(ctx context.Context, page, pageSize int) ([]*Product, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, userID int, items []OrderItem) (*Order, error)
	GetOrder(ctx context.Context, id int) (*Order, error)
	GetUserOrders(ctx context.Context, userID int) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) error
}

// 외부 서비스 인터페이스
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendOrderConfirmation(order *Order, user *User) error
}

type PaymentService interface {
	ProcessPayment(orderID int, amount float64, token string) error
	RefundPayment(orderID int) error
}

type CacheService interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}

type NotificationService interface {
	SendPushNotification(userID int, title, message string) error
	SendSMS(phoneNumber, message string) error
}

// ============================================================================
// Repository 구현체 (Adapter)
// ============================================================================

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	// 실제 구현에서는 SQL 쿼리 수행
	return &User{
		ID:        id,
		Email:     fmt.Sprintf("user%d@example.com", id),
		Name:      fmt.Sprintf("User %d", id),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return &User{
		ID:        1,
		Email:     email,
		Name:      "Test User",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	user.ID = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	users := make([]*User, 0, limit)
	for i := 0; i < limit; i++ {
		users = append(users, &User{
			ID:        offset + i + 1,
			Email:     fmt.Sprintf("user%d@example.com", offset+i+1),
			Name:      fmt.Sprintf("User %d", offset+i+1),
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return users, nil
}

// Mock Repository for testing
type MockUserRepository struct {
	users map[int]*User
}

func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		users: make(map[int]*User),
	}
}

func (r *MockUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) Create(ctx context.Context, user *User) error {
	user.ID = len(r.users) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Update(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Delete(ctx context.Context, id int) error {
	delete(r.users, id)
	return nil
}

func (r *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
	users := make([]*User, 0)
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// ============================================================================
// Service 구현체
// ============================================================================

type UserServiceImpl struct {
	userRepo UserRepository
	cache    CacheService
	email    EmailService
}

func NewUserService(userRepo UserRepository, cache CacheService, email EmailService) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		cache:    cache,
		email:    email,
	}
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id int) (*User, error) {
	// 캐시 확인
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached, err := s.cache.Get(cacheKey); err == nil {
		if user, ok := cached.(*User); ok {
			return user, nil
		}
	}

	// Repository에서 조회
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 캐시 저장
	s.cache.Set(cacheKey, user, 5*time.Minute)

	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, email, name, role string) (*User, error) {
	user := &User{
		Email: email,
		Name:  name,
		Role:  role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 환영 이메일 발송
	s.email.SendEmail(email, "Welcome!", fmt.Sprintf("Welcome %s!", name))

	return user, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id int, name string) (*User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 캐시 무효화
	cacheKey := fmt.Sprintf("user:%d", id)
	s.cache.Delete(cacheKey)

	return user, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 캐시 무효화
	cacheKey := fmt.Sprintf("user:%d", id)
	s.cache.Delete(cacheKey)

	return nil
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, page, pageSize int) ([]*User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(ctx, pageSize, offset)
}

// ============================================================================
// 외부 서비스 구현체
// ============================================================================

type SMTPEmailService struct {
	host     string
	port     int
	username string
	password string
}

func NewSMTPEmailService(host string, port int, username, password string) EmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPEmailService) SendEmail(to, subject, body string) error {
	log.Printf("Sending email to %s: %s", to, subject)
	return nil
}

func (s *SMTPEmailService) SendOrderConfirmation(order *Order, user *User) error {
	subject := fmt.Sprintf("Order #%d Confirmation", order.ID)
	body := fmt.Sprintf("Dear %s, your order has been confirmed.", user.Name)
	return s.SendEmail(user.Email, subject, body)
}

type MockEmailService struct{}

func NewMockEmailService() EmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendEmail(to, subject, body string) error {
	log.Printf("[MOCK] Email sent to %s: %s", to, subject)
	return nil
}

func (s *MockEmailService) SendOrderConfirmation(order *Order, user *User) error {
	log.Printf("[MOCK] Order confirmation sent to %s", user.Email)
	return nil
}

type InMemoryCacheService struct {
	store map[string]interface{}
}

func NewInMemoryCacheService() CacheService {
	return &InMemoryCacheService{
		store: make(map[string]interface{}),
	}
}

func (c *InMemoryCacheService) Get(key string) (interface{}, error) {
	if val, exists := c.store[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("key not found")
}

func (c *InMemoryCacheService) Set(key string, value interface{}, expiration time.Duration) error {
	c.store[key] = value
	// 실제 구현에서는 expiration 처리
	return nil
}

func (c *InMemoryCacheService) Delete(key string) error {
	delete(c.store, key)
	return nil
}

// ============================================================================
// DI Container / Factory
// ============================================================================

type Container struct {
	config            *Config
	db                *sql.DB
	userRepository    UserRepository
	productRepository ProductRepository
	orderRepository   OrderRepository
	userService       UserService
	productService    ProductService
	orderService      OrderService
	emailService      EmailService
	paymentService    PaymentService
	cacheService      CacheService
	notificationService NotificationService
}

type Config struct {
	DatabaseURL string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	Environment string
}

// Factory functions
func NewContainer(config *Config) (*Container, error) {
	c := &Container{config: config}

	// Initialize database
	if config.Environment != "test" {
		db, err := sql.Open("postgres", config.DatabaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		c.db = db
	}

	return c, nil
}

func (c *Container) GetUserRepository() UserRepository {
	if c.userRepository == nil {
		if c.config.Environment == "test" {
			c.userRepository = NewMockUserRepository()
		} else {
			c.userRepository = NewPostgresUserRepository(c.db)
		}
	}
	return c.userRepository
}

func (c *Container) GetEmailService() EmailService {
	if c.emailService == nil {
		if c.config.Environment == "test" {
			c.emailService = NewMockEmailService()
		} else {
			c.emailService = NewSMTPEmailService(
				c.config.SMTPHost,
				c.config.SMTPPort,
				c.config.SMTPUser,
				c.config.SMTPPass,
			)
		}
	}
	return c.emailService
}

func (c *Container) GetCacheService() CacheService {
	if c.cacheService == nil {
		c.cacheService = NewInMemoryCacheService()
	}
	return c.cacheService
}

func (c *Container) GetUserService() UserService {
	if c.userService == nil {
		c.userService = NewUserService(
			c.GetUserRepository(),
			c.GetCacheService(),
			c.GetEmailService(),
		)
	}
	return c.userService
}

// ============================================================================
// HTTP Handlers
// ============================================================================

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var id int
	if err := c.ScanParam("id", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name" binding:"required"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Name, req.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var id int
	if err := c.ScanParam("id", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var id int
	if err := c.ScanParam("id", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	var p, ps int
	fmt.Sscanf(page, "%d", &p)
	fmt.Sscanf(pageSize, "%d", &ps)

	if p < 1 {
		p = 1
	}
	if ps < 1 || ps > 100 {
		ps = 10
	}

	users, err := h.userService.ListUsers(c.Request.Context(), p, ps)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
		"page":  p,
		"page_size": ps,
	})
}

// ============================================================================
// Router Setup with DI
// ============================================================================

func SetupRouter(container *Container) *gin.Engine {
	router := gin.Default()

	// Initialize handlers with injected services
	userHandler := NewUserHandler(container.GetUserService())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"environment": container.config.Environment,
		})
	})

	// DI information endpoint
	router.GET("/di/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"pattern": "Constructor Injection + Factory",
			"benefits": []string{
				"Testability",
				"Loose coupling",
				"Single responsibility",
				"Dependency inversion",
			},
			"services": []string{
				"UserRepository",
				"UserService",
				"EmailService",
				"CacheService",
			},
		})
	})

	// User routes
	users := router.Group("/users")
	{
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.GET("", userHandler.ListUsers)
	}

	// Demonstrate different injection patterns
	patterns := router.Group("/patterns")
	{
		// Constructor injection example
		patterns.GET("/constructor", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"pattern": "Constructor Injection",
				"description": "Dependencies are provided through the constructor",
				"example": "NewUserService(repo, cache, email)",
			})
		})

		// Factory pattern example
		patterns.GET("/factory", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"pattern": "Factory Pattern",
				"description": "Container creates and manages dependencies",
				"example": "container.GetUserService()",
			})
		})

		// Interface segregation example
		patterns.GET("/interface", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"pattern": "Interface Segregation",
				"description": "Clients depend on interfaces, not concrete implementations",
				"example": "UserService interface with multiple implementations",
			})
		})
	}

	return router
}

// ============================================================================
// Main Application
// ============================================================================

func main() {
	// Load configuration
	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost/testdb"),
		SMTPHost:    getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:    587,
		SMTPUser:    getEnv("SMTP_USER", "test@example.com"),
		SMTPPass:    getEnv("SMTP_PASS", "password"),
		Environment: getEnv("APP_ENV", "development"),
	}

	// Create DI container
	container, err := NewContainer(config)
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}

	// Setup router with dependencies
	router := SetupRouter(container)

	// Start server
	log.Printf("🚀 Server starting on :8080 in %s mode", config.Environment)
	log.Println("📦 Dependency Injection Pattern: Constructor Injection + Factory")
	log.Println("🔧 Services initialized with interface-based design")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Package for helper
package main

import "os"