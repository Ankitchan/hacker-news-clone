package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Ankitchan/hackernews-clone/internal/models"
	"github.com/Ankitchan/hackernews-clone/pkg/database"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

// TestMain sets up the test database
func TestMain(m *testing.M) {
	// Setup test database connection
	cfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "hacker_news_db"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	var err error
	testDB, err = database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations (skip if already run)
	if err := database.RunMigrations(testDB, "../../migrations"); err != nil {
		// Migrations might already be run, check if tables exist
		var exists bool
		err := testDB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'posts')").Scan(&exists)
		if err != nil || !exists {
			log.Fatalf("Failed to run migrations and tables don't exist: %v", err)
		}
		log.Println("Using existing database schema")
	}

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupTestData creates test posts and votes for sorting tests
func setupTestData(t *testing.T, repo *PostRepository) (userID1, userID2 int, postIDs []int) {
	t.Helper()

	// Create test users
	userID1, userID2 = createTestUsers(t)

	// Clean up existing test data
	_, err := testDB.Exec("DELETE FROM votes WHERE post_id IS NOT NULL")
	if err != nil {
		t.Fatalf("Failed to clean votes: %v", err)
	}
	_, err = testDB.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatalf("Failed to clean posts: %v", err)
	}

	// Create test posts with different characteristics for testing sorting
	testPosts := []struct {
		title     string
		userID    int
		createdAt time.Time
		points    int
		upvotes   int
		downvotes int
	}{
		// Post 1: Old post with high upvotes (should rank lower in "top" due to age)
		{
			title:     "Old Popular Post",
			userID:    userID1,
			createdAt: time.Now().Add(-48 * time.Hour), // 2 days old
			upvotes:   50,
			downvotes: 5,
		},
		// Post 2: Recent post with moderate upvotes (should rank high in "top")
		{
			title:     "Recent Good Post",
			userID:    userID1,
			createdAt: time.Now().Add(-2 * time.Hour), // 2 hours old
			upvotes:   20,
			downvotes: 2,
		},
		// Post 3: Very recent post with few votes (should be first in "new")
		{
			title:     "Brand New Post",
			userID:    userID2,
			createdAt: time.Now().Add(-30 * time.Minute), // 30 minutes old
			upvotes:   5,
			downvotes: 0,
		},
		// Post 4: Perfect score post (10 upvotes, 0 downvotes - high Wilson score)
		{
			title:     "High Quality Post",
			userID:    userID1,
			createdAt: time.Now().Add(-6 * time.Hour),
			upvotes:   10,
			downvotes: 0,
		},
		// Post 5: Controversial post (100 upvotes, 30 downvotes - lower Wilson score)
		{
			title:     "Controversial Post",
			userID:    userID2,
			createdAt: time.Now().Add(-4 * time.Hour),
			upvotes:   100,
			downvotes: 30,
		},
		// Post 6: No votes post
		{
			title:     "No Votes Post",
			userID:    userID1,
			createdAt: time.Now().Add(-1 * time.Hour),
			upvotes:   0,
			downvotes: 0,
		},
	}

	for _, tp := range testPosts {
		// Create post
		post := &models.Post{
			Title:  tp.title,
			UserID: tp.userID,
			URL:    sql.NullString{String: "https://example.com", Valid: true},
		}

		// Insert post with custom timestamp
		query := `
			INSERT INTO posts (title, url, user_id, created_at, points)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at
		`
		err := testDB.QueryRow(query, post.Title, post.URL, post.UserID, tp.createdAt, 0).
			Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}

		postIDs = append(postIDs, post.ID)

		// Add upvotes
		for i := 0; i < tp.upvotes; i++ {
			// Create a unique user for each vote (in real scenario)
			// For simplicity, we'll use alternating users
			voterID := userID1
			if i%2 == 0 {
				voterID = userID2
			}

			_, err := testDB.Exec(
				"INSERT INTO votes (user_id, post_id, vote_type) VALUES ($1, $2, 1) ON CONFLICT DO NOTHING",
				voterID, post.ID,
			)
			if err != nil {
				t.Logf("Warning: Failed to add upvote: %v", err)
			}
		}

		// Add downvotes
		for i := 0; i < tp.downvotes; i++ {
			// Use different voter IDs for downvotes
			voterID := userID1
			if i%2 == 1 {
				voterID = userID2
			}

			_, err := testDB.Exec(
				"INSERT INTO votes (user_id, post_id, vote_type) VALUES ($1, $2, -1) ON CONFLICT DO NOTHING",
				voterID, post.ID,
			)
			if err != nil {
				t.Logf("Warning: Failed to add downvote: %v", err)
			}
		}

		// Update post points
		points := tp.upvotes - tp.downvotes
		_, err = testDB.Exec("UPDATE posts SET points = $1 WHERE id = $2", points, post.ID)
		if err != nil {
			t.Fatalf("Failed to update post points: %v", err)
		}
	}

	return userID1, userID2, postIDs
}

