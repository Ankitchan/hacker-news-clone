package repository

import (
	"database/sql"
	"fmt"

	"github.com/Ankitchan/hackernews-clone/internal/models"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// CreateOrUpdate creates a new vote or updates an existing one
func (r *VoteRepository) CreateOrUpdate(vote *models.Vote) error {
	// Check if vote already exists
	existing, err := r.GetUserVote(vote.UserID, vote.PostID, vote.CommentID)
	if err != nil && err.Error() != "vote not found" {
		return err
	}

	if existing != nil {
		// Update existing vote
		return r.Update(vote)
	}

	// Create new vote
	query := `
		INSERT INTO votes (user_id, post_id, comment_id, vote_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err = r.db.QueryRow(
		query,
		vote.UserID,
		vote.PostID,
		vote.CommentID,
		vote.VoteType,
	).Scan(&vote.ID, &vote.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

// Update updates an existing vote
func (r *VoteRepository) Update(vote *models.Vote) error {
	query := `
		UPDATE votes
		SET vote_type = $1
		WHERE user_id = $2 AND (
			(post_id = $3 AND comment_id IS NULL) OR
			(comment_id = $4 AND post_id IS NULL)
		)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		vote.VoteType,
		vote.UserID,
		vote.PostID,
		vote.CommentID,
	).Scan(&vote.ID, &vote.CreatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("vote not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update vote: %w", err)
	}

	return nil
}

// Delete deletes a vote
func (r *VoteRepository) Delete(id int) error {
	query := `DELETE FROM votes WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

// DeleteUserVote deletes a user's vote on a post or comment
func (r *VoteRepository) DeleteUserVote(userID int, postID, commentID *int) error {
	query := `
		DELETE FROM votes
		WHERE user_id = $1 AND (
			(post_id = $2 AND comment_id IS NULL) OR
			(comment_id = $3 AND post_id IS NULL)
		)
	`

	result, err := r.db.Exec(query, userID, postID, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vote not found")
	}

	return nil
}

// GetUserVote retrieves a user's vote on a post or comment
func (r *VoteRepository) GetUserVote(userID int, postID, commentID *int) (*models.Vote, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, vote_type, created_at
		FROM votes
		WHERE user_id = $1 AND (
			(post_id = $2 AND comment_id IS NULL) OR
			(comment_id = $3 AND post_id IS NULL)
		)
	`

	vote := &models.Vote{}
	err := r.db.QueryRow(query, userID, postID, commentID).Scan(
		&vote.ID,
		&vote.UserID,
		&vote.PostID,
		&vote.CommentID,
		&vote.VoteType,
		&vote.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vote not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	return vote, nil
}

// GetPostVotes retrieves all votes for a post
func (r *VoteRepository) GetPostVotes(postID int) ([]models.Vote, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, vote_type, created_at
		FROM votes
		WHERE post_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post votes: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(
			&vote.ID,
			&vote.UserID,
			&vote.PostID,
			&vote.CommentID,
			&vote.VoteType,
			&vote.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	return votes, rows.Err()
}

// GetCommentVotes retrieves all votes for a comment
func (r *VoteRepository) GetCommentVotes(commentID int) ([]models.Vote, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, vote_type, created_at
		FROM votes
		WHERE comment_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment votes: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(
			&vote.ID,
			&vote.UserID,
			&vote.PostID,
			&vote.CommentID,
			&vote.VoteType,
			&vote.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	return votes, rows.Err()
}

// CalculatePostPoints calculates the total points for a post
func (r *VoteRepository) CalculatePostPoints(postID int) (int, error) {
	query := `
		SELECT COALESCE(SUM(vote_type), 0)
		FROM votes
		WHERE post_id = $1
	`

	var points int
	err := r.db.QueryRow(query, postID).Scan(&points)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate post points: %w", err)
	}

	return points, nil
}

// CalculateCommentPoints calculates the total points for a comment
func (r *VoteRepository) CalculateCommentPoints(commentID int) (int, error) {
	query := `
		SELECT COALESCE(SUM(vote_type), 0)
		FROM votes
		WHERE comment_id = $1
	`

	var points int
	err := r.db.QueryRow(query, commentID).Scan(&points)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate comment points: %w", err)
	}

	return points, nil
}

// GetUserVotesOnPosts retrieves all user votes on posts
func (r *VoteRepository) GetUserVotesOnPosts(userID int) ([]models.Vote, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, vote_type, created_at
		FROM votes
		WHERE user_id = $1 AND post_id IS NOT NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user post votes: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(
			&vote.ID,
			&vote.UserID,
			&vote.PostID,
			&vote.CommentID,
			&vote.VoteType,
			&vote.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	return votes, rows.Err()
}

// GetUserVotesOnComments retrieves all user votes on comments
func (r *VoteRepository) GetUserVotesOnComments(userID int) ([]models.Vote, error) {
	query := `
		SELECT id, user_id, post_id, comment_id, vote_type, created_at
		FROM votes
		WHERE user_id = $1 AND comment_id IS NOT NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comment votes: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(
			&vote.ID,
			&vote.UserID,
			&vote.PostID,
			&vote.CommentID,
			&vote.VoteType,
			&vote.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}

	return votes, rows.Err()
}
