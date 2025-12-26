package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg Config) (*sql.DB, error) {
	var dsn string
	if cfg.Password != "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.SSLMode,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.DBName,
			cfg.SSLMode,
		)
	}

	log.Printf("DSN (without password): host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Verify which database we're connected to
	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		return nil, fmt.Errorf("error verifying database: %w", err)
	}
	log.Printf("Successfully connected to database: %s", currentDB)

	DB = db
	return db, nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
