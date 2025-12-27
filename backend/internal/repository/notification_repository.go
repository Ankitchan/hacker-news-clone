package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// GetByUserID retrieves all notifications for a specific user with pagination
func (r *NotificationRepository) GetByUserID(userID, limit, offset int) ([]models.Notification, int, int, error) {
	// Get total count and unread count
	var totalCount, unreadCount int
	countQuery := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_read = FALSE) as unread
		FROM notifications
		WHERE user_id = $1
	`
	err := r.db.QueryRow(countQuery, userID).Scan(&totalCount, &unreadCount)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get notification counts: %w", err)
	}

	// Get notifications with actor username
	query := `
		SELECT
			n.id, n.user_id, n.actor_id, n.post_id, n.comment_id,
			n.notification_type, n.message, n.is_read, n.created_at,
			u.username as actor_username
		FROM notifications n
		JOIN users u ON n.actor_id = u.id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.ActorID,
			&n.PostID,
			&n.CommentID,
			&n.NotificationType,
			&n.Message,
			&n.IsRead,
			&n.CreatedAt,
			&n.ActorUsername,
		)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, totalCount, unreadCount, nil
}

// GetUnreadByUserID retrieves only unread notifications for a specific user
func (r *NotificationRepository) GetUnreadByUserID(userID, limit, offset int) ([]models.Notification, int, error) {
	// Get unread count
	var unreadCount int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	err := r.db.QueryRow(countQuery, userID).Scan(&unreadCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	// Get unread notifications
	query := `
		SELECT
			n.id, n.user_id, n.actor_id, n.post_id, n.comment_id,
			n.notification_type, n.message, n.is_read, n.created_at,
			u.username as actor_username
		FROM notifications n
		JOIN users u ON n.actor_id = u.id
		WHERE n.user_id = $1 AND n.is_read = FALSE
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.ActorID,
			&n.PostID,
			&n.CommentID,
			&n.NotificationType,
			&n.Message,
			&n.IsRead,
			&n.CreatedAt,
			&n.ActorUsername,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, unreadCount, nil
}

// MarkAsRead marks specific notifications as read
func (r *NotificationRepository) MarkAsRead(userID int, notificationIDs []int) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	// Build PostgreSQL array literal string
	arrayStr := "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(notificationIDs)), ","), "[]") + "}"

	// Build the query with placeholders
	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = $1 AND id = ANY($2::int[])
	`

	result, err := r.db.Exec(query, userID, arrayStr)
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no notifications were marked as read")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(userID int) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	return nil
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(userID, notificationID int) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`

	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}
