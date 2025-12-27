package models

import (
	"database/sql"
	"time"
)

type NotificationType string

const (
	NotificationCommentOnPost NotificationType = "comment_on_post"
	NotificationCommentReply  NotificationType = "comment_reply"
	NotificationCommentEdit   NotificationType = "comment_edit"
)

type Notification struct {
	ID               int              `json:"id"`
	UserID           int              `json:"user_id"`
	ActorID          int              `json:"actor_id"`
	ActorUsername    string           `json:"actor_username"` // Joined from users table
	PostID           sql.NullInt64    `json:"post_id,omitempty"`
	CommentID        sql.NullInt64    `json:"comment_id,omitempty"`
	NotificationType NotificationType `json:"notification_type"`
	Message          string           `json:"message"`
	IsRead           bool             `json:"is_read"`
	CreatedAt        time.Time        `json:"created_at"`
}

type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	TotalCount    int            `json:"total_count"`
	UnreadCount   int            `json:"unread_count"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
}

type MarkAsReadRequest struct {
	NotificationIDs []int `json:"notification_ids"`
}
