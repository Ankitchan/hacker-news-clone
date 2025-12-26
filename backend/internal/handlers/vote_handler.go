package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/internal/repository"
	"github.com/Ankitchan/hackernews-clone/internal/utils"
)

type VoteHandler struct {
	voteRepo    *repository.VoteRepository
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
}

func NewVoteHandler(db *sql.DB) *VoteHandler {
	return &VoteHandler{
		voteRepo:    repository.NewVoteRepository(db),
		postRepo:    repository.NewPostRepository(db),
		commentRepo: repository.NewCommentRepository(db),
	}
}

// VoteOnPost handles voting on a post
func (h *VoteHandler) VoteOnPost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var voteData models.VoteCreate
	if err := utils.ParseJSONBody(r, &voteData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate vote type
	if voteData.VoteType != models.VoteUp && voteData.VoteType != models.VoteDown {
		utils.RespondWithError(w, http.StatusBadRequest, "Vote type must be 1 (upvote) or -1 (downvote)")
		return
	}

	// Verify post exists
	_, err = h.postRepo.GetByID(postID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Create or update vote
	vote := &models.Vote{
		UserID:   claims.UserID,
		PostID:   &postID,
		VoteType: voteData.VoteType,
	}

	if err := h.voteRepo.CreateOrUpdate(vote); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process vote")
		return
	}

	// Calculate new points
	points, err := h.voteRepo.CalculatePostPoints(postID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to calculate points")
		return
	}

	// Update post points
	if err := h.postRepo.UpdatePoints(postID, points); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update post points")
		return
	}

	response := models.VoteResponse{
		Success: true,
		Points:  points,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// VoteOnComment handles voting on a comment
func (h *VoteHandler) VoteOnComment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var voteData models.VoteCreate
	if err := utils.ParseJSONBody(r, &voteData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate vote type
	if voteData.VoteType != models.VoteUp && voteData.VoteType != models.VoteDown {
		utils.RespondWithError(w, http.StatusBadRequest, "Vote type must be 1 (upvote) or -1 (downvote)")
		return
	}

	// Verify comment exists
	_, err = h.commentRepo.GetByID(commentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	// Create or update vote
	vote := &models.Vote{
		UserID:    claims.UserID,
		CommentID: &commentID,
		VoteType:  voteData.VoteType,
	}

	if err := h.voteRepo.CreateOrUpdate(vote); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process vote")
		return
	}

	// Calculate new points
	points, err := h.voteRepo.CalculateCommentPoints(commentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to calculate points")
		return
	}

	// Update comment points
	if err := h.commentRepo.UpdatePoints(commentID, points); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update comment points")
		return
	}

	response := models.VoteResponse{
		Success: true,
		Points:  points,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// UnvotePost removes a vote from a post
func (h *VoteHandler) UnvotePost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Verify post exists
	_, err = h.postRepo.GetByID(postID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Delete vote
	if err := h.voteRepo.DeleteUserVote(claims.UserID, &postID, nil); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Vote not found")
		return
	}

	// Calculate new points
	points, err := h.voteRepo.CalculatePostPoints(postID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to calculate points")
		return
	}

	// Update post points
	if err := h.postRepo.UpdatePoints(postID, points); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update post points")
		return
	}

	response := models.VoteResponse{
		Success: true,
		Points:  points,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// UnvoteComment removes a vote from a comment
func (h *VoteHandler) UnvoteComment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Verify comment exists
	_, err = h.commentRepo.GetByID(commentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	// Delete vote
	if err := h.voteRepo.DeleteUserVote(claims.UserID, nil, &commentID); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Vote not found")
		return
	}

	// Calculate new points
	points, err := h.voteRepo.CalculateCommentPoints(commentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to calculate points")
		return
	}

	// Update comment points
	if err := h.commentRepo.UpdatePoints(commentID, points); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update comment points")
		return
	}

	response := models.VoteResponse{
		Success: true,
		Points:  points,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetUserVote retrieves a user's vote on a post or comment
func (h *VoteHandler) GetUserVoteOnPost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	vote, err := h.voteRepo.GetUserVote(claims.UserID, &postID, nil)
	if err != nil {
		// No vote found
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"vote": nil})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, vote)
}

// GetUserVoteOnComment retrieves a user's vote on a comment
func (h *VoteHandler) GetUserVoteOnComment(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	commentID, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	vote, err := h.voteRepo.GetUserVote(claims.UserID, nil, &commentID)
	if err != nil {
		// No vote found
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"vote": nil})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, vote)
}
