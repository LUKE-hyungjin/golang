package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
}

func main() {
	r := gin.Default()

	// Basic routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Path param
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// Query param
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		c.JSON(http.StatusOK, gin.H{"q": q})
	})

	// POST JSON binding
	r.POST("/users", func(c *gin.Context) {
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"name": body.Name, "email": body.Email})
	})

	// PUT
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"updated": id})
	})

	// DELETE
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Status(http.StatusNoContent)
		_ = id
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
