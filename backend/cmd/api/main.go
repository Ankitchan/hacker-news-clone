package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ankitchan/hackernews-clone/internal/middleware"
	"github.com/Ankitchan/hackernews-clone/internal/routes"
	"github.com/Ankitchan/hackernews-clone/pkg/auth"
	"github.com/Ankitchan/hackernews-clone/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get configuration from environment
	port := getEnv("PORT", "8080")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "hacker_news_db")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")
	jwtSecret := getEnv("JWT_SECRET", "")
	jwtExpirationHours := 72 // Default 3 days
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	// Validate required environment variables
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize JWT
	auth.InitJWT(jwtSecret, jwtExpirationHours)
	log.Println("JWT configuration initialized")

	// Connect to database
	dbConfig := database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}

	log.Printf("Connecting to database: %s on %s:%s", dbName, dbHost, dbPort)
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup router with all routes
	router := routes.SetupRoutes(db)

	// Setup CORS
	allowedOrigins := strings.Split(corsOrigins, ",")
	corsHandler := middleware.SetupCORS(allowedOrigins)

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)

	// Wrap router with CORS
	handler := corsHandler.Handler(router)

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Allowed CORS origins: %v", allowedOrigins)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
