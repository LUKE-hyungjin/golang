package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ========== Models ==========

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Posts     []Post    `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comments  []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Tags      []Tag     `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"not null"`
	PostID    uint      `json:"post_id"`
	Post      *Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"unique;not null"`
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;"`
}

// ========== Database ==========

type Database struct {
	*gorm.DB
}

func NewDatabase(dsn string, config *gorm.Config) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) Migrate() error {
	return db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Tag{})
}

// Test database with transaction support
type TestDatabase struct {
	*Database
	tx *gorm.DB
}

func NewTestDatabase() (*TestDatabase, error) {
	// Use in-memory database for tests
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := NewDatabase(":memory:", config)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		return nil, err
	}

	return &TestDatabase{Database: db}, nil
}

func (tdb *TestDatabase) Begin() {
	tdb.tx = tdb.DB.Begin()
}

func (tdb *TestDatabase) Rollback() {
	if tdb.tx != nil {
		tdb.tx.Rollback()
	}
}

func (tdb *TestDatabase) Commit() {
	if tdb.tx != nil {
		tdb.tx.Commit()
	}
}

func (tdb *TestDatabase) GetDB() *gorm.DB {
	if tdb.tx != nil {
		return tdb.tx
	}
	return tdb.DB
}

// ========== Repositories ==========

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Posts").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindByID(id uint) (*Post, error) {
	var post Post
	err := r.db.Preload("User").Preload("Comments.User").Preload("Tags").
		First(&post, id).Error
	return &post, err
}

func (r *PostRepository) List(limit, offset int) ([]Post, error) {
	var posts []Post
	err := r.db.Preload("User").Preload("Tags").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Update(post *Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&Post{}, id).Error
}

// ========== Services ==========

type BlogService struct {
	userRepo *UserRepository
	postRepo *PostRepository
	db       *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{
		userRepo: NewUserRepository(db),
		postRepo: NewPostRepository(db),
		db:       db,
	}
}

func (s *BlogService) CreateUserWithPost(username, email, password, title, content string) (*User, error) {
	// Use transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	post := &Post{
		Title:   title,
		Content: content,
		UserID:  user.ID,
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with associations
	s.userRepo.FindByID(user.ID)
	return user, nil
}

// ========== Handlers ==========

type BlogHandler struct {
	service *BlogService
}

func NewBlogHandler(service *BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (h *BlogHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // Should be hashed in production
	}

	if err := h.service.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *BlogHandler) GetUser(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required"`
	}{ID: id}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.userRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *BlogHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title   string   `json:"title" binding:"required"`
		Content string   `json:"content"`
		UserID  uint     `json:"user_id" binding:"required"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	// Handle tags
	if len(req.Tags) > 0 {
		var tags []Tag
		for _, tagName := range req.Tags {
			var tag Tag
			h.service.db.FirstOrCreate(&tag, Tag{Name: tagName})
			tags = append(tags, tag)
		}
		post.Tags = tags
	}

	if err := h.service.postRepo.Create(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *BlogHandler) GetPost(c *gin.Context) {
	var id uint
	if err := c.ShouldBindUri(&struct {
		ID uint `uri:"id" binding:"required"`
	}{ID: id}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.service.postRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *BlogHandler) ListPosts(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	posts, err := h.service.postRepo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *BlogHandler) CreateComment(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		PostID  uint   `json:"post_id" binding:"required"`
		UserID  uint   `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  req.UserID,
	}

	if err := h.service.db.Create(comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// Health check endpoint
func (h *BlogHandler) HealthCheck(c *gin.Context) {
	sqlDB, err := h.service.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "blog-api",
	})
}

// ========== Router Setup ==========

func SetupRouter(handler *BlogHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Users
		v1.POST("/users", handler.CreateUser)
		v1.GET("/users/:id", handler.GetUser)

		// Posts
		v1.POST("/posts", handler.CreatePost)
		v1.GET("/posts", handler.ListPosts)
		v1.GET("/posts/:id", handler.GetPost)

		// Comments
		v1.POST("/comments", handler.CreateComment)
	}

	return router
}

// ========== Test Helpers ==========

type TestServer struct {
	Router  *gin.Engine
	DB      *TestDatabase
	Service *BlogService
	Handler *BlogHandler
}

func NewTestServer() (*TestServer, error) {
	gin.SetMode(gin.TestMode)

	db, err := NewTestDatabase()
	if err != nil {
		return nil, err
	}

	service := NewBlogService(db.GetDB())
	handler := NewBlogHandler(service)
	router := SetupRouter(handler)

	return &TestServer{
		Router:  router,
		DB:      db,
		Service: service,
		Handler: handler,
	}, nil
}

func (ts *TestServer) Cleanup() {
	if ts.DB != nil {
		sqlDB, _ := ts.DB.DB.DB()
		sqlDB.Close()
	}
}

