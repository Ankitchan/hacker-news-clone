# Hacker News Clone - Backend

A fully functional Hacker News clone backend built with Go and PostgreSQL.

## Features

✅ **User Authentication**
- Signup with username, email, and password
- Secure password hashing with bcrypt (automatic salt generation)
- JWT-based authentication
- Token refresh mechanism

✅ **Posts**
- Create URL or text posts (or both)
- Pagination support
- **Advanced Sorting Algorithms:**
  - **New**: Chronological order (newest first)
  - **Top**: Hacker News ranking formula `(P-1)^0.8 / (T+2)^1.8`
    - Balances popularity with freshness
    - Time decay ensures front page stays current
  - **Best**: Wilson Score Interval (95% confidence)
    - Statistical quality ranking
    - 10 upvotes/0 downvotes ranks higher than 100 up/30 down
- Search functionality
- Edit and delete own posts

✅ **Comments**
- Threaded comment system
- Automatic depth calculation
- Nested comment structure
- Edit and delete own comments

✅ **Voting**
- Upvote/downvote on posts and comments
- Automatic points calculation
- Vote tracking per user

✅ **Security & Protection**
- Rate limiting middleware (DDoS protection)
- IP-based request throttling
- Token bucket algorithm (10 req/sec, burst of 20)
- Automatic cleanup to prevent memory leaks

## Tech Stack

- **Language:** Go 1.25.5
- **Database:** PostgreSQL 18.1
- **Router:** Gorilla Mux
- **Authentication:** JWT (golang-jwt/jwt/v5)
- **Password Hashing:** bcrypt
- **CORS:** rs/cors
- **Rate Limiting:** golang.org/x/time/rate

## Project Structure

```
backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Auth, CORS, rate limiting middleware
│   ├── models/          # Data models
│   ├── repository/      # Database operations
│   ├── routes/          # API route definitions
│   └── utils/           # Helper functions
├── migrations/          # Database migration files
└── pkg/
    ├── auth/            # JWT and password utilities
    └── database/        # Database connection and migrations
```

## Setup

### Prerequisites

- Go 1.25.5 or higher
- PostgreSQL 18.1 or higher

### Installation

1. **Clone the repository**
   ```bash
   cd hacker_news_clone/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup PostgreSQL**
   ```bash
   # Start PostgreSQL
   sudo systemctl start postgresql

   # Create database
   sudo -u postgres createdb hacker_news_db
   ```

4. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. **Build the application**
   ```bash
   go build -o bin/api ./cmd/api
   ```

6. **Run the server**
   ```bash
   ./bin/api
   ```

The server will start on `http://localhost:8080` (configurable via PORT environment variable).

## Environment Variables

See [.env.example](.env.example) for all configuration options:

- `PORT` - Server port (default: 8080)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - Secret key for JWT tokens (**CHANGE IN PRODUCTION**)
- `JWT_EXPIRATION_HOURS` - Token expiration time
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins

## API Documentation

See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for complete API reference.

### Quick Examples

**Signup:**
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'
```

**Create Post:**
```bash
curl -X POST http://localhost:8080/api/posts \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{"title":"My Post","url":"https://example.com"}'
```

**Vote on Post:**
```bash
curl -X POST http://localhost:8080/api/votes/posts/1 \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -d '{"vote_type":1}'
```

**Get Posts:**
```bash
# Get newest posts
curl http://localhost:8080/api/posts?sort=new&page=1&page_size=20

# Get top posts (HN algorithm)
curl http://localhost:8080/api/posts?sort=top&page=1&page_size=20

# Get best posts (Wilson Score)
curl http://localhost:8080/api/posts?sort=best&page=1&page_size=20
```

## Sorting Algorithms

### Top - Hacker News Ranking Formula

Posts are ranked using the authentic Hacker News algorithm:

```
Score = (P - 1)^0.8 / (T + 2)^1.8
```

Where:
- **P** = Points (upvotes - downvotes)
- **T** = Age in hours
- **-1** = Ignores the submitter's automatic upvote
- **+2** = Prevents brand new posts from instantly topping the list
- **0.8 exponent on points** = Diminishing returns for additional votes
- **1.8 exponent on time** = Gravity factor - ensures old posts decay

**Key behavior**: Because 1.8 > 0.8, time eventually wins. Recent posts with moderate scores rank higher than old posts with many votes, keeping the front page fresh.

### Best - Wilson Score Interval

Posts are ranked using the Wilson Score confidence interval (95% confidence level):

```
Lower bound = (phat + z²/2n - z√[(phat(1-phat) + z²/4n)/n]) / (1 + z²/n)
```

Where:
- **phat** = upvotes / total_votes (proportion of positive votes)
- **z** = 1.96 (95% confidence z-score)
- **n** = total votes

**Key behavior**: This statistical approach calculates the minimum quality score a post is likely to have. A post with **10 upvotes and 0 downvotes** ranks **higher** than one with **100 upvotes and 30 downvotes**, because it has a better quality ratio with statistical confidence.

## Database Schema

The application automatically runs migrations on startup. Schema includes:

- **users** - User accounts
- **posts** - Posts with URL or text content
- **comments** - Threaded comments with automatic depth
- **votes** - Voting system for posts and comments
- **schema_migrations** - Migration tracking

## Security Features

- **Password hashing** with bcrypt (automatic unique salt per password)
- **JWT token-based authentication** with refresh mechanism
- **Rate limiting** to prevent DDoS attacks (10 req/sec per IP, burst of 20)
- **CORS protection** with configurable origins
- **SQL injection prevention** (parameterized queries)
- **Authorization checks** (users can only edit/delete their own content)
- **Proxy-aware IP detection** (X-Forwarded-For, X-Real-IP headers)

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests verbosely
go test ./... -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Coverage

- **pkg/auth**: 84.2% coverage
  - Password hashing and verification
  - JWT token generation and validation
  - Token expiration and refresh

- **internal/utils**: 100% coverage
  - Request/response helpers
  - Context utilities
  - Parameter extraction

- **internal/repository**: Comprehensive sorting algorithm tests
  - GetByNew (chronological sorting)
  - GetByTop (HN ranking formula)
  - GetByBest (Wilson Score Interval)
  - Pagination and edge cases
  - Performance benchmarks

See [TESTING.md](TESTING.md) for detailed testing documentation and [TEST_SUMMARY.md](TEST_SUMMARY.md) for latest test results.

### Manual Testing

The API has been manually tested with:
- User registration and login
- Post creation (URL and text)
- Comment creation (flat and threaded)
- Voting system
- Pagination
- Authorization

## Development

**Build:**
```bash
go build -o bin/api ./cmd/api
```

**Run with live reload (using air):**
```bash
air
```

**Format code:**
```bash
go fmt ./...
```

## Production Considerations

Before deploying to production:

1. Change `JWT_SECRET` to a strong random value
2. Use a strong database password
3. Enable SSL for PostgreSQL (`DB_SSL_MODE=require`)
4. Set appropriate CORS origins
5. Use HTTPS for the API
6. Review rate limiting settings (currently 10 req/sec, burst 20)
7. Set up proper logging and monitoring
8. Configure reverse proxy headers for accurate IP detection

### Rate Limiting Configuration

The application includes built-in rate limiting with the following defaults:
- **Rate:** 10 requests per second per IP
- **Burst:** 20 requests (allows temporary spikes)
- **Algorithm:** Token bucket
- **Cleanup:** Automatic every 5 minutes

To adjust rate limits, modify the values in `internal/routes/routes.go`:
```go
rateLimiter := middleware.NewRateLimiter(rate.Limit(10), 20)
```

## License

MIT
