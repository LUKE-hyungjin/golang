package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
}

type PatchUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users  = make(map[string]User)
	mu     sync.RWMutex
	nextID int64 = 1
)

func getNextID() string {
	id := strconv.FormatInt(nextID, 10)
	nextID++
	return id
}

func main() {
	r := gin.Default()

	// Basic routes
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Users - list
	r.GET("/users", func(c *gin.Context) {
		mu.RLock()
		defer mu.RUnlock()

		list := make([]User, 0, len(users))
		for _, u := range users {
			list = append(list, u)
		}
		c.JSON(http.StatusOK, list)
	})

	// Users - get by id
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		mu.RLock()
		u, ok := users[id]
		mu.RUnlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	// Query param
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")
		c.JSON(http.StatusOK, gin.H{"q": q})
	})

	// Users - create
	r.POST("/users", func(c *gin.Context) {
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		id := getNextID()
		u := User{ID: id, Name: body.Name, Email: body.Email}
		users[id] = u
		mu.Unlock()

		c.JSON(http.StatusCreated, gin.H{"id": id, "user": u})
	})

	// Users - update (full)
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		u, ok := users[id]
		if !ok {
			mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		u.Name = body.Name
		u.Email = body.Email
		users[id] = u
		mu.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "User updated", "user": u})
	})

	// Users - partial update
	r.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var body PatchUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mu.Lock()
		u, ok := users[id]
		if !ok {
			mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if body.Name != nil {
			u.Name = *body.Name
		}
		if body.Email != nil {
			u.Email = *body.Email
		}
		users[id] = u
		mu.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "User partially updated", "user": u})
	})

	// Users - delete
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		mu.Lock()
		if _, ok := users[id]; !ok {
			mu.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		delete(users, id)
		mu.Unlock()
		c.Status(http.StatusNoContent)
	})

	if err := r.Run(":3001"); err != nil {
		panic(err)
	}
}
