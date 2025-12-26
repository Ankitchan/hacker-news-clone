package models

import (
	"time"
)

type VoteType int

const (
	VoteUp   VoteType = 1
	VoteDown VoteType = -1
)

type Vote struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PostID    *int      `json:"post_id,omitempty"`    // Either PostID or CommentID will be set
	CommentID *int      `json:"comment_id,omitempty"` // Either PostID or CommentID will be set
	VoteType  VoteType  `json:"vote_type"`            // 1 for upvote, -1 for downvote
	CreatedAt time.Time `json:"created_at"`
}

type VoteCreate struct {
	PostID    *int     `json:"post_id,omitempty"`
	CommentID *int     `json:"comment_id,omitempty"`
	VoteType  VoteType `json:"vote_type"` // 1 for upvote, -1 for downvote
	UserID    int      `json:"-"`         // Set from auth context
}

type VoteResponse struct {
	Success bool `json:"success"`
	Points  int  `json:"points"` // Updated total points
}
