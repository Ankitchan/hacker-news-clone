-- Drop triggers
DROP TRIGGER IF EXISTS set_comment_depth ON comments;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_comment_depth();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_votes_comment_id;
DROP INDEX IF EXISTS idx_votes_post_id;
DROP INDEX IF EXISTS idx_votes_user_id;

DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;

DROP INDEX IF EXISTS idx_posts_points;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_user_id;

-- Drop tables (in reverse order due to foreign key constraints)
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
