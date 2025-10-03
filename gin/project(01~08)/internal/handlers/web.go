package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gin-project/internal/models"
)

// RenderHome renders the home page
func RenderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Blog Community Platform",
		"year":  time.Now().Year(),
	})
}

// RenderPosts renders the posts listing page
func RenderPosts(c *gin.Context) {
	var postList []*models.Post
	for _, post := range posts {
		if post.IsPublished {
			if author, exists := users[getUsernameByID(post.AuthorID)]; exists {
				post.Author = author
			}
			postList = append(postList, post)
		}
	}

	c.HTML(http.StatusOK, "posts.html", gin.H{
		"title": "Blog Posts",
		"posts": postList,
	})
}

// RenderPost renders a single post page
func RenderPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Post Not Found",
		})
		return
	}

	post, exists := posts[id]
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Post Not Found",
		})
		return
	}

	// Get author info
	if author, exists := users[getUsernameByID(post.AuthorID)]; exists {
		post.Author = author
	}

	// Get comments
	var commentList []*models.Comment
	for _, comment := range comments {
		if comment.PostID == id {
			if author, exists := users[getUsernameByID(comment.AuthorID)]; exists {
				comment.Author = author
			}
			commentList = append(commentList, comment)
		}
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":    post.Title,
		"post":     post,
		"comments": commentList,
	})
}

// RenderAdminDashboard renders the admin dashboard
func RenderAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"title": "Admin Dashboard",
		"stats": gin.H{
			"total_users":    len(users),
			"total_posts":    len(posts),
			"total_comments": len(comments),
		},
	})
}

// RenderLogin renders the login page
func RenderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// RenderRegister renders the registration page
func RenderRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}