func (ts *TestServer) SeedTestData() {
	// Create test users
	users := []User{
		{Username: "alice", Email: "alice@example.com", Password: "password123"},
		{Username: "bob", Email: "bob@example.com", Password: "password123"},
		{Username: "charlie", Email: "charlie@example.com", Password: "password123"},
	}

	for _, user := range users {
		ts.DB.Create(&user)
	}

	// Create test posts
	posts := []Post{
		{Title: "First Post", Content: "Hello World", UserID: 1},
		{Title: "Second Post", Content: "Testing Integration", UserID: 1},
		{Title: "Bob's Post", Content: "Bob's content", UserID: 2},
	}

	for _, post := range posts {
		ts.DB.Create(&post)
	}

	// Create test tags
	tags := []Tag{
		{Name: "golang"},
		{Name: "testing"},
		{Name: "gin"},
	}

	for _, tag := range tags {
		ts.DB.Create(&tag)
	}
}

// ========== Integration Tests ==========

func TestHealthCheck_Integration(t *testing.T) {
	server, err := NewTestServer()
	require.NoError(t, err)
	defer server.Cleanup()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestCreateUser_Integration(t *testing.T) {
	server, err := NewTestServer()
	require.NoError(t, err)
	defer server.Cleanup()

	user := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
	assert.NotZero(t, response.ID)

	// Verify in database
	var dbUser User
	err = server.DB.First(&dbUser, response.ID).Error
	require.NoError(t, err)
	assert.Equal(t, response.Username, dbUser.Username)
}

func TestCreatePost_WithTags_Integration(t *testing.T) {
	server, err := NewTestServer()
	require.NoError(t, err)
	defer server.Cleanup()

	// Create a user first
	user := &User{
		Username: "poster",
		Email:    "poster@example.com",
		Password: "password123",
	}
	server.DB.Create(user)

	// Create post with tags
	post := map[string]interface{}{
		"title":   "Test Post",
		"content": "Test Content",
		"user_id": user.ID,
		"tags":    []string{"test", "integration", "golang"},
	}
	jsonBody, _ := json.Marshal(post)

	req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response Post
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Test Post", response.Title)
	assert.Len(t, response.Tags, 3)

	// Verify tags in database
	var dbPost Post
	err = server.DB.Preload("Tags").First(&dbPost, response.ID).Error
	require.NoError(t, err)
	assert.Len(t, dbPost.Tags, 3)
}

func TestUserPostFlow_Integration(t *testing.T) {
	server, err := NewTestServer()
	require.NoError(t, err)
	defer server.Cleanup()

	// Step 1: Create user
	user := map[string]string{
		"username": "flowuser",
		"email":    "flow@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser User
	json.Unmarshal(w.Body.Bytes(), &createdUser)

	// Step 2: Create post for user
	post := map[string]interface{}{
		"title":   "Flow Test Post",
		"content": "Flow test content",
		"user_id": createdUser.ID,
	}
	jsonBody, _ = json.Marshal(post)

	req, _ = http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdPost Post
	json.Unmarshal(w.Body.Bytes(), &createdPost)

	// Step 3: Add comment to post
	comment := map[string]interface{}{
		"content": "Great post!",
		"post_id": createdPost.ID,
		"user_id": createdUser.ID,
	}
	jsonBody, _ = json.Marshal(comment)

	req, _ = http.NewRequest("POST", "/api/v1/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 4: Get post with all associations
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%d", createdPost.ID), nil)
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fullPost Post
	json.Unmarshal(w.Body.Bytes(), &fullPost)
	assert.Equal(t, "Flow Test Post", fullPost.Title)
	assert.NotNil(t, fullPost.User)
	assert.Len(t, fullPost.Comments, 1)
}

// ========== Test Suite ==========

type BlogIntegrationSuite struct {
	suite.Suite
	server *TestServer
}

func (suite *BlogIntegrationSuite) SetupSuite() {
	server, err := NewTestServer()
	suite.Require().NoError(err)
	suite.server = server
}

func (suite *BlogIntegrationSuite) TearDownSuite() {
	suite.server.Cleanup()
}

func (suite *BlogIntegrationSuite) SetupTest() {
	// Start transaction for each test
	suite.server.DB.Begin()
}

func (suite *BlogIntegrationSuite) TearDownTest() {
	// Rollback transaction after each test
	suite.server.DB.Rollback()
}

func (suite *BlogIntegrationSuite) TestCompleteScenario() {
	// Seed initial data
	suite.server.SeedTestData()

	// Test listing posts
	req, _ := http.NewRequest("GET", "/api/v1/posts?limit=2", nil)
	w := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	posts := response["posts"].([]interface{})
	suite.Len(posts, 2)

	// Test getting specific user with posts
	req, _ = http.NewRequest("GET", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	suite.server.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var user User
	json.Unmarshal(w.Body.Bytes(), &user)
	suite.Equal("alice", user.Username)
	suite.NotEmpty(user.Posts)
}

func (suite *BlogIntegrationSuite) TestConcurrentRequests() {
	done := make(chan bool, 10)

	// Create 10 concurrent requests
	for i := 0; i < 10; i++ {
		go func(index int) {
			user := map[string]string{
				"username": fmt.Sprintf("concurrent%d", index),
				"email":    fmt.Sprintf("concurrent%d@example.com", index),
				"password": "password123",
			}
			jsonBody, _ := json.Marshal(user)

			req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(w, req)

			suite.Equal(http.StatusCreated, w.Code)
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users were created
	var count int64
	suite.server.DB.Model(&User{}).Where("username LIKE ?", "concurrent%").Count(&count)
	suite.Equal(int64(10), count)
}

func TestBlogIntegrationSuite(t *testing.T) {
	suite.Run(t, new(BlogIntegrationSuite))
}

// ========== Benchmark Tests ==========

func BenchmarkCreateUser_Integration(b *testing.B) {
	server, _ := NewTestServer()
	defer server.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := map[string]string{
			"username": fmt.Sprintf("bench%d", i),
			"email":    fmt.Sprintf("bench%d@example.com", i),
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(user)

		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Errorf("Expected status 201, got %d", w.Code)
		}
	}
}

func BenchmarkListPosts_Integration(b *testing.B) {
	server, _ := NewTestServer()
	defer server.Cleanup()
	server.SeedTestData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}

// ========== Test Fixtures ==========

type TestFixtures struct {
	Users    []User
	Posts    []Post
	Comments []Comment
	Tags     []Tag
}

func LoadTestFixtures(db *gorm.DB) *TestFixtures {
	fixtures := &TestFixtures{
		Users: []User{
			{Username: "admin", Email: "admin@example.com", Password: "admin123"},
			{Username: "editor", Email: "editor@example.com", Password: "editor123"},
			{Username: "viewer", Email: "viewer@example.com", Password: "viewer123"},
		},
		Posts: []Post{
			{Title: "Welcome Post", Content: "Welcome to our blog", UserID: 1},
			{Title: "Tutorial", Content: "How to use our platform", UserID: 1},
			{Title: "News", Content: "Latest updates", UserID: 2},
		},
		Comments: []Comment{
			{Content: "Great post!", PostID: 1, UserID: 2},
			{Content: "Thanks for sharing", PostID: 1, UserID: 3},
		},
		Tags: []Tag{
			{Name: "welcome"},
			{Name: "tutorial"},
			{Name: "news"},
		},
	}

	// Load fixtures into database
	for _, user := range fixtures.Users {
		db.Create(&user)
	}
	for _, post := range fixtures.Posts {
		db.Create(&post)
	}
	for _, comment := range fixtures.Comments {
		db.Create(&comment)
	}
	for _, tag := range fixtures.Tags {
		db.Create(&tag)
	}

	return fixtures
}

// ========== Test with Context ==========

func TestWithTimeout_Integration(t *testing.T) {
	server, err := NewTestServer()
	require.NoError(t, err)
	defer server.Cleanup()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Simulate slow operation
	go func() {
		time.Sleep(200 * time.Millisecond)
		user := map[string]string{
			"username": "slowuser",
			"email":    "slow@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(user)

		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
	}()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	case <-time.After(300 * time.Millisecond):
		t.Error("Context should have timed out")
	}
}

// ========== Main Function ==========

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		fmt.Println("Running integration tests...")
		fmt.Println("\nTo run tests, use:")
		fmt.Println("  go test -v                    # Run all tests")
		fmt.Println("  go test -run Integration      # Run integration tests only")
		fmt.Println("  go test -run Suite           # Run test suite")
		fmt.Println("  go test -bench=.             # Run benchmarks")
		fmt.Println("  go test -cover               # Check coverage")
		return
	}

	// Production setup
	gin.SetMode(gin.ReleaseMode)

	dbFile := "blog.db"
	if os.Getenv("DB_FILE") != "" {
		dbFile = os.Getenv("DB_FILE")
	}

	db, err := NewDatabase(dbFile, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	service := NewBlogService(db.DB)
	handler := NewBlogHandler(service)
	router := SetupRouter(handler)

	// Seed some initial data if database is empty
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding initial data...")
		LoadTestFixtures(db.DB)
	}

	fmt.Println("Server starting on :8080...")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health")
	fmt.Println("  POST   /api/v1/users")
	fmt.Println("  GET    /api/v1/users/:id")
	fmt.Println("  POST   /api/v1/posts")
	fmt.Println("  GET    /api/v1/posts")
	fmt.Println("  GET    /api/v1/posts/:id")
	fmt.Println("  POST   /api/v1/comments")
	fmt.Println("\nRun with 'test' argument to see test instructions")

	log.Fatal(router.Run(":8080"))
}