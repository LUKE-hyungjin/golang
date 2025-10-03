package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================================
// 모델 정의 (버전별 스키마)
// ============================================================================

// V1 Models (초기 버전)
type UserV1 struct {
	ID        uint      `gorm:"primarykey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Username  string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// V2 Models (필드 추가)
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Username    string         `gorm:"uniqueIndex;not null" json:"username"`
	Name        string         `json:"name"`
	Bio         string         `gorm:"type:text" json:"bio"`
	Avatar      string         `json:"avatar"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IsAdmin     bool           `gorm:"default:false" json:"is_admin"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Posts       []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}

type Post struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Excerpt     string         `json:"excerpt"`
	CoverImage  string         `json:"cover_image"`
	Published   bool           `gorm:"default:false;index" json:"published"`
	PublishedAt *time.Time     `json:"published_at"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`
	UserID      uint           `json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID  *uint          `json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag          `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Posts       []Post         `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex;not null" json:"slug"`
	Posts     []Post         `gorm:"many2many:post_tags;" json:"posts,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ============================================================================
// 마이그레이션 시스템
// ============================================================================

type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time
}

type MigrationFunc struct {
	Version string
	Name    string
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

type Migrator struct {
	db         *gorm.DB
	migrations []MigrationFunc
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []MigrationFunc{},
	}
}

func (m *Migrator) AddMigration(migration MigrationFunc) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) Migrate() error {
	// 마이그레이션 테이블 생성
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	for _, migration := range m.migrations {
		var count int64
		m.db.Model(&Migration{}).Where("version = ?", migration.Version).Count(&count)

		if count == 0 {
			log.Printf("🔄 Applying migration: %s - %s", migration.Version, migration.Name)

			// 트랜잭션으로 마이그레이션 실행
			err := m.db.Transaction(func(tx *gorm.DB) error {
				if err := migration.Up(tx); err != nil {
					return err
				}

				// 마이그레이션 기록
				record := Migration{
					Version:   migration.Version,
					Name:      migration.Name,
					AppliedAt: time.Now(),
				}
				return tx.Create(&record).Error
			})

			if err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Version, err)
			}

			log.Printf("✅ Migration %s completed", migration.Version)
		}
	}

	return nil
}

func (m *Migrator) Rollback(version string) error {
	// 특정 버전으로 롤백
	var migration Migration
	if err := m.db.Where("version = ?", version).First(&migration).Error; err != nil {
		return fmt.Errorf("migration %s not found", version)
	}

	// 해당 마이그레이션 찾기
	var targetMigration *MigrationFunc
	for _, mig := range m.migrations {
		if mig.Version == version {
			targetMigration = &mig
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration function for %s not found", version)
	}

	log.Printf("🔄 Rolling back migration: %s", version)

	// 트랜잭션으로 롤백 실행
	err := m.db.Transaction(func(tx *gorm.DB) error {
		if err := targetMigration.Down(tx); err != nil {
			return err
		}

		// 마이그레이션 기록 삭제
		return tx.Delete(&migration).Error
	})

	if err != nil {
		return fmt.Errorf("rollback %s failed: %w", version, err)
	}

	log.Printf("✅ Rollback %s completed", version)
	return nil
}

func (m *Migrator) Status() ([]Migration, error) {
	var migrations []Migration
	err := m.db.Order("applied_at desc").Find(&migrations).Error
	return migrations, err
}

// ============================================================================
// 마이그레이션 정의
// ============================================================================

func GetMigrations() []MigrationFunc {
	return []MigrationFunc{
		{
			Version: "001_create_users_table",
			Name:    "Create users table",
			Up: func(db *gorm.DB) error {
				return db.AutoMigrate(&UserV1{})
			},
			Down: func(db *gorm.DB) error {
				return db.Migrator().DropTable(&UserV1{})
			},
		},
		{
			Version: "002_add_user_fields",
			Name:    "Add name, bio, avatar fields to users",
			Up: func(db *gorm.DB) error {
				// 새 필드 추가
				if !db.Migrator().HasColumn(&User{}, "name") {
					db.Migrator().AddColumn(&User{}, "name")
				}
				if !db.Migrator().HasColumn(&User{}, "bio") {
					db.Migrator().AddColumn(&User{}, "bio")
				}
				if !db.Migrator().HasColumn(&User{}, "avatar") {
					db.Migrator().AddColumn(&User{}, "avatar")
				}
				if !db.Migrator().HasColumn(&User{}, "is_admin") {
					db.Migrator().AddColumn(&User{}, "is_admin")
				}
				return nil
			},
			Down: func(db *gorm.DB) error {
				db.Migrator().DropColumn(&User{}, "name")
				db.Migrator().DropColumn(&User{}, "bio")
				db.Migrator().DropColumn(&User{}, "avatar")
				db.Migrator().DropColumn(&User{}, "is_admin")
				return nil
			},
		},
		{
			Version: "003_create_posts_categories_tags",
			Name:    "Create posts, categories, and tags tables",
			Up: func(db *gorm.DB) error {
				return db.AutoMigrate(&Category{}, &Tag{}, &Post{})
			},
			Down: func(db *gorm.DB) error {
				db.Migrator().DropTable("post_tags")
				db.Migrator().DropTable(&Post{})
				db.Migrator().DropTable(&Tag{})
				db.Migrator().DropTable(&Category{})
				return nil
			},
		},
		{
			Version: "004_add_soft_delete",
			Name:    "Add soft delete to all tables",
			Up: func(db *gorm.DB) error {
				// DeletedAt 필드는 AutoMigrate로 자동 추가됨
				return db.AutoMigrate(&User{}, &Post{}, &Category{}, &Tag{})
			},
			Down: func(db *gorm.DB) error {
				// Soft delete 필드 제거는 권장하지 않음
				return nil
			},
		},
		{
			Version: "005_add_post_metrics",
			Name:    "Add like_count and published_at to posts",
			Up: func(db *gorm.DB) error {
				if !db.Migrator().HasColumn(&Post{}, "like_count") {
					db.Migrator().AddColumn(&Post{}, "like_count")
				}
				if !db.Migrator().HasColumn(&Post{}, "published_at") {
					db.Migrator().AddColumn(&Post{}, "published_at")
				}
				if !db.Migrator().HasColumn(&Post{}, "excerpt") {
					db.Migrator().AddColumn(&Post{}, "excerpt")
				}
				if !db.Migrator().HasColumn(&Post{}, "cover_image") {
					db.Migrator().AddColumn(&Post{}, "cover_image")
				}
				return nil
			},
			Down: func(db *gorm.DB) error {
				db.Migrator().DropColumn(&Post{}, "like_count")
				db.Migrator().DropColumn(&Post{}, "published_at")
				db.Migrator().DropColumn(&Post{}, "excerpt")
				db.Migrator().DropColumn(&Post{}, "cover_image")
				return nil
			},
		},
	}
}

// ============================================================================
// 시드 데이터 생성
// ============================================================================

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Seed() error {
	log.Println("🌱 Starting seed process...")

	// 1. Categories 생성
	if err := s.seedCategories(); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	// 2. Tags 생성
	if err := s.seedTags(); err != nil {
		return fmt.Errorf("failed to seed tags: %w", err)
	}

	// 3. Users 생성
	if err := s.seedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// 4. Posts 생성
	if err := s.seedPosts(); err != nil {
		return fmt.Errorf("failed to seed posts: %w", err)
	}

	log.Println("✅ Seed process completed!")
	return nil
}

func (s *Seeder) seedCategories() error {
	categories := []Category{
		{Name: "Technology", Slug: "technology", Description: "Tech articles and tutorials"},
		{Name: "Programming", Slug: "programming", Description: "Programming languages and concepts"},
		{Name: "Web Development", Slug: "web-dev", Description: "Web development topics"},
		{Name: "Mobile Development", Slug: "mobile-dev", Description: "Mobile app development"},
		{Name: "DevOps", Slug: "devops", Description: "DevOps and infrastructure"},
		{Name: "Database", Slug: "database", Description: "Database design and optimization"},
		{Name: "Security", Slug: "security", Description: "Cybersecurity topics"},
		{Name: "AI/ML", Slug: "ai-ml", Description: "Artificial Intelligence and Machine Learning"},
	}

	for _, cat := range categories {
		var count int64
		s.db.Model(&Category{}).Where("slug = ?", cat.Slug).Count(&count)
		if count == 0 {
			if err := s.db.Create(&cat).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("✅ Seeded %d categories", len(categories))
	return nil
}

func (s *Seeder) seedTags() error {
	tags := []Tag{
		{Name: "Go", Slug: "go"},
		{Name: "Python", Slug: "python"},
		{Name: "JavaScript", Slug: "javascript"},
		{Name: "Docker", Slug: "docker"},
		{Name: "Kubernetes", Slug: "kubernetes"},
		{Name: "React", Slug: "react"},
		{Name: "Vue", Slug: "vue"},
		{Name: "Node.js", Slug: "nodejs"},
		{Name: "PostgreSQL", Slug: "postgresql"},
		{Name: "MongoDB", Slug: "mongodb"},
		{Name: "Redis", Slug: "redis"},
		{Name: "AWS", Slug: "aws"},
		{Name: "GCP", Slug: "gcp"},
		{Name: "CI/CD", Slug: "ci-cd"},
		{Name: "Testing", Slug: "testing"},
	}

	for _, tag := range tags {
		var count int64
		s.db.Model(&Tag{}).Where("slug = ?", tag.Slug).Count(&count)
		if count == 0 {
			if err := s.db.Create(&tag).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("✅ Seeded %d tags", len(tags))
	return nil
}

func (s *Seeder) seedUsers() error {
	// Admin user
	admin := User{
		Email:    "admin@example.com",
		Username: "admin",
		Name:     "Admin User",
		Bio:      "System Administrator",
		IsAdmin:  true,
		IsActive: true,
	}

	var count int64
	s.db.Model(&User{}).Where("email = ?", admin.Email).Count(&count)
	if count == 0 {
		if err := s.db.Create(&admin).Error; err != nil {
			return err
		}
	}

	// Regular users with faker
	for i := 0; i < 10; i++ {
		user := User{
			Email:    faker.Email(),
			Username: faker.Username(),
			Name:     faker.Name(),
			Bio:      faker.Sentence(),
			Avatar:   fmt.Sprintf("https://i.pravatar.cc/150?img=%d", i+1),
			IsActive: true,
		}

		var count int64
		s.db.Model(&User{}).Where("email = ?", user.Email).Count(&count)
		if count == 0 {
			if err := s.db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user: %v", err)
				continue
			}
		}
	}

	var total int64
	s.db.Model(&User{}).Count(&total)
	log.Printf("✅ Seeded users (total: %d)", total)
	return nil
}

func (s *Seeder) seedPosts() error {
	var users []User
	var categories []Category
	var tags []Tag

	s.db.Find(&users)
	s.db.Find(&categories)
	s.db.Find(&tags)

	if len(users) == 0 || len(categories) == 0 {
		return fmt.Errorf("users or categories not found")
	}

	// 각 사용자마다 포스트 생성
	for _, user := range users {
		numPosts := rand.Intn(5) + 1 // 1-5개 포스트

		for i := 0; i < numPosts; i++ {
			title := faker.Sentence()
			content := faker.Paragraph()
			published := rand.Float32() > 0.3 // 70% 확률로 published

			post := Post{
				Title:      title,
				Content:    content,
				Excerpt:    truncateString(content, 150),
				Slug:       slugify(title) + fmt.Sprintf("-%d", time.Now().Unix()),
				Published:  published,
				ViewCount:  rand.Intn(1000),
				LikeCount:  rand.Intn(100),
				UserID:     user.ID,
				CategoryID: &categories[rand.Intn(len(categories))].ID,
				CoverImage: fmt.Sprintf("https://picsum.photos/800/400?random=%d", rand.Intn(1000)),
			}

			if published {
				now := time.Now()
				post.PublishedAt = &now
			}

			if err := s.db.Create(&post).Error; err != nil {
				log.Printf("Failed to create post: %v", err)
				continue
			}

			// 랜덤 태그 추가 (1-3개)
			numTags := rand.Intn(3) + 1
			selectedTags := make([]Tag, 0)
			for j := 0; j < numTags && j < len(tags); j++ {
				selectedTags = append(selectedTags, tags[rand.Intn(len(tags))])
			}

			if len(selectedTags) > 0 {
				s.db.Model(&post).Association("Tags").Append(selectedTags)
			}
		}
	}

	var total int64
	s.db.Model(&Post{}).Count(&total)
	log.Printf("✅ Seeded posts (total: %d)", total)
	return nil
}

func (s *Seeder) Clean() error {
	log.Println("🧹 Cleaning database...")

	// 역순으로 삭제 (외래키 제약 고려)
	tables := []interface{}{
		&Post{},
		&Tag{},
		&Category{},
		&User{},
	}

	for _, table := range tables {
		if err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			return fmt.Errorf("failed to clean table: %w", err)
		}
	}

	// Many-to-many 중간 테이블 직접 삭제
	s.db.Exec("DELETE FROM post_tags")

	log.Println("✅ Database cleaned!")
	return nil
}

// Load seed data from JSON file
func (s *Seeder) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	var seedData struct {
		Users      []User     `json:"users"`
		Categories []Category `json:"categories"`
		Tags       []Tag      `json:"tags"`
		Posts      []Post     `json:"posts"`
	}

	if err := json.Unmarshal(data, &seedData); err != nil {
		return fmt.Errorf("failed to unmarshal seed data: %w", err)
	}

	// Insert data in order
	for _, user := range seedData.Users {
		s.db.Create(&user)
	}
	for _, cat := range seedData.Categories {
		s.db.Create(&cat)
	}
	for _, tag := range seedData.Tags {
		s.db.Create(&tag)
	}
	for _, post := range seedData.Posts {
		s.db.Create(&post)
	}

	log.Printf("✅ Loaded seed data from %s", filename)
	return nil
}

// Export current data to JSON file
func (s *Seeder) ExportToFile(filename string) error {
	var users []User
	var categories []Category
	var tags []Tag
	var posts []Post

	s.db.Find(&users)
	s.db.Find(&categories)
	s.db.Find(&tags)
	s.db.Find(&posts)

	seedData := map[string]interface{}{
		"users":      users,
		"categories": categories,
		"tags":       tags,
		"posts":      posts,
	}

	data, err := json.MarshalIndent(seedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("✅ Exported data to %s", filename)
	return nil
}

// ============================================================================
// HTTP Handlers
// ============================================================================

type MigrationHandler struct {
	migrator *Migrator
	seeder   *Seeder
}

func NewMigrationHandler(migrator *Migrator, seeder *Seeder) *MigrationHandler {
	return &MigrationHandler{
		migrator: migrator,
		seeder:   seeder,
	}
}

func (h *MigrationHandler) GetStatus(c *gin.Context) {
	migrations, err := h.migrator.Status()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get migration status"})
		return
	}

	c.JSON(200, gin.H{
		"applied_migrations": migrations,
		"total":             len(migrations),
	})
}

func (h *MigrationHandler) RunMigrations(c *gin.Context) {
	if err := h.migrator.Migrate(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Migrations completed successfully"})
}

func (h *MigrationHandler) Rollback(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		c.JSON(400, gin.H{"error": "Version is required"})
		return
	}

	if err := h.migrator.Rollback(version); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Rolled back to %s", version)})
}

func (h *MigrationHandler) Seed(c *gin.Context) {
	if err := h.seeder.Seed(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Database seeded successfully"})
}

func (h *MigrationHandler) Clean(c *gin.Context) {
	if err := h.seeder.Clean(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Database cleaned successfully"})
}

func (h *MigrationHandler) Reset(c *gin.Context) {
	// Clean and reseed
	if err := h.seeder.Clean(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := h.seeder.Seed(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Database reset successfully"})
}

func (h *MigrationHandler) Export(c *gin.Context) {
	filename := c.DefaultQuery("file", "seed_data.json")

	// Ensure file is in current directory
	filename = filepath.Base(filename)

	if err := h.seeder.ExportToFile(filename); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Data exported successfully",
		"file":    filename,
	})
}

func (h *MigrationHandler) Import(c *gin.Context) {
	filename := c.DefaultQuery("file", "seed_data.json")

	// Ensure file is in current directory
	filename = filepath.Base(filename)

	if err := h.seeder.LoadFromFile(filename); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Data imported successfully",
		"file":    filename,
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

func slugify(text string) string {
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

func truncateString(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[:length] + "..."
}

// ============================================================================
// Router Setup
// ============================================================================

func SetupRouter(handler *MigrationHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// Migration routes
	migrations := router.Group("/migrations")
	{
		migrations.GET("/status", handler.GetStatus)
		migrations.POST("/run", handler.RunMigrations)
		migrations.POST("/rollback/:version", handler.Rollback)
	}

	// Seed routes
	seed := router.Group("/seed")
	{
		seed.POST("/run", handler.Seed)
		seed.POST("/clean", handler.Clean)
		seed.POST("/reset", handler.Reset)
		seed.POST("/export", handler.Export)
		seed.POST("/import", handler.Import)
	}

	// Info route
	router.GET("/info", func(c *gin.Context) {
		var userCount, postCount, categoryCount, tagCount int64
		handler.seeder.db.Model(&User{}).Count(&userCount)
		handler.seeder.db.Model(&Post{}).Count(&postCount)
		handler.seeder.db.Model(&Category{}).Count(&categoryCount)
		handler.seeder.db.Model(&Tag{}).Count(&tagCount)

		c.JSON(200, gin.H{
			"database": "SQLite",
			"file":     "blog.db",
			"stats": gin.H{
				"users":      userCount,
				"posts":      postCount,
				"categories": categoryCount,
				"tags":       tagCount,
			},
		})
	})

	return router
}

// ============================================================================
// Main
// ============================================================================

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Database connection
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Initialize migrator and seeder
	migrator := NewMigrator(db)

	// Add all migrations
	for _, migration := range GetMigrations() {
		migrator.AddMigration(migration)
	}

	seeder := NewSeeder(db)

	// Run migrations on startup
	if err := migrator.Migrate(); err != nil {
		log.Printf("Migration failed: %v", err)
	}

	// Initialize handler
	handler := NewMigrationHandler(migrator, seeder)

	// Setup router
	router := SetupRouter(handler)

	// Start server
	log.Println("🚀 Server starting on :8080")
	log.Println("📊 Database: SQLite (blog.db)")
	log.Println("🔄 Migrations: Ready")
	log.Println("🌱 Seeder: Ready")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}