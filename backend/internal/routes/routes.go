package routes

import (
	"database/sql"
	"net/http"

	"github.com/Ankitchan/hackernews-clone/internal/handlers"
	"github.com/Ankitchan/hackernews-clone/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	postHandler := handlers.NewPostHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	voteHandler := handlers.NewVoteHandler(db)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Health check
	api.HandleFunc("/health", healthCheck).Methods("GET")

	// ===========================
	// Authentication Routes (Public)
	// ===========================
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Protected auth routes
	authProtected := authRoutes.PathPrefix("").Subrouter()
	authProtected.Use(middleware.AuthMiddleware)
	authProtected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	authProtected.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// ===========================
	// Post Routes
	// ===========================
	postRoutes := api.PathPrefix("/posts").Subrouter()

	// Public post routes
	postRoutes.HandleFunc("", postHandler.GetAll).Methods("GET")
	postRoutes.HandleFunc("/new", postHandler.GetByNew).Methods("GET")
	postRoutes.HandleFunc("/top", postHandler.GetByTop).Methods("GET")
	postRoutes.HandleFunc("/best", postHandler.GetByBest).Methods("GET")
	postRoutes.HandleFunc("/search", postHandler.Search).Methods("GET")
	postRoutes.HandleFunc("/{id:[0-9]+}", postHandler.GetByID).Methods("GET")

	// Protected post routes (require authentication)
	postProtected := postRoutes.PathPrefix("").Subrouter()
	postProtected.Use(middleware.AuthMiddleware)
	postProtected.HandleFunc("", postHandler.Create).Methods("POST")
	postProtected.HandleFunc("/{id:[0-9]+}", postHandler.Update).Methods("PUT", "PATCH")
	postProtected.HandleFunc("/{id:[0-9]+}", postHandler.Delete).Methods("DELETE")

	// ===========================
	// Comment Routes
	// ===========================
	commentRoutes := api.PathPrefix("/comments").Subrouter()

	// Public comment routes
	commentRoutes.HandleFunc("/{id:[0-9]+}", commentHandler.GetByID).Methods("GET")
	commentRoutes.HandleFunc("/{id:[0-9]+}/replies", commentHandler.GetReplies).Methods("GET")

	// Get comments by post ID (public)
	api.HandleFunc("/posts/{post_id:[0-9]+}/comments", commentHandler.GetByPostID).Methods("GET")

	// Create comment for a specific post (protected)
	postCommentProtected := api.PathPrefix("/posts/{post_id:[0-9]+}/comments").Subrouter()
	postCommentProtected.Use(middleware.AuthMiddleware)
	postCommentProtected.HandleFunc("", commentHandler.CreateForPost).Methods("POST")

	// Protected comment routes
	commentProtected := commentRoutes.PathPrefix("").Subrouter()
	commentProtected.Use(middleware.AuthMiddleware)
	commentProtected.HandleFunc("", commentHandler.Create).Methods("POST")
	commentProtected.HandleFunc("/{id:[0-9]+}", commentHandler.Update).Methods("PUT", "PATCH")
	commentProtected.HandleFunc("/{id:[0-9]+}", commentHandler.Delete).Methods("DELETE")

	// ===========================
	// Vote Routes (All Protected)
	// ===========================
	voteRoutes := api.PathPrefix("/votes").Subrouter()
	voteRoutes.Use(middleware.AuthMiddleware)

	// Post votes
	voteRoutes.HandleFunc("/posts/{id:[0-9]+}", voteHandler.VoteOnPost).Methods("POST")
	voteRoutes.HandleFunc("/posts/{id:[0-9]+}", voteHandler.UnvotePost).Methods("DELETE")
	voteRoutes.HandleFunc("/posts/{id:[0-9]+}/user", voteHandler.GetUserVoteOnPost).Methods("GET")

	// Comment votes
	voteRoutes.HandleFunc("/comments/{id:[0-9]+}", voteHandler.VoteOnComment).Methods("POST")
	voteRoutes.HandleFunc("/comments/{id:[0-9]+}", voteHandler.UnvoteComment).Methods("DELETE")
	voteRoutes.HandleFunc("/comments/{id:[0-9]+}/user", voteHandler.GetUserVoteOnComment).Methods("GET")

	return router
}

// healthCheck returns the health status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","message":"Hacker News Clone API is running"}`))
}
