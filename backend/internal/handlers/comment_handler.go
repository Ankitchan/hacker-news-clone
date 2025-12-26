package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/internal/repository"
	"github.com/Ankitchan/hackernews-clone/internal/utils"
)

type CommentHandler struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

func NewCommentHandler(db *sql.DB) *CommentHandler {
	return &CommentHandler{
		commentRepo: repository.NewCommentRepository(db),
		postRepo:    repository.NewPostRepository(db),
	}
}

// Create creates a new comment
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var commentData models.CommentCreate
	if err := utils.ParseJSONBody(r, &commentData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if strings.TrimSpace(commentData.Text) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Comment text is required")
		return
	}

	commentData.Text = strings.TrimSpace(commentData.Text)

	if commentData.PostID <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Valid post ID is required")
		return
	}

	// Verify post exists
	_, err := h.postRepo.GetByID(commentData.PostID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// If parent_id is provided, verify it exists
	if commentData.ParentID != nil && *commentData.ParentID > 0 {
		_, err := h.commentRepo.GetByID(*commentData.ParentID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Parent comment not found")
			return
		}
	}

	// Create comment
	comment := &models.Comment{
		PostID: commentData.PostID,
		UserID: claims.UserID,
		Text:   commentData.Text,
	}

	if commentData.ParentID != nil && *commentData.ParentID > 0 {
		comment.ParentID = sql.NullInt64{Int64: int64(*commentData.ParentID), Valid: true}
	}

	if err := h.commentRepo.Create(comment); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Get the full comment with username
	fullComment, err := h.commentRepo.GetByID(comment.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve created comment")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, fullComment)
}

// CreateForPost creates a new comment for a specific post (post_id from URL)
func (h *CommentHandler) CreateForPost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get post ID from URL
	postID, err := utils.GetIDParam(r, "post_id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse request body
	var requestData struct {
		Text     string `json:"text"`
		ParentID *int   `json:"parent_id,omitempty"`
	}
	if err := utils.ParseJSONBody(r, &requestData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if strings.TrimSpace(requestData.Text) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Comment text is required")
		return
	}

	requestData.Text = strings.TrimSpace(requestData.Text)

	// Verify post exists
	_, err = h.postRepo.GetByID(postID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// If parent_id is provided, verify it exists
	if requestData.ParentID != nil && *requestData.ParentID > 0 {
		_, err := h.commentRepo.GetByID(*requestData.ParentID)
		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Parent comment not found")
			return
		}
	}

	// Create comment
	comment := &models.Comment{
		PostID: postID,
		UserID: claims.UserID,
		Text:   requestData.Text,
	}

	if requestData.ParentID != nil && *requestData.ParentID > 0 {
		comment.ParentID = sql.NullInt64{Int64: int64(*requestData.ParentID), Valid: true}
	}

	if err := h.commentRepo.Create(comment); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Get the full comment with username
	fullComment, err := h.commentRepo.GetByID(comment.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve created comment")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, fullComment)
}

// GetByID retrieves a single comment
func (h *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comment)
}

// GetByPostID retrieves all comments for a post
func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := utils.GetIDParam(r, "post_id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if threaded view is requested
	threaded := utils.GetQueryParam(r, "threaded", "true")

	var comments []models.Comment
	if threaded == "true" {
		comments, err = h.commentRepo.GetByPostIDThreaded(postID)
	} else {
		comments, err = h.commentRepo.GetByPostID(postID)
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve comments")
		return
	}

	response := models.CommentTree{
		Comments: comments,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetReplies retrieves replies to a comment
func (h *CommentHandler) GetReplies(w http.ResponseWriter, r *http.Request) {
	commentID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	replies, err := h.commentRepo.GetReplies(commentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve replies")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, replies)
}

// Update updates a comment
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing comment
	comment, err := h.commentRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	// Check ownership
	if comment.UserID != claims.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only edit your own comments")
		return
	}

	var updateData models.CommentUpdate
	if err := utils.ParseJSONBody(r, &updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate text
	if strings.TrimSpace(updateData.Text) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Comment text cannot be empty")
		return
	}

	comment.Text = strings.TrimSpace(updateData.Text)

	if err := h.commentRepo.Update(comment); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update comment")
		return
	}

	// Get updated comment
	updatedComment, err := h.commentRepo.GetByID(comment.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve updated comment")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedComment)
}

// Delete deletes a comment
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing comment
	comment, err := h.commentRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	// Check ownership
	if comment.UserID != claims.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only delete your own comments")
		return
	}

	if err := h.commentRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "Comment deleted successfully")
}
