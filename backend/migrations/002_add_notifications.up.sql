-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL CHECK (notification_type IN ('comment_on_post', 'comment_reply', 'comment_edit')),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Ensure either post_id or comment_id is set
    CONSTRAINT check_notification_target CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR
        (post_id IS NULL AND comment_id IS NOT NULL)
    )
);

-- Create indexes for better query performance
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;

-- Create function to create notification for comment on post
CREATE OR REPLACE FUNCTION notify_post_owner_on_comment()
RETURNS TRIGGER AS $$
DECLARE
    post_owner_id INTEGER;
    actor_username TEXT;
    post_title TEXT;
BEGIN
    -- Get the post owner
    SELECT user_id, title INTO post_owner_id, post_title
    FROM posts
    WHERE id = NEW.post_id;

    -- Get the actor's username
    SELECT username INTO actor_username
    FROM users
    WHERE id = NEW.user_id;

    -- Only create notification if the comment author is not the post owner
    IF post_owner_id IS NOT NULL AND post_owner_id != NEW.user_id THEN
        INSERT INTO notifications (user_id, actor_id, post_id, notification_type, message)
        VALUES (
            post_owner_id,
            NEW.user_id,
            NEW.post_id,
            'comment_on_post',
            actor_username || ' commented on your post: "' || post_title || '"'
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to create notification for comment reply
CREATE OR REPLACE FUNCTION notify_parent_comment_owner_on_reply()
RETURNS TRIGGER AS $$
DECLARE
    parent_owner_id INTEGER;
    actor_username TEXT;
    parent_comment_text TEXT;
BEGIN
    -- Only process if this is a reply (has parent_id)
    IF NEW.parent_id IS NOT NULL THEN
        -- Get the parent comment owner
        SELECT user_id, text INTO parent_owner_id, parent_comment_text
        FROM comments
        WHERE id = NEW.parent_id;

        -- Get the actor's username
        SELECT username INTO actor_username
        FROM users
        WHERE id = NEW.user_id;

        -- Only create notification if the reply author is not the parent comment owner
        IF parent_owner_id IS NOT NULL AND parent_owner_id != NEW.user_id THEN
            INSERT INTO notifications (user_id, actor_id, comment_id, notification_type, message)
            VALUES (
                parent_owner_id,
                NEW.user_id,
                NEW.parent_id,
                'comment_reply',
                actor_username || ' replied to your comment'
            );
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to create notification for comment edit
CREATE OR REPLACE FUNCTION notify_on_comment_edit()
RETURNS TRIGGER AS $$
DECLARE
    post_owner_id INTEGER;
    actor_username TEXT;
BEGIN
    -- Only process if the comment text actually changed
    IF OLD.text IS DISTINCT FROM NEW.text THEN
        -- Get the post owner if this is a top-level comment
        IF NEW.parent_id IS NULL THEN
            SELECT user_id INTO post_owner_id
            FROM posts
            WHERE id = NEW.post_id;

            -- Get the actor's username
            SELECT username INTO actor_username
            FROM users
            WHERE id = NEW.user_id;

            -- Only create notification if the editor is not the post owner
            IF post_owner_id IS NOT NULL AND post_owner_id != NEW.user_id THEN
                INSERT INTO notifications (user_id, actor_id, post_id, notification_type, message)
                VALUES (
                    post_owner_id,
                    NEW.user_id,
                    NEW.post_id,
                    'comment_edit',
                    actor_username || ' edited their comment on your post'
                );
            END IF;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER trigger_notify_post_owner_on_comment
    AFTER INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION notify_post_owner_on_comment();

CREATE TRIGGER trigger_notify_parent_comment_owner_on_reply
    AFTER INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION notify_parent_comment_owner_on_reply();

CREATE TRIGGER trigger_notify_on_comment_edit
    AFTER UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION notify_on_comment_edit();
