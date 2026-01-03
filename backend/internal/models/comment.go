package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID          int           `json:"id"`
	PostID      int           `json:"post_id"`
	UserID      int           `json:"user_id"`
	Username    string        `json:"username"`            // Joined from users table
	ParentID    sql.NullInt64 `json:"parent_id,omitempty"` // Null for top-level comments
	Text        string        `json:"text"`
	Points      int           `json:"points"`       // Net votes (for database compatibility)
	TotalPoints int           `json:"total_points"` // Cumulative points (this comment + all descendants)
	Depth       int           `json:"depth"`        // 0 for top-level, increases with nesting
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Replies     []Comment     `json:"replies,omitempty"` // For threaded display
}

type CommentCreate struct {
	PostID   int    `json:"post_id"`
	ParentID *int   `json:"parent_id,omitempty"`
	Text     string `json:"text"`
	UserID   int    `json:"-"` // Set from auth context
}

type CommentUpdate struct {
	Text string `json:"text"`
}

type CommentTree struct {
	Comments []Comment `json:"comments"`
}
