package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all pending up migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	files, err := getMigrationFiles(migrationsPath, ".up.sql")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Execute pending migrations
	for _, file := range files {
		migrationName := strings.TrimSuffix(file, ".up.sql")

		// Skip if already applied
		if contains(appliedMigrations, migrationName) {
			log.Printf("Migration %s already applied, skipping\n", migrationName)
			continue
		}

		log.Printf("Running migration: %s\n", migrationName)

		// Read migration file
		filePath := filepath.Join(migrationsPath, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
		}

		// Record migration
		if err := recordMigration(db, migrationName); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
		}

		log.Printf("Migration %s applied successfully\n", migrationName)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the schema_migrations table
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			migration VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// getAppliedMigrations returns a list of already applied migrations
func getAppliedMigrations(db *sql.DB) ([]string, error) {
	query := "SELECT migration FROM schema_migrations ORDER BY migration"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		if err := rows.Scan(&migration); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, rows.Err()
}

// recordMigration records a migration as applied
func recordMigration(db *sql.DB, migration string) error {
	query := "INSERT INTO schema_migrations (migration) VALUES ($1)"
	_, err := db.Exec(query, migration)
	return err
}

// getMigrationFiles returns sorted list of migration files
func getMigrationFiles(migrationsPath, suffix string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), suffix) {
			files = append(files, d.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// contains checks if a string exists in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
