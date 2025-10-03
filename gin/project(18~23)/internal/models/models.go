package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `json:"role" gorm:"default:user"`
	Active    bool           `json:"active" gorm:"default:true"`
	Posts     []Post         `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Post represents a blog post
type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content"`
	Slug        string         `json:"slug" gorm:"uniqueIndex"`
	Published   bool           `json:"published" gorm:"default:false"`
	UserID      uint           `json:"user_id"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tags        []Tag          `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Tag represents a post tag
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null"`
	Posts     []Post    `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	CreatedAt time.Time `json:"created_at"`
}

// RefreshToken stores refresh tokens
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditLog stores audit logs
type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName sets the table name for User
func (User) TableName() string {
	return "users"
}

// TableName sets the table name for Post
func (Post) TableName() string {
	return "posts"
}