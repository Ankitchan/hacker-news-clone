-- Drop triggers
DROP TRIGGER IF EXISTS trigger_notify_on_comment_edit ON comments;
DROP TRIGGER IF EXISTS trigger_notify_parent_comment_owner_on_reply ON comments;
DROP TRIGGER IF EXISTS trigger_notify_post_owner_on_comment ON comments;

-- Drop functions
DROP FUNCTION IF EXISTS notify_on_comment_edit();
DROP FUNCTION IF EXISTS notify_parent_comment_owner_on_reply();
DROP FUNCTION IF EXISTS notify_post_owner_on_comment();

-- Drop indexes
DROP INDEX IF EXISTS idx_notifications_user_unread;
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id;

-- Drop table
DROP TABLE IF EXISTS notifications;
