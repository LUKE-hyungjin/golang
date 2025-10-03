package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id" binding:"required"`
	AuthorID  int       `json:"author_id"`
	Author    *User     `json:"author,omitempty"`
	Content   string    `json:"content" binding:"required,min=1,max=1000"`
	ParentID  *int      `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	PostID   int    `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID *int   `json:"parent_id,omitempty"`
}
