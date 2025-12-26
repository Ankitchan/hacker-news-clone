package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/internal/repository"
	"github.com/Ankitchan/hackernews-clone/internal/utils"
)

type PostHandler struct {
	postRepo *repository.PostRepository
}

func NewPostHandler(db *sql.DB) *PostHandler {
	return &PostHandler{
		postRepo: repository.NewPostRepository(db),
	}
}

// Create creates a new post
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var postData models.PostCreate
	if err := utils.ParseJSONBody(r, &postData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if strings.TrimSpace(postData.Title) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	postData.Title = strings.TrimSpace(postData.Title)

	// Validate that either URL or Text is provided (but not necessarily both)
	hasURL := postData.URL != nil && strings.TrimSpace(*postData.URL) != ""
	hasText := postData.Text != nil && strings.TrimSpace(*postData.Text) != ""

	if !hasURL && !hasText {
		utils.RespondWithError(w, http.StatusBadRequest, "Either URL or text content is required")
		return
	}

	// Create post
	post := &models.Post{
		Title:  postData.Title,
		UserID: claims.UserID,
	}

	if hasURL {
		trimmedURL := strings.TrimSpace(*postData.URL)
		post.URL = sql.NullString{String: trimmedURL, Valid: true}
	}

	if hasText {
		trimmedText := strings.TrimSpace(*postData.Text)
		post.Text = sql.NullString{String: trimmedText, Valid: true}
	}

	if err := h.postRepo.Create(post); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Get the full post with username
	fullPost, err := h.postRepo.GetByID(post.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve created post")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, fullPost)
}

// GetByID retrieves a single post
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDParam(r, "id")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, post)
}

// GetAll retrieves all posts with pagination
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, err := utils.GetQueryParamInt(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := utils.GetQueryParamInt(r, "page_size", 20)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, totalCount, err := h.postRepo.GetAll(pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	response := models.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetByNew retrieves posts sorted by newest
func (h *PostHandler) GetByNew(w http.ResponseWriter, r *http.Request) {
	page, err := utils.GetQueryParamInt(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := utils.GetQueryParamInt(r, "page_size", 20)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, totalCount, err := h.postRepo.GetByNew(pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	response := models.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetByTop retrieves posts sorted by points
func (h *PostHandler) GetByTop(w http.ResponseWriter, r *http.Request) {
	page, err := utils.GetQueryParamInt(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := utils.GetQueryParamInt(r, "page_size", 20)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, totalCount, err := h.postRepo.GetByTop(pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	response := models.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetByBest retrieves posts sorted by best algorithm
func (h *PostHandler) GetByBest(w http.ResponseWriter, r *http.Request) {
	page, err := utils.GetQueryParamInt(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := utils.GetQueryParamInt(r, "page_size", 20)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, totalCount, err := h.postRepo.GetByBest(pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	response := models.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// Search searches posts
func (h *PostHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := utils.GetQueryParam(r, "q", "")
	if strings.TrimSpace(query) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	page, err := utils.GetQueryParamInt(r, "page", 1)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := utils.GetQueryParamInt(r, "page_size", 20)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	posts, totalCount, err := h.postRepo.Search(query, pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to search posts")
		return
	}

	response := models.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// Update updates a post
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	// Get existing post
	post, err := h.postRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Check ownership
	if post.UserID != claims.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only edit your own posts")
		return
	}

	var updateData models.PostUpdate
	if err := utils.ParseJSONBody(r, &updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields if provided
	if updateData.Title != nil {
		post.Title = strings.TrimSpace(*updateData.Title)
		if post.Title == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Title cannot be empty")
			return
		}
	}

	if updateData.URL != nil {
		trimmedURL := strings.TrimSpace(*updateData.URL)
		post.URL = sql.NullString{String: trimmedURL, Valid: trimmedURL != ""}
	}

	if updateData.Text != nil {
		trimmedText := strings.TrimSpace(*updateData.Text)
		post.Text = sql.NullString{String: trimmedText, Valid: trimmedText != ""}
	}

	if err := h.postRepo.Update(post); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// Get updated post with full details
	updatedPost, err := h.postRepo.GetByID(post.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve updated post")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedPost)
}

// Delete deletes a post
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	// Get existing post
	post, err := h.postRepo.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Check ownership
	if post.UserID != claims.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "You can only delete your own posts")
		return
	}

	if err := h.postRepo.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "Post deleted successfully")
}
