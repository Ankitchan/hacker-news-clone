package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/internal/repository"
	"github.com/Ankitchan/hackernews-clone/internal/utils"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: repository.NewNotificationRepository(db),
	}
}

// GetAll retrieves all notifications for the authenticated user
func (h *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
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

	notifications, totalCount, unreadCount, err := h.notificationRepo.GetByUserID(claims.UserID, pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	response := models.NotificationList{
		Notifications: notifications,
		TotalCount:    totalCount,
		UnreadCount:   unreadCount,
		Page:          page,
		PageSize:      pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetUnread retrieves only unread notifications for the authenticated user
func (h *NotificationHandler) GetUnread(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
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

	notifications, unreadCount, err := h.notificationRepo.GetUnreadByUserID(claims.UserID, pageSize, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve unread notifications")
		return
	}

	response := models.NotificationList{
		Notifications: notifications,
		TotalCount:    unreadCount,
		UnreadCount:   unreadCount,
		Page:          page,
		PageSize:      pageSize,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// MarkAsRead marks specific notifications as read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var reqData models.MarkAsReadRequest
	if err := utils.ParseJSONBody(r, &reqData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(reqData.NotificationIDs) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No notification IDs provided")
		return
	}

	err := h.notificationRepo.MarkAsRead(claims.UserID, reqData.NotificationIDs)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark notifications as read")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "Notifications marked as read")
}

// MarkAllAsRead marks all notifications as read for the authenticated user
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.notificationRepo.MarkAllAsRead(claims.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "All notifications marked as read")
}

// Delete deletes a specific notification
func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = h.notificationRepo.Delete(claims.UserID, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	utils.RespondWithSuccess(w, http.StatusOK, nil, "Notification deleted successfully")
}

// GetUnreadCount returns the count of unread notifications
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	count, err := h.notificationRepo.GetUnreadCount(claims.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]int{"unread_count": count})
}