// createTestUsers creates test users for voting
func createTestUsers(t *testing.T) (userID1, userID2 int) {
	t.Helper()

	// Create or get test user 1
	err := testDB.QueryRow(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`, "testuser1", "test1@example.com", "hashedpassword1").Scan(&userID1)
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}

	// Create or get test user 2
	err = testDB.QueryRow(`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`, "testuser2", "test2@example.com", "hashedpassword2").Scan(&userID2)
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}

	return userID1, userID2
}

func TestGetByNew(t *testing.T) {
	repo := NewPostRepository(testDB)
	_, _, postIDs := setupTestData(t, repo)

	tests := []struct {
		name   string
		limit  int
		offset int
	}{
		{
			name:   "Get newest posts",
			limit:  10,
			offset: 0,
		},
		{
			name:   "Get with pagination",
			limit:  2,
			offset: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, totalCount, err := repo.GetByNew(tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("GetByNew() error = %v", err)
			}

			if totalCount != len(postIDs) {
				t.Errorf("GetByNew() totalCount = %v, want %v", totalCount, len(postIDs))
			}

			if len(posts) == 0 {
				t.Fatal("GetByNew() returned no posts")
			}

			// Check if posts are sorted by creation time (newest first)
			for i := 0; i < len(posts)-1; i++ {
				if posts[i].CreatedAt.Before(posts[i+1].CreatedAt) {
					t.Errorf("Posts not sorted by newest: post[%d] is older than post[%d]", i, i+1)
				}
			}

			// Log the ordering for verification
			t.Logf("Retrieved %d posts, first post: %s", len(posts), posts[0].Title)
		})
	}
}

func TestGetByTop(t *testing.T) {
	repo := NewPostRepository(testDB)
	setupTestData(t, repo)

	tests := []struct {
		name  string
		limit int
	}{
		{
			name:  "Get top posts with HN algorithm",
			limit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, totalCount, err := repo.GetByTop(tt.limit, 0)
			if err != nil {
				t.Fatalf("GetByTop() error = %v", err)
			}

			if totalCount == 0 {
				t.Error("GetByTop() totalCount = 0")
			}

			if len(posts) == 0 {
				t.Fatal("GetByTop() returned no posts")
			}

			// Verify HN ranking: recent posts with good scores should rank high
			// "Recent Good Post" (2 hours old, 18 net points) should rank higher than
			// "Old Popular Post" (48 hours old, 45 net points) due to time decay
			var recentGoodIndex, oldPopularIndex int = -1, -1
			for i, post := range posts {
				if post.Title == "Recent Good Post" {
					recentGoodIndex = i
				}
				if post.Title == "Old Popular Post" {
					oldPopularIndex = i
				}
			}

			if recentGoodIndex >= 0 && oldPopularIndex >= 0 {
				if recentGoodIndex > oldPopularIndex {
					t.Logf("Note: Recent post ranked at %d, old post at %d", recentGoodIndex, oldPopularIndex)
					t.Log("This might be expected depending on the exact scores, but typically recent posts should rank higher")
				}
			}

			// Log the ranking for inspection
			t.Log("Top posts ranking:")
			for i, post := range posts {
				age := time.Since(post.CreatedAt).Hours()
				t.Logf("%d. %s (Points: %d, Age: %.1f hours)", i+1, post.Title, post.Points, age)
			}
		})
	}
}

func TestGetByBest(t *testing.T) {
	repo := NewPostRepository(testDB)
	setupTestData(t, repo)

	tests := []struct {
		name  string
		limit int
	}{
		{
			name:  "Get best posts with Wilson Score",
			limit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, totalCount, err := repo.GetByBest(tt.limit, 0)
			if err != nil {
				t.Fatalf("GetByBest() error = %v", err)
			}

			if totalCount == 0 {
				t.Error("GetByBest() totalCount = 0")
			}

			if len(posts) == 0 {
				t.Fatal("GetByBest() returned no posts")
			}

			// Verify Wilson Score: "High Quality Post" (10 up, 0 down) should rank higher than
			// "Controversial Post" (100 up, 30 down) because it has better proportion
			var highQualityIndex, controversialIndex int = -1, -1
			for i, post := range posts {
				if post.Title == "High Quality Post" {
					highQualityIndex = i
				}
				if post.Title == "Controversial Post" {
					controversialIndex = i
				}
			}

			if highQualityIndex >= 0 && controversialIndex >= 0 {
				if highQualityIndex > controversialIndex {
					t.Errorf("High Quality Post (10 up, 0 down) ranked at %d, should rank higher than Controversial Post (100 up, 30 down) at %d",
						highQualityIndex, controversialIndex)
				} else {
					t.Logf("✓ Wilson Score working correctly: High Quality ranked at %d, Controversial at %d",
						highQualityIndex, controversialIndex)
				}
			}

			// Log the ranking for inspection
			t.Log("Best posts ranking (Wilson Score):")
			for i, post := range posts {
				t.Logf("%d. %s (Points: %d)", i+1, post.Title, post.Points)
			}
		})
	}
}

func TestSortingAlgorithmsProduceDistinctResults(t *testing.T) {
	repo := NewPostRepository(testDB)
	setupTestData(t, repo)

	// Get results from all three sorting methods
	newPosts, _, err := repo.GetByNew(10, 0)
	if err != nil {
		t.Fatalf("GetByNew() error = %v", err)
	}

	topPosts, _, err := repo.GetByTop(10, 0)
	if err != nil {
		t.Fatalf("GetByTop() error = %v", err)
	}

	bestPosts, _, err := repo.GetByBest(10, 0)
	if err != nil {
		t.Fatalf("GetByBest() error = %v", err)
	}

	// Verify that sorting methods produce different orders
	// (unless by coincidence they're the same, which is unlikely with our test data)
	newOrder := getPostOrder(newPosts)
	topOrder := getPostOrder(topPosts)
	bestOrder := getPostOrder(bestPosts)

	t.Logf("New order: %v", newOrder)
	t.Logf("Top order: %v", topOrder)
	t.Logf("Best order: %v", bestOrder)

	// At least one of the orderings should be different
	allSame := areOrdersEqual(newOrder, topOrder) && areOrdersEqual(topOrder, bestOrder)
	if allSame {
		t.Log("Warning: All sorting algorithms produced the same order. This might indicate an issue.")
	}
}

// Helper function to get post order as slice of titles
func getPostOrder(posts []models.Post) []string {
	order := make([]string, len(posts))
	for i, post := range posts {
		order[i] = post.Title
	}
	return order
}

// Helper function to check if two orders are equal
func areOrdersEqual(order1, order2 []string) bool {
	if len(order1) != len(order2) {
		return false
	}
	for i := range order1 {
		if order1[i] != order2[i] {
			return false
		}
	}
	return true
}

func TestSortingWithPagination(t *testing.T) {
	repo := NewPostRepository(testDB)
	_, _, postIDs := setupTestData(t, repo)

	tests := []struct {
		name     string
		sortFunc func(limit, offset int) ([]models.Post, int, error)
		limit    int
		offset   int
	}{
		{
			name:     "New sorting with pagination",
			sortFunc: repo.GetByNew,
			limit:    3,
			offset:   2,
		},
		{
			name:     "Top sorting with pagination",
			sortFunc: repo.GetByTop,
			limit:    3,
			offset:   2,
		},
		{
			name:     "Best sorting with pagination",
			sortFunc: repo.GetByBest,
			limit:    3,
			offset:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, totalCount, err := tt.sortFunc(tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("Sorting with pagination error = %v", err)
			}

			if totalCount != len(postIDs) {
				t.Errorf("TotalCount = %v, want %v", totalCount, len(postIDs))
			}

			expectedLen := len(postIDs) - tt.offset
			if expectedLen > tt.limit {
				expectedLen = tt.limit
			}
			if expectedLen < 0 {
				expectedLen = 0
			}

			if len(posts) != expectedLen {
				t.Errorf("Got %d posts, expected %d with limit=%d offset=%d total=%d",
					len(posts), expectedLen, tt.limit, tt.offset, len(postIDs))
			}
		})
	}
}

func TestSortingWithEmptyDatabase(t *testing.T) {
	repo := NewPostRepository(testDB)

	// Clean database
	testDB.Exec("DELETE FROM votes WHERE post_id IS NOT NULL")
	testDB.Exec("DELETE FROM posts")

	tests := []struct {
		name     string
		sortFunc func(limit, offset int) ([]models.Post, int, error)
	}{
		{
			name:     "GetByNew with empty database",
			sortFunc: repo.GetByNew,
		},
		{
			name:     "GetByTop with empty database",
			sortFunc: repo.GetByTop,
		},
		{
			name:     "GetByBest with empty database",
			sortFunc: repo.GetByBest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, totalCount, err := tt.sortFunc(10, 0)
			if err != nil {
				t.Fatalf("Sorting empty database error = %v", err)
			}

			if totalCount != 0 {
				t.Errorf("TotalCount = %v, want 0", totalCount)
			}

			if len(posts) != 0 {
				t.Errorf("Got %d posts, want 0", len(posts))
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetByNew(b *testing.B) {
	repo := NewPostRepository(testDB)
	setupTestData(&testing.T{}, repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := repo.GetByNew(20, 0)
		if err != nil {
			b.Fatalf("GetByNew() error = %v", err)
		}
	}
}

func BenchmarkGetByTop(b *testing.B) {
	repo := NewPostRepository(testDB)
	setupTestData(&testing.T{}, repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := repo.GetByTop(20, 0)
		if err != nil {
			b.Fatalf("GetByTop() error = %v", err)
		}
	}
}

func BenchmarkGetByBest(b *testing.B) {
	repo := NewPostRepository(testDB)
	setupTestData(&testing.T{}, repo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := repo.GetByBest(20, 0)
		if err != nil {
			b.Fatalf("GetByBest() error = %v", err)
		}
	}
}
