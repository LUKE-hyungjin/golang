package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users  []User
	nextID = 1
)

func findUserIndexByID(id string) int {
	for i, u := range users {
		if u.ID == id {
			return i
		}
	}
	return -1
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idx := findUserIndexByID(id)
		if idx == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, users[idx])
	})

	r.POST("/users", func(c *gin.Context) {
		var body User
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		id := strconv.Itoa(nextID)
		nextID++
		u := User{ID: id, Name: body.Name, Email: body.Email}
		users = append(users, u)
		c.JSON(http.StatusCreated, gin.H{"id": id, "user": u})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var body User
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		idx := findUserIndexByID(id)
		if idx == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		u := users[idx]
		u.Name = body.Name
		u.Email = body.Email
		users[idx] = u
		c.JSON(http.StatusOK, gin.H{"user": u})
	})

	r.DELETE("users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idx := findUserIndexByID(id)
		if idx == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		users = append(users[:idx], users[idx+1:]...)
		c.Status(http.StatusNoContent)
	})

	if err := r.Run(":3002"); err != nil {
		panic(err)
	}
}
