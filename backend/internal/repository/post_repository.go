package repository

import (
	"database/sql"
	"fmt"

	"github.com/Ankitchan/hackernews-clone/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(post *models.Post) error {
	query := `
		INSERT INTO posts (title, url, text, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, points, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		post.Title,
		post.URL,
		post.Text,
		post.UserID,
	).Scan(&post.ID, &post.Points, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID retrieves a post by ID with user information
func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	query := `
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = $1
		GROUP BY p.id, u.username
	`

	post := &models.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.URL,
		&post.Text,
		&post.UserID,
		&post.Points,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Username,
		&post.CommentCount,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

// GetAll retrieves all posts with pagination
func (r *PostRepository) GetAll(limit, offset int) ([]models.Post, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM posts`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get post count: %w", err)
	}

	// Get posts
	query := `
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.URL,
			&post.Text,
			&post.UserID,
			&post.Points,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, totalCount, nil
}

// GetByNew retrieves posts sorted by newest first
func (r *PostRepository) GetByNew(limit, offset int) ([]models.Post, int, error) {
	return r.getPostsWithOrder(limit, offset, "p.created_at DESC")
}

// GetByTop retrieves posts sorted by Hacker News ranking formula
// Formula: Score = (Points - 1)^0.8 / (Age + 2)^1.8
// Where Age is in hours, Points ignores the submitter's upvote
func (r *PostRepository) GetByTop(limit, offset int) ([]models.Post, int, error) {
	// HN ranking formula: (P-1)^0.8 / (T+2)^1.8
	// P = points (minus submitter's vote), T = age in hours
	orderBy := `(
		POWER(GREATEST(p.points - 1, 0), 0.8) /
		POWER(GREATEST(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - p.created_at))/3600, 0) + 2, 1.8)
	) DESC`
	return r.getPostsWithOrder(limit, offset, orderBy)
}

// GetByBest retrieves posts sorted by Wilson Score Interval
// This algorithm calculates a confidence interval for the true quality score
// A post with 10 upvotes and 0 downvotes ranks higher than 100 up and 30 down
func (r *PostRepository) GetByBest(limit, offset int) ([]models.Post, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM posts`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get post count: %w", err)
	}

	// Wilson Score Interval formula (95% confidence)
	// Lower bound of Wilson score confidence interval for a Bernoulli parameter
	// where: phat = upvotes/total, z = 1.96 (95% confidence), n = total votes
	query := `
		WITH vote_counts AS (
			SELECT
				post_id,
				COUNT(*) FILTER (WHERE vote_type = 1) AS upvotes,
				COUNT(*) FILTER (WHERE vote_type = -1) AS downvotes,
				COUNT(*) AS total_votes
			FROM votes
			WHERE post_id IS NOT NULL
			GROUP BY post_id
		)
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN vote_counts vc ON p.id = vc.post_id
		GROUP BY p.id, u.username, vc.upvotes, vc.downvotes, vc.total_votes
		ORDER BY
			CASE
				WHEN COALESCE(vc.total_votes, 0) = 0 THEN 0
				ELSE (
					(COALESCE(vc.upvotes, 0)::float / COALESCE(vc.total_votes, 1)::float + 1.96 * 1.96 / (2 * COALESCE(vc.total_votes, 1)) -
					1.96 * SQRT((COALESCE(vc.upvotes, 0)::float / COALESCE(vc.total_votes, 1)::float *
					(1 - COALESCE(vc.upvotes, 0)::float / COALESCE(vc.total_votes, 1)::float) +
					1.96 * 1.96 / (4 * COALESCE(vc.total_votes, 1))) / COALESCE(vc.total_votes, 1))) /
					(1 + 1.96 * 1.96 / COALESCE(vc.total_votes, 1))
				)
			END DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.URL,
			&post.Text,
			&post.UserID,
			&post.Points,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, totalCount, nil
}

// getPostsWithOrder is a helper function to get posts with custom ordering
func (r *PostRepository) getPostsWithOrder(limit, offset int, orderBy string) ([]models.Post, int, error) {
	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM posts`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get post count: %w", err)
	}

	// Get posts with custom ordering
	query := fmt.Sprintf(`
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		GROUP BY p.id, u.username
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.URL,
			&post.Text,
			&post.UserID,
			&post.Points,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, totalCount, nil
}

// GetByUserID retrieves all posts by a specific user
func (r *PostRepository) GetByUserID(userID, limit, offset int) ([]models.Post, error) {
	query := `
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.URL,
			&post.Text,
			&post.UserID,
			&post.Points,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// Search searches posts by title or text
func (r *PostRepository) Search(searchTerm string, limit, offset int) ([]models.Post, int, error) {
	searchPattern := "%" + searchTerm + "%"

	// Get total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM posts
		WHERE title ILIKE $1 OR text ILIKE $1
	`
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search count: %w", err)
	}

	// Search posts
	query := `
		SELECT
			p.id, p.title, p.url, p.text, p.user_id, p.points, p.created_at, p.updated_at,
			u.username,
			COALESCE(COUNT(DISTINCT c.id), 0) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.title ILIKE $1 OR p.text ILIKE $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.URL,
			&post.Text,
			&post.UserID,
			&post.Points,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, totalCount, rows.Err()
}

// Update updates a post
func (r *PostRepository) Update(post *models.Post) error {
	query := `
		UPDATE posts
		SET title = $1, url = $2, text = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		post.Title,
		post.URL,
		post.Text,
		post.ID,
	).Scan(&post.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("post not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// UpdatePoints updates the points for a post
func (r *PostRepository) UpdatePoints(postID, points int) error {
	query := `UPDATE posts SET points = $1 WHERE id = $2`

	result, err := r.db.Exec(query, points, postID)
	if err != nil {
		return fmt.Errorf("failed to update post points: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// Delete deletes a post
func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}
