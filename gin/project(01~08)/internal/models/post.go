package models

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" binding:"required,min=3,max=200"`
	Content   string    `json:"content" binding:"required"`
	AuthorID  int       `json:"author_id"`
	Author    *User     `json:"author,omitempty"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	ImageURL  string    `json:"image_url"`
	ViewCount int       `json:"view_count"`
	IsPublished bool    `json:"is_published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required,min=3,max=200"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title    *string  `json:"title,omitempty"`
	Content  *string  `json:"content,omitempty"`
	Category *string  `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}
