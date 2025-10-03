package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gin-project/internal/models"
)

// In-memory storage
var comments = map[int]*models.Comment{
	1: {
		ID:        1,
		PostID:    1,
		AuthorID:  2,
		Content:   "Great post! Looking forward to more content.",
		CreatedAt: time.Now().AddDate(0, 0, -5),
		UpdatedAt: time.Now().AddDate(0, 0, -5),
	},
	2: {
		ID:        2,
		PostID:    1,
		AuthorID:  1,
		Content:   "Thanks for your support!",
		ParentID:  func() *int { i := 1; return &i }(),
		CreatedAt: time.Now().AddDate(0, 0, -4),
		UpdatedAt: time.Now().AddDate(0, 0, -4),
	},
}

var nextCommentID = 3

// GetComments returns all comments for a post
func GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	// Verify post exists
	if _, exists := posts[postID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	var commentList []*models.Comment
	for _, comment := range comments {
		if comment.PostID == postID {
			// Add author info
			if author, exists := users[getUsernameByID(comment.AuthorID)]; exists {
				comment.Author = author
			}
			commentList = append(commentList, comment)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": commentList,
		"total":    len(commentList),
	})
}

// CreateComment creates a new comment
func CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Verify post exists
	if _, exists := posts[req.PostID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	userID, _ := c.Get("user_id")

	comment := &models.Comment{
		ID:        nextCommentID,
		PostID:    req.PostID,
		AuthorID:  userID.(int),
		Content:   req.Content,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	nextCommentID++

	comments[comment.ID] = comment

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// DeleteComment deletes a comment
func DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	comment, exists := comments[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comment not found",
		})
		return
	}

	// Check if user is author or admin
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if comment.AuthorID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to delete this comment",
		})
		return
	}

	delete(comments, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
