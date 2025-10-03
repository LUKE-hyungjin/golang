package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================================
// 모델 정의
// ============================================================================

// Base 모델 - 모든 모델이 상속
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 모델
type User struct {
	Base
	Email    string    `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Username string    `gorm:"uniqueIndex;not null;size:50" json:"username" binding:"required,min=3,max=50"`
	Name     string    `gorm:"size:100" json:"name" binding:"required"`
	Age      int       `json:"age" binding:"min=0,max=150"`
	Bio      string    `gorm:"type:text" json:"bio"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// Post 모델
type Post struct {
	Base
	Title      string    `gorm:"not null;size:200" json:"title" binding:"required"`
	Content    string    `gorm:"type:text" json:"content" binding:"required"`
	Slug       string    `gorm:"uniqueIndex;not null" json:"slug"`
	Published  bool      `gorm:"default:false;index" json:"published"`
	ViewCount  int       `gorm:"default:0" json:"view_count"`
	UserID     uint      `json:"user_id" binding:"required"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tags       []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	Comments   []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CategoryID *uint     `json:"category_id"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// Category 모델
type Category struct {
	Base
	Name        string `gorm:"uniqueIndex;not null;size:50" json:"name" binding:"required"`
	Description string `gorm:"type:text" json:"description"`
	Posts       []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// Tag 모델
type Tag struct {
	Base
	Name  string `gorm:"uniqueIndex;not null;size:30" json:"name" binding:"required"`
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// Comment 모델
type Comment struct {
	Base
	Content string `gorm:"type:text;not null" json:"content" binding:"required"`
	UserID  uint   `json:"user_id" binding:"required"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PostID  uint   `json:"post_id" binding:"required"`
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// ============================================================================
// 데이터베이스 연결 및 초기화
// ============================================================================

type Database struct {
	*gorm.DB
}

func NewDatabase(debug bool) (*Database, error) {
	// SQLite 연결
	logLevel := logger.Error
	if debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 마이그레이션
	if err := db.AutoMigrate(&User{}, &Post{}, &Category{}, &Tag{}, &Comment{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return &Database{db}, nil
}

// ============================================================================
// Repository 패턴
// ============================================================================

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// Create - 사용자 생성
func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

// FindByID - ID로 사용자 조회
func (r *UserRepository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Posts").Preload("Comments").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail - 이메일로 사용자 조회
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll - 모든 사용자 조회 (페이지네이션)
func (r *UserRepository) FindAll(offset, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	// 전체 개수
	r.db.Model(&User{}).Count(&total)

	// 페이지네이션 적용
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// Update - 사용자 업데이트
func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

// UpdateFields - 특정 필드만 업데이트
func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete - 사용자 삭제 (소프트 삭제)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

// HardDelete - 사용자 완전 삭제
func (r *UserRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&User{}, id).Error
}

// PostRepository
type PostRepository struct {
	db *Database
}

func NewPostRepository(db *Database) *PostRepository {
	return &PostRepository{db: db}
}

// Create - 포스트 생성
func (r *PostRepository) Create(post *Post) error {
	// Slug 자동 생성
	if post.Slug == "" {
		post.Slug = fmt.Sprintf("%s-%d", slugify(post.Title), time.Now().Unix())
	}
	return r.db.Create(post).Error
}

// FindByID - ID로 포스트 조회
func (r *PostRepository) FindByID(id uint) (*Post, error) {
	var post Post
	err := r.db.Preload("User").
		Preload("Tags").
		Preload("Category").
		Preload("Comments.User").
		First(&post, id).Error
	if err != nil {
		return nil, err
	}

	// 조회수 증가
	r.db.Model(&post).Update("view_count", post.ViewCount+1)

	return &post, nil
}

// FindBySlug - Slug로 포스트 조회
func (r *PostRepository) FindBySlug(slug string) (*Post, error) {
	var post Post
	err := r.db.Where("slug = ?", slug).
		Preload("User").
		Preload("Tags").
		Preload("Category").
		First(&post).Error
	return &post, err
}

// FindAll - 모든 포스트 조회 (필터링 + 페이지네이션)
func (r *PostRepository) FindAll(filters map[string]interface{}, offset, limit int) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := r.db.Model(&Post{})

	// 필터링
	if published, ok := filters["published"].(bool); ok {
		query = query.Where("published = ?", published)
	}
	if userID, ok := filters["user_id"].(uint); ok {
		query = query.Where("user_id = ?", userID)
	}
	if categoryID, ok := filters["category_id"].(uint); ok {
		query = query.Where("category_id = ?", categoryID)
	}

	// 전체 개수
	query.Count(&total)

	// 페이지네이션 및 조회
	err := query.
		Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	return posts, total, err
}

// Update - 포스트 업데이트
func (r *PostRepository) Update(post *Post) error {
	return r.db.Save(post).Error
}

// Delete - 포스트 삭제
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&Post{}, id).Error
}

// AddTag - 포스트에 태그 추가
func (r *PostRepository) AddTag(postID uint, tagID uint) error {
	var post Post
	var tag Tag

	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}
	if err := r.db.First(&tag, tagID).Error; err != nil {
		return err
	}

	return r.db.Model(&post).Association("Tags").Append(&tag)
}

// RemoveTag - 포스트에서 태그 제거
func (r *PostRepository) RemoveTag(postID uint, tagID uint) error {
	var post Post
	var tag Tag

	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}
	if err := r.db.First(&tag, tagID).Error; err != nil {
		return err
	}

	return r.db.Model(&post).Association("Tags").Delete(&tag)
}

// ============================================================================
// Service 레이어
// ============================================================================

type BlogService struct {
	userRepo *UserRepository
	postRepo *PostRepository
	db       *Database
}

func NewBlogService(db *Database) *BlogService {
	return &BlogService{
		userRepo: NewUserRepository(db),
		postRepo: NewPostRepository(db),
		db:       db,
	}
}

// GetUserWithPosts - 사용자와 포스트 함께 조회
func (s *BlogService) GetUserWithPosts(userID uint) (*User, error) {
	var user User
	err := s.db.Preload("Posts", "published = ?", true).
		Preload("Posts.Category").
		First(&user, userID).Error
	return &user, err
}

// GetPopularPosts - 인기 포스트 조회
func (s *BlogService) GetPopularPosts(limit int) ([]Post, error) {
	var posts []Post
	err := s.db.Where("published = ?", true).
		Order("view_count DESC").
		Limit(limit).
		Preload("User").
		Find(&posts).Error
	return posts, err
}

// SearchPosts - 포스트 검색
func (s *BlogService) SearchPosts(keyword string) ([]Post, error) {
	var posts []Post
	searchTerm := "%" + keyword + "%"
	err := s.db.Where("title LIKE ? OR content LIKE ?", searchTerm, searchTerm).
		Where("published = ?", true).
		Preload("User").
		Find(&posts).Error
	return posts, err
}

// ============================================================================
// HTTP Handlers
// ============================================================================

type Handler struct {
	service *BlogService
}

func NewHandler(service *BlogService) *Handler {
	return &Handler{service: service}
}

// User Handlers
func (h *Handler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.userRepo.Create(&user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.userRepo.FindByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func (h *Handler) GetUsers(c *gin.Context) {
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

	offset := (p - 1) * ps

	users, total, err := h.service.userRepo.FindAll(offset, ps)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{
		"users":      users,
		"total":      total,
		"page":       p,
		"page_size":  ps,
		"total_pages": (total + int64(ps) - 1) / int64(ps),
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.userRepo.UpdateFields(id, updates); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	hard := c.Query("hard") == "true"

	var err error
	if hard {
		err = h.service.userRepo.HardDelete(id)
	} else {
		err = h.service.userRepo.Delete(id)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// Post Handlers
func (h *Handler) CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.postRepo.Create(&post); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(201, post)
}

func (h *Handler) GetPost(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.service.postRepo.FindByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, post)
}

func (h *Handler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.service.postRepo.FindBySlug(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	published := c.DefaultQuery("published", "")
	userID := c.DefaultQuery("user_id", "")
	categoryID := c.DefaultQuery("category_id", "")

	var p, ps int
	fmt.Sscanf(page, "%d", &p)
	fmt.Sscanf(pageSize, "%d", &ps)

	if p < 1 {
		p = 1
	}
	if ps < 1 || ps > 100 {
		ps = 10
	}

	offset := (p - 1) * ps

	filters := make(map[string]interface{})
	if published != "" {
		filters["published"] = published == "true"
	}
	if userID != "" {
		var uid uint
		fmt.Sscanf(userID, "%d", &uid)
		filters["user_id"] = uid
	}
	if categoryID != "" {
		var cid uint
		fmt.Sscanf(categoryID, "%d", &cid)
		filters["category_id"] = cid
	}

	posts, total, err := h.service.postRepo.FindAll(filters, offset, ps)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(200, gin.H{
		"posts":      posts,
		"total":      total,
		"page":       p,
		"page_size":  ps,
		"total_pages": (total + int64(ps) - 1) / int64(ps),
	})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	var post Post
	if err := h.service.db.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	post.ID = id
	if err := h.service.postRepo.Update(&post); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(200, post)
}

func (h *Handler) DeletePost(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		c.JSON(400, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.service.postRepo.Delete(id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(200, gin.H{"message": "Post deleted successfully"})
}

// Search Handler
func (h *Handler) SearchPosts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(400, gin.H{"error": "Search keyword is required"})
		return
	}

	posts, err := h.service.SearchPosts(keyword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to search posts"})
		return
	}

	c.JSON(200, gin.H{
		"keyword": keyword,
		"results": posts,
		"count":   len(posts),
	})
}

// Popular Posts Handler
func (h *Handler) GetPopularPosts(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	var l int
	fmt.Sscanf(limit, "%d", &l)

	if l < 1 || l > 50 {
		l = 10
	}

	posts, err := h.service.GetPopularPosts(l)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch popular posts"})
		return
	}

	c.JSON(200, posts)
}

// Advanced Query Examples
func (h *Handler) GetAdvancedQueries(c *gin.Context) {
	examples := []gin.H{
		{
			"name": "Raw SQL",
			"description": "Execute raw SQL queries",
			"example": `db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)`,
		},
		{
			"name": "Joins",
			"description": "Join tables",
			"example": `db.Joins("JOIN posts ON posts.user_id = users.id").Find(&users)`,
		},
		{
			"name": "Subquery",
			"description": "Use subqueries",
			"example": `db.Where("id IN (?)", db.Table("posts").Select("user_id")).Find(&users)`,
		},
		{
			"name": "Aggregation",
			"description": "Count, Sum, Avg, etc.",
			"example": `db.Model(&Post{}).Select("user_id, COUNT(*) as post_count").Group("user_id").Scan(&results)`,
		},
		{
			"name": "Batch Operations",
			"description": "Batch insert/update",
			"example": `db.CreateInBatches(users, 100)`,
		},
		{
			"name": "Hooks",
			"description": "Before/After hooks",
			"example": `func (u *User) BeforeCreate(tx *gorm.DB) error { ... }`,
		},
		{
			"name": "Scopes",
			"description": "Reusable query conditions",
			"example": `db.Scopes(Published, Popular).Find(&posts)`,
		},
	}

	c.JSON(200, examples)
}

// ============================================================================
// Helper Functions
// ============================================================================

func slugify(text string) string {
	// 간단한 slug 생성 (실제로는 더 복잡한 로직 필요)
	result := ""
	for _, char := range text {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			result += string(char)
		} else if char == ' ' {
			result += "-"
		}
	}
	return result
}

// ============================================================================
// Router Setup
// ============================================================================

func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"database": "SQLite",
			"orm": "GORM",
		})
	})

	// Database info
	router.GET("/db/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"database": "SQLite",
			"file": "blog.db",
			"models": []string{"User", "Post", "Category", "Tag", "Comment"},
			"features": []string{
				"Auto Migration",
				"Soft Delete",
				"Associations",
				"Hooks",
				"Transactions",
			},
		})
	})

	// User routes
	users := router.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}

	// Post routes
	posts := router.Group("/posts")
	{
		posts.POST("", handler.CreatePost)
		posts.GET("", handler.GetPosts)
		posts.GET("/:id", handler.GetPost)
		posts.GET("/slug/:slug", handler.GetPostBySlug)
		posts.PUT("/:id", handler.UpdatePost)
		posts.DELETE("/:id", handler.DeletePost)
	}

	// Search and filters
	router.GET("/search", handler.SearchPosts)
	router.GET("/popular", handler.GetPopularPosts)

	// Advanced queries examples
	router.GET("/examples/queries", handler.GetAdvancedQueries)

	return router
}

// ============================================================================
// Main
// ============================================================================

func main() {
	// 데이터베이스 연결
	db, err := NewDatabase(true)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 서비스 초기화
	service := NewBlogService(db)

	// 핸들러 초기화
	handler := NewHandler(service)

	// 라우터 설정
	router := SetupRouter(handler)

	// 서버 시작
	log.Println("🚀 Server starting on :8080")
	log.Println("📊 Database: SQLite (blog.db)")
	log.Println("🔧 ORM: GORM with auto-migration")

	if err := router.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}