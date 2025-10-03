package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ========== Application Code ==========

// Models
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username" binding:"required,min=3,max=20"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Repository interface for testing
type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	List(limit, offset int) ([]User, error)
}

// Mock repository for testing
type MockUserRepository struct {
	users map[uint]*User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*User),
	}
}

func (r *MockUserRepository) FindByID(id uint) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *MockUserRepository) FindByEmail(email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) Create(user *User) error {
	if user.ID == 0 {
		user.ID = uint(len(r.users) + 1)
	}
	user.CreatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Update(user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Delete(id uint) error {
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *MockUserRepository) List(limit, offset int) ([]User, error) {
	var users []User
	count := 0
	for _, user := range r.users {
		if count >= offset && len(users) < limit {
			users = append(users, *user)
		}
		count++
	}
	return users, nil
}

// Service layer
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Handlers
type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}{ID: id}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}{ID: id}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.service.repo.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required,min=1"`
	}{ID: id}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	users, err := h.service.repo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"total":  len(users),
	})
}

// Login handler for authentication testing
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate authentication
	if req.Email == "test@example.com" && req.Password == "password123" {
		user := User{
			ID:       1,
			Username: "testuser",
			Email:    req.Email,
		}
		c.JSON(http.StatusOK, LoginResponse{
			Token: "fake-jwt-token",
			User:  user,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// File upload handler for multipart testing
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Check file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files allowed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename,
		"size":     file.Size,
		"message":  "Avatar uploaded successfully",
	})
}

// Middleware for testing
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Simulate token validation
		if token != "Bearer valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", uint(1))
		c.Next()
	}
}

// Setup router for testing
func SetupRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Public routes
	router.POST("/login", handler.Login)
	router.POST("/users", handler.CreateUser)
	router.GET("/users", handler.ListUsers)
	router.GET("/users/:id", handler.GetUser)

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.PUT("/users/:id", handler.UpdateUser)
		protected.DELETE("/users/:id", handler.DeleteUser)
		protected.POST("/users/:id/avatar", handler.UploadAvatar)
	}

	return router
}

// ========== Test Code ==========

// Test helper functions
func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequestWithHeaders(r http.Handler, method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Basic unit tests
func TestGetUser_Success(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare test data
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	repo.Create(testUser)

	// Perform request
	w := performRequest(router, "GET", "/users/1", nil)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, testUser.Username, response.Username)
	assert.Equal(t, testUser.Email, response.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Perform request
	w := performRequest(router, "GET", "/users/999", nil)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

func TestCreateUser_Success(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare request body
	newUser := User{
		Username: "newuser",
		Email:    "new@example.com",
	}
	jsonBody, _ := json.Marshal(newUser)

	// Perform request
	w := performRequest(router, "POST", "/users", bytes.NewBuffer(jsonBody))

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, newUser.Username, response.Username)
	assert.Equal(t, newUser.Email, response.Email)
	assert.NotZero(t, response.ID)
}

func TestCreateUser_ValidationError(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare invalid request body
	invalidUser := User{
		Username: "ab", // Too short
		Email:    "invalid-email",
	}
	jsonBody, _ := json.Marshal(invalidUser)

	// Perform request
	w := performRequest(router, "POST", "/users", bytes.NewBuffer(jsonBody))

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_WithAuth(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare test data
	testUser := &User{
		ID:       1,
		Username: "oldname",
		Email:    "old@example.com",
	}
	repo.Create(testUser)

	// Update data
	updateData := User{
		Username: "newname",
		Email:    "new@example.com",
	}
	jsonBody, _ := json.Marshal(updateData)

	// Perform request with authentication
	headers := map[string]string{
		"Authorization": "Bearer valid-token",
		"Content-Type":  "application/json",
	}
	w := performRequestWithHeaders(router, "PUT", "/users/1", bytes.NewBuffer(jsonBody), headers)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, updateData.Username, response.Username)
	assert.Equal(t, updateData.Email, response.Email)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Update data
	updateData := User{
		Username: "newname",
		Email:    "new@example.com",
	}
	jsonBody, _ := json.Marshal(updateData)

	// Perform request without authentication
	w := performRequest(router, "PUT", "/users/1", bytes.NewBuffer(jsonBody))

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteUser_WithAuth(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare test data
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	repo.Create(testUser)

	// Perform request with authentication
	headers := map[string]string{
		"Authorization": "Bearer valid-token",
	}
	w := performRequestWithHeaders(router, "DELETE", "/users/1", nil, headers)

	// Assertions
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify user is deleted
	_, err := repo.FindByID(1)
	assert.Error(t, err)
}

func TestListUsers_WithPagination(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Create multiple users
	for i := 1; i <= 5; i++ {
		user := &User{
			ID:       uint(i),
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
		}
		repo.Create(user)
	}

	// Test with pagination
	w := performRequest(router, "GET", "/users?limit=2&offset=1", nil)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	users := response["users"].([]interface{})
	assert.Len(t, users, 2)
	assert.Equal(t, float64(2), response["limit"])
	assert.Equal(t, float64(1), response["offset"])
}

func TestLogin_Success(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare login request
	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(loginReq)

	// Perform request
	w := performRequest(router, "POST", "/login", bytes.NewBuffer(jsonBody))

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, loginReq.Email, response.User.Email)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Prepare invalid login request
	loginReq := LoginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginReq)

	// Perform request
	w := performRequest(router, "POST", "/login", bytes.NewBuffer(jsonBody))

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid credentials", response["error"])
}

