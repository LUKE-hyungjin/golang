package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gin-project/internal/models"
)

// In-memory storage (replace with database in production)
var users = map[string]*models.User{
	"admin": {
		ID:        1,
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "password123",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	"user": {
		ID:        2,
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password123",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var nextUserID = 3

// Register handles user registration
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if username already exists
	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// Create new user
	user := &models.User{
		ID:        nextUserID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	nextUserID++

	users[req.Username] = user

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Find user
	user, exists := users[req.Username]
	if !exists || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate token (simplified - use JWT in production)
	token := "valid-token-123"
	if user.Role == "admin" {
		token = "admin-token-456"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// GetProfile returns current user's profile
func GetProfile(c *gin.Context) {
	username, _ := c.Get("username")
	user, exists := users[username.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
