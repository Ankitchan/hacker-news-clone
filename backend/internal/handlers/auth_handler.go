package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/internal/repository"
	"github.com/Ankitchan/hackernews-clone/internal/utils"
	"github.com/Ankitchan/hackernews-clone/pkg/auth"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(db),
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var signupData models.UserSignup

	if err := utils.ParseJSONBody(r, &signupData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if signupData.Username == "" || signupData.Email == "" || signupData.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Trim and validate
	signupData.Username = strings.TrimSpace(signupData.Username)
	signupData.Email = strings.TrimSpace(signupData.Email)

	if len(signupData.Username) < 3 {
		utils.RespondWithError(w, http.StatusBadRequest, "Username must be at least 3 characters")
		return
	}

	// Validate password strength
	if err := auth.ValidatePasswordStrength(signupData.Password); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if email already exists
	emailExists, err := h.userRepo.EmailExists(signupData.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check email")
		return
	}
	if emailExists {
		utils.RespondWithError(w, http.StatusConflict, "Email already registered")
		return
	}

	// Check if username already exists
	usernameExists, err := h.userRepo.UsernameExists(signupData.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check username")
		return
	}
	if usernameExists {
		utils.RespondWithError(w, http.StatusConflict, "Username already taken")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(signupData.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := &models.User{
		Username:     signupData.Username,
		Email:        signupData.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.userRepo.Create(user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"token": token,
		"user": models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.UserLogin

	if err := utils.ParseJSONBody(r, &loginData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if (loginData.Email == "" && loginData.Username == "") || loginData.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username/email and password are required")
		return
	}

	loginData.Email = strings.TrimSpace(loginData.Email)
	loginData.Username = strings.TrimSpace(loginData.Username)

	// Get user by email or username
	var user *models.User
	var err error

	if loginData.Email != "" {
		user, err = h.userRepo.GetByEmail(loginData.Email)
	} else {
		user, err = h.userRepo.GetByUsername(loginData.Username)
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, loginData.Password); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"token": token,
		"user": models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// GetProfile returns the authenticated user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get full user details
	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// RefreshToken generates a new token for the authenticated user
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := utils.GetUserFromRequest(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Generate new token
	token, err := auth.GenerateToken(claims.UserID, claims.Username, claims.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
