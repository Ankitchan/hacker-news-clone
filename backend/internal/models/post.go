package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	URL          sql.NullString `json:"url,omitempty"`          // Optional URL for link posts
	Text         sql.NullString `json:"text,omitempty"`         // Optional text for text posts
	UserID       int            `json:"user_id"`
	Username     string         `json:"username"`               // Joined from users table
	Points       int            `json:"points"`                 // Net votes (upvotes - downvotes)
	CommentCount int            `json:"comment_count"`          // Total number of comments
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type PostCreate struct {
	Title  string  `json:"title"`
	URL    *string `json:"url,omitempty"`
	Text   *string `json:"text,omitempty"`
	UserID int     `json:"-"` // Set from auth context, not from request body
}

type PostUpdate struct {
	Title *string `json:"title,omitempty"`
	URL   *string `json:"url,omitempty"`
	Text  *string `json:"text,omitempty"`
}

type PostList struct {
	Posts      []Post `json:"posts"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}
