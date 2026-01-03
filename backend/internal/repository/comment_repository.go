package repository

import (
	"database/sql"
	"fmt"

	"github.com/Ankitchan/hackernews-clone/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create creates a new comment
func (r *CommentRepository) Create(comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, parent_id, text)
		VALUES ($1, $2, $3, $4)
		RETURNING id, points, depth, created_at, updated_at
	`

	var parentID interface{}
	if comment.ParentID.Valid {
		parentID = comment.ParentID.Int64
	} else {
		parentID = nil
	}

	err := r.db.QueryRow(
		query,
		comment.PostID,
		comment.UserID,
		parentID,
		comment.Text,
	).Scan(&comment.ID, &comment.Points, &comment.Depth, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetByID retrieves a comment by ID
func (r *CommentRepository) GetByID(id int) (*models.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_id, c.text, c.points, c.depth,
			c.created_at, c.updated_at,
			u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	comment := &models.Comment{}
	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.ParentID,
		&comment.Text,
		&comment.Points,
		&comment.Depth,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Username,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

// GetByPostID retrieves all comments for a post (flat list)
func (r *CommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_id, c.text, c.points, c.depth,
			c.created_at, c.updated_at,
			u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.ParentID,
			&comment.Text,
			&comment.Points,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// GetByPostIDThreaded retrieves comments for a post in threaded structure
func (r *CommentRepository) GetByPostIDThreaded(postID int) ([]models.Comment, error) {
	// Get all comments for the post
	comments, err := r.GetByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Build threaded structure
	return buildCommentTree(comments), nil
}

// buildCommentTree converts a flat list of comments into a threaded structure
func buildCommentTree(comments []models.Comment) []models.Comment {
	// Create a map for quick lookup and initialize replies
	commentMap := make(map[int]*models.Comment)
	for i := range comments {
		comments[i].Replies = []models.Comment{}
		comments[i].TotalPoints = comments[i].Points // Initialize with own points
		commentMap[comments[i].ID] = &comments[i]
	}

	// Track root comment IDs
	var rootIDs []int

	// Attach children to their parents
	for i := range comments {
		if comments[i].ParentID.Valid {
			// This is a child comment
			parentID := int(comments[i].ParentID.Int64)
			if parent, ok := commentMap[parentID]; ok {
				// Append pointer reference so nested updates propagate
				parent.Replies = append(parent.Replies, *commentMap[comments[i].ID])
			}
		} else {
			// Track root comment IDs
			rootIDs = append(rootIDs, comments[i].ID)
		}
	}

	// Collect root comments - need to build recursively to capture all nested levels
	var rootComments []models.Comment
	for _, id := range rootIDs {
		if comment, ok := commentMap[id]; ok {
			rootComments = append(rootComments, buildCommentWithReplies(comment, commentMap))
		}
	}

	return rootComments
}

// buildCommentWithReplies recursively builds a comment with all its nested replies
func buildCommentWithReplies(comment *models.Comment, commentMap map[int]*models.Comment) models.Comment {
	result := *comment
	result.Replies = []models.Comment{}

	// Initialize with comment's own points
	result.TotalPoints = result.Points

	// Find all direct children
	for id, c := range commentMap {
		if c.ParentID.Valid && int(c.ParentID.Int64) == comment.ID {
			// Recursively build this child with its replies
			childComment := buildCommentWithReplies(commentMap[id], commentMap)
			result.Replies = append(result.Replies, childComment)

			// Add child's total points to parent's total
			result.TotalPoints += childComment.TotalPoints
		}
	}

	return result
}

// GetByUserID retrieves all comments by a specific user
func (r *CommentRepository) GetByUserID(userID int, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_id, c.text, c.points, c.depth,
			c.created_at, c.updated_at,
			u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.ParentID,
			&comment.Text,
			&comment.Points,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// GetReplies retrieves all direct replies to a comment
func (r *CommentRepository) GetReplies(commentID int) ([]models.Comment, error) {
	query := `
		SELECT
			c.id, c.post_id, c.user_id, c.parent_id, c.text, c.points, c.depth,
			c.created_at, c.updated_at,
			u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.parent_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(query, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.ParentID,
			&comment.Text,
			&comment.Points,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

// Update updates a comment's text
func (r *CommentRepository) Update(comment *models.Comment) error {
	query := `
		UPDATE comments
		SET text = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, comment.Text, comment.ID).Scan(&comment.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("comment not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// UpdatePoints updates the points for a comment
func (r *CommentRepository) UpdatePoints(commentID, points int) error {
	query := `UPDATE comments SET points = $1 WHERE id = $2`

	result, err := r.db.Exec(query, points, commentID)
	if err != nil {
		return fmt.Errorf("failed to update comment points: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// Delete deletes a comment
func (r *CommentRepository) Delete(id int) error {
	query := `DELETE FROM comments WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found")
	}

	return nil
}

// GetCountByPostID gets the total number of comments for a post
func (r *CommentRepository) GetCountByPostID(postID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int
	err := r.db.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment count: %w", err)
	}

	return count, nil
}
