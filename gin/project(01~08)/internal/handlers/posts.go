package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gin-project/internal/models"
)

// In-memory storage (replace with database in production)
var posts = map[int]*models.Post{
	1: {
		ID:       1,
		Title:    "Welcome to Our Blog Platform",
		Content:  "This is our first blog post. Welcome to the community!",
		AuthorID: 1,
		Category: "Announcement",
		Tags:     []string{"welcome", "first-post"},
		ImageURL: "/static/images/blog1.jpg",
		ViewCount: 42,
		IsPublished: true,
		CreatedAt: time.Now().AddDate(0, -1, 0),
		UpdatedAt: time.Now().AddDate(0, -1, 0),
	},
	2: {
		ID:       2,
		Title:    "Getting Started with Gin Framework",
		Content:  "Gin is a high-performance HTTP web framework written in Go. In this post, we'll explore its core features...",
		AuthorID: 1,
		Category: "Tutorial",
		Tags:     []string{"golang", "gin", "tutorial"},
		ImageURL: "/static/images/blog2.jpg",
		ViewCount: 128,
		IsPublished: true,
		CreatedAt: time.Now().AddDate(0, 0, -7),
		UpdatedAt: time.Now().AddDate(0, 0, -7),
	},
}

var nextPostID = 3

// GetPosts returns all posts with pagination
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	var postList []*models.Post
	for _, post := range posts {
		if post.IsPublished {
			if category == "" || post.Category == category {
				// Add author info
				if author, exists := users[getUsernameByID(post.AuthorID)]; exists {
					post.Author = author
				}
				postList = append(postList, post)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": postList,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       len(postList),
			"total_pages": (len(postList) + limit - 1) / limit,
		},
	})
}

// GetPost returns a single post by ID
func GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	post, exists := posts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	// Increment view count
	post.ViewCount++

	// Add author info
	if author, exists := users[getUsernameByID(post.AuthorID)]; exists {
		post.Author = author
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// CreatePost creates a new blog post
func CreatePost(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")

	post := &models.Post{
		ID:          nextPostID,
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    userID.(int),
		Category:    req.Category,
		Tags:        req.Tags,
		ViewCount:   0,
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	nextPostID++

	posts[post.ID] = post

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// UpdatePost updates an existing post
func UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	post, exists := posts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	// Check if user is author or admin
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if post.AuthorID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to update this post",
		})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Category != nil {
		post.Category = *req.Category
	}
	if req.Tags != nil {
		post.Tags = req.Tags
	}

	post.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

// DeletePost deletes a post
func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	post, exists := posts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	// Check if user is author or admin
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if post.AuthorID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this post",
		})
		return
	}

	delete(posts, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

// Helper function
func getUsernameByID(id int) string {
	for _, user := range users {
		if user.ID == id {
			return user.Username
		}
	}
	return ""
}