func TestUploadAvatar_Success(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("avatar", "test.jpg")
	require.NoError(t, err)

	// Write fake image data
	_, err = part.Write([]byte("fake-image-data"))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Perform request with authentication
	req, _ := http.NewRequest("POST", "/users/1/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "test.jpg", response["filename"])
	assert.Equal(t, "Avatar uploaded successfully", response["message"])
}

// Test Suite using testify/suite
type UserHandlerTestSuite struct {
	suite.Suite
	router  *gin.Engine
	repo    *MockUserRepository
	handler *UserHandler
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.repo = NewMockUserRepository()
	service := NewUserService(suite.repo)
	suite.handler = NewUserHandler(service)
	suite.router = SetupRouter(suite.handler)
}

func (suite *UserHandlerTestSuite) TestUserCRUDFlow() {
	// Create user
	newUser := User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	jsonBody, _ := json.Marshal(newUser)

	w := performRequest(suite.router, "POST", "/users", bytes.NewBuffer(jsonBody))
	suite.Equal(http.StatusCreated, w.Code)

	var createdUser User
	json.Unmarshal(w.Body.Bytes(), &createdUser)
	userID := createdUser.ID

	// Get user
	w = performRequest(suite.router, "GET", fmt.Sprintf("/users/%d", userID), nil)
	suite.Equal(http.StatusOK, w.Code)

	// Update user (with auth)
	updateData := User{
		Username: "updateduser",
		Email:    "updated@example.com",
	}
	jsonBody, _ = json.Marshal(updateData)

	headers := map[string]string{
		"Authorization": "Bearer valid-token",
		"Content-Type":  "application/json",
	}
	w = performRequestWithHeaders(suite.router, "PUT", fmt.Sprintf("/users/%d", userID), bytes.NewBuffer(jsonBody), headers)
	suite.Equal(http.StatusOK, w.Code)

	// Delete user (with auth)
	w = performRequestWithHeaders(suite.router, "DELETE", fmt.Sprintf("/users/%d", userID), nil, headers)
	suite.Equal(http.StatusNoContent, w.Code)

	// Verify deletion
	w = performRequest(suite.router, "GET", fmt.Sprintf("/users/%d", userID), nil)
	suite.Equal(http.StatusNotFound, w.Code)
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

// Table-driven tests
func TestUserValidation_TableDriven(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	tests := []struct {
		name         string
		user         User
		expectedCode int
		expectedErr  string
	}{
		{
			name: "Valid user",
			user: User{
				Username: "validuser",
				Email:    "valid@example.com",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Username too short",
			user: User{
				Username: "ab",
				Email:    "valid@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Username too long",
			user: User{
				Username: strings.Repeat("a", 21),
				Email:    "valid@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			user: User{
				Username: "validuser",
				Email:    "invalid-email",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty username",
			user: User{
				Username: "",
				Email:    "valid@example.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty email",
			user: User{
				Username: "validuser",
				Email:    "",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.user)
			w := performRequest(router, "POST", "/users", bytes.NewBuffer(jsonBody))
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

// Benchmark tests
func BenchmarkGetUser(b *testing.B) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Create test user
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	repo.Create(testUser)

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest(router, "GET", "/users/1", nil)
		if w.Code != http.StatusOK {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkCreateUser(b *testing.B) {
	// Setup
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	user := User{
		Username: "benchuser",
		Email:    "bench@example.com",
	}
	jsonBody, _ := json.Marshal(user)

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := performRequest(router, "POST", "/users", bytes.NewBuffer(jsonBody))
		if w.Code != http.StatusCreated {
			b.Errorf("Expected status 201, got %d", w.Code)
		}
	}
}

// Custom test reporter
type CustomTestReporter struct {
	passed int
	failed int
	tests  []string
}

func (r *CustomTestReporter) Record(name string, passed bool) {
	r.tests = append(r.tests, name)
	if passed {
		r.passed++
	} else {
		r.failed++
	}
}

func (r *CustomTestReporter) Summary() {
	fmt.Printf("\n=== Test Summary ===\n")
	fmt.Printf("Total: %d | Passed: %d | Failed: %d\n", r.passed+r.failed, r.passed, r.failed)
	fmt.Printf("Tests run: %v\n", r.tests)
}

// Example of custom test runner
func runCustomTests() {
	reporter := &CustomTestReporter{}

	// Run individual tests
	tests := []struct {
		name string
		fn   func() bool
	}{
		{"TestGetUser", func() bool {
			// Test implementation
			return true // simplified
		}},
		{"TestCreateUser", func() bool {
			// Test implementation
			return true // simplified
		}},
		{"TestUpdateUser", func() bool {
			// Test implementation
			return true // simplified
		}},
	}

	for _, test := range tests {
		passed := test.fn()
		reporter.Record(test.name, passed)
	}

	reporter.Summary()
}

// Main function to demonstrate testing
func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("Running tests...")

		// Run custom tests
		runCustomTests()

		// You would normally use `go test` command instead
		fmt.Println("\nTo run actual tests, use: go test -v")
		fmt.Println("To run with coverage: go test -cover")
		fmt.Println("To run benchmarks: go test -bench=.")
		fmt.Println("To run specific test: go test -run TestGetUser")
		fmt.Println("To run test suite: go test -run TestUserHandlerSuite")
		return
	}

	// Normal application startup
	gin.SetMode(gin.ReleaseMode)

	repo := NewMockUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupRouter(handler)

	// Seed some data
	repo.Create(&User{
		Username: "admin",
		Email:    "admin@example.com",
	})

	fmt.Println("Server starting on :8080...")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /login")
	fmt.Println("  GET    /users")
	fmt.Println("  GET    /users/:id")
	fmt.Println("  POST   /users")
	fmt.Println("  PUT    /users/:id (requires auth)")
	fmt.Println("  DELETE /users/:id (requires auth)")
	fmt.Println("  POST   /users/:id/avatar (requires auth)")
	fmt.Println("\nRun with 'test' argument to see test examples")

	log.Fatal(router.Run(":8080"))
}