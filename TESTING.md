# Testing Documentation

## Test Coverage

The project follows Test-Driven Development (TDD) practices with comprehensive unit tests for all core packages.

### Current Coverage

| Package | Coverage | Test Files |
|---------|----------|------------|
| `pkg/auth` | 84.2% | `password_test.go`, `jwt_test.go` |
| `internal/utils` | 100.0% | `auth_test.go`, `request_test.go`, `response_test.go` |
| `internal/repository` | ✅ Tested | `post_repository_sorting_test.go` |

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Run Tests with Detailed Coverage Report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Specific Package Tests
```bash
# Auth package
go test ./pkg/auth/... -v

# Utils package
go test ./internal/utils/... -v
```

### Run Tests with Race Detection
```bash
go test ./... -race
```

## Test Structure

### Auth Package Tests (`pkg/auth`)

#### Password Tests (`password_test.go`)
- ✅ `TestHashPassword` - Validates password hashing
  - Valid passwords
  - Long passwords (up to 72 chars)
  - Special characters
- ✅ `TestVerifyPassword` - Validates password verification
  - Correct password
  - Incorrect password
  - Empty password
- ✅ `TestValidatePasswordStrength` - Validates password requirements
  - Minimum length (8 chars)
  - Maximum length (72 chars)
  - Empty passwords
- ✅ `TestHashPasswordUniqueness` - Ensures unique salts
  - Same password generates different hashes
  - Both hashes verify correctly

#### JWT Tests (`jwt_test.go`)
- ✅ `TestInitJWT` - JWT configuration initialization
- ✅ `TestGenerateToken` - Token generation
  - Valid user data
  - Special characters in username/email
- ✅ `TestValidateToken` - Token validation
  - Valid tokens
  - Invalid tokens
  - Empty tokens
  - Malformed tokens
- ✅ `TestTokenExpiration` - Expiration handling
- ✅ `TestRefreshToken` - Token refresh mechanism
- ✅ `TestValidateTokenWithWrongSecret` - Security validation

### Utils Package Tests (`internal/utils`)

#### Auth Utils Tests (`auth_test.go`)
- ✅ `TestGetUserFromContext` - Context extraction
- ✅ `TestGetUserFromRequest` - HTTP request context
- ✅ `TestUserContextKey` - Context key validation

#### Request Utils Tests (`request_test.go`)
- ✅ `TestParseJSONBody` - JSON parsing
  - Valid JSON
  - Invalid JSON
  - Empty body
- ✅ `TestGetIDParam` - URL parameter extraction
  - Valid IDs
  - Invalid IDs (non-numeric, negative, zero)
- ✅ `TestGetQueryParam` - Query parameter extraction
- ✅ `TestGetQueryParamInt` - Integer query parameters

#### Response Utils Tests (`response_test.go`)
- ✅ `TestRespondWithError` - Error responses
  - Different status codes
  - Error messages
  - Content-Type headers
- ✅ `TestRespondWithJSON` - JSON responses
  - Structs, maps, slices
  - Status codes
- ✅ `TestRespondWithSuccess` - Success responses

## Test Patterns

### Table-Driven Tests
All tests use table-driven test patterns for comprehensive coverage:

```go
tests := []struct {
    name     string
    input    string
    want     result
    wantErr  bool
}{
    {
        name:    "valid input",
        input:   "test",
        want:    expected,
        wantErr: false,
    },
    // ... more test cases
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### Security Tests
- Password hashing with unique salts
- JWT token validation with wrong secrets
- Token expiration
- Authorization context validation

### Edge Cases
- Empty inputs
- Invalid formats
- Boundary conditions (min/max lengths)
- Nil values

## Continuous Integration

### Pre-commit Checks
```bash
#!/bin/bash
# Run before committing

# Format code
go fmt ./...

# Run tests
go test ./... -cover

# Run linter
golint ./...

# Check for vulnerabilities
go vet ./...
```

### CI Pipeline Recommendations
1. Run tests on multiple Go versions
2. Check test coverage (aim for >80%)
3. Run race detector
4. Run static analysis tools
5. Generate coverage reports

## Writing New Tests

### Guidelines
1. Follow table-driven test pattern
2. Test both success and failure cases
3. Test edge cases and boundary conditions
4. Use descriptive test names
5. Keep tests independent and isolated
6. Mock external dependencies

### Example Test
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "success case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "error case",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewFeature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Repository Package Tests (`internal/repository`)

#### Sorting Algorithm Tests (`post_repository_sorting_test.go`)

Comprehensive integration tests for all three sorting algorithms with real database operations.

##### Test Suite Overview

- ✅ **TestGetByNew** - Chronological sorting
  - Verifies posts are sorted by creation time (newest first)
  - Tests pagination (limit/offset)
  - Validates ordering consistency

- ✅ **TestGetByTop** - Hacker News ranking formula
  - Tests formula: `(P-1)^0.8 / (T+2)^1.8`
  - Validates time decay behavior
  - Confirms recent posts rank higher than old popular posts
  - Logs detailed ranking for manual inspection

- ✅ **TestGetByBest** - Wilson Score Interval
  - Tests 95% confidence interval calculation
  - Validates statistical quality ranking
  - Confirms 10 up/0 down ranks higher than 100 up/30 down
  - Verifies edge cases with no votes

- ✅ **TestSortingAlgorithmsProduceDistinctResults**
  - Ensures each algorithm produces unique orderings
  - Validates algorithm independence
  - Logs comparisons for all three methods

- ✅ **TestSortingWithPagination**
  - Tests limit/offset for all sorting methods
  - Validates correct result counts
  - Ensures pagination doesn't break ordering

- ✅ **TestSortingWithEmptyDatabase**
  - Edge case: no posts in database
  - Tests all three sorting methods
  - Validates graceful handling of empty results

##### Test Data Setup

The tests create realistic post scenarios:

```go
// Post 1: Old Popular Post
- Age: 48 hours
- Votes: 50 upvotes, 5 downvotes
- Purpose: Tests time decay in "top" algorithm

// Post 2: Recent Good Post
- Age: 2 hours
- Votes: 20 upvotes, 2 downvotes
- Purpose: Should rank high in "top" due to recency

// Post 3: Brand New Post
- Age: 30 minutes
- Votes: 5 upvotes, 0 downvotes
- Purpose: Tests freshness handling

// Post 4: High Quality Post
- Age: 6 hours
- Votes: 10 upvotes, 0 downvotes
- Purpose: Perfect ratio for Wilson Score

// Post 5: Controversial Post
- Age: 4 hours
- Votes: 100 upvotes, 30 downvotes
- Purpose: High total, lower quality ratio

// Post 6: No Votes Post
- Age: 1 hour
- Votes: 0 upvotes, 0 downvotes
- Purpose: Edge case testing
```

##### Performance Benchmarks

```bash
BenchmarkGetByNew-4     5    4.5ms/op  (baseline)
BenchmarkGetByTop-4     5    5.5ms/op  (21% slower)
BenchmarkGetByBest-4    5    7.3ms/op  (62% slower - complex Wilson Score)
```

##### Running Repository Tests

```bash
# Run all repository tests
go test ./internal/repository -v

# Run specific sorting tests
go test ./internal/repository -run TestGetBy

# Run with benchmarks
go test ./internal/repository -bench=.

# Run with race detection
go test ./internal/repository -race
```

##### Key Validations

1. **HN Formula Correctness**: Verifies recent posts with moderate scores rank higher than old posts with high scores
2. **Wilson Score Accuracy**: Confirms statistical quality ranking works as expected
3. **Ordering Consistency**: All sorts maintain stable, deterministic ordering
4. **Edge Case Handling**: Empty database, zero votes, pagination boundaries
5. **Performance**: Benchmarks ensure algorithms perform within acceptable ranges

## Future Test Plans

### Handler Tests (Planned)
- HTTP endpoint testing
- Request/response validation
- Authentication/authorization
- Error handling

### Handler Tests
- HTTP endpoint testing
- Request/response validation
- Authentication/authorization
- Error handling

### Middleware Tests
- Auth middleware
- CORS middleware
- Logging middleware

### Integration Tests
- End-to-end API testing
- Database integration
- Full request lifecycle

## Test Metrics

### Current Status
- ✅ Unit tests implemented for core packages
- ✅ 84%+ average coverage
- ✅ All edge cases covered
- ✅ Security tests included
- ✅ **Repository sorting algorithm tests (NEW)**
  - Integration tests with real database
  - All 3 sorting methods tested
  - Performance benchmarks included
- ⏳ Integration tests (planned)
- ⏳ Handler tests (planned)

## Troubleshooting

### Common Issues

**Tests timing out:**
```bash
go test ./... -timeout 30s
```

**Race conditions:**
```bash
go test ./... -race
```

**Verbose output:**
```bash
go test ./... -v
```

**Run specific test:**
```bash
go test -run TestFunctionName ./pkg/auth
```

## Best Practices

1. **Write tests first** - Follow TDD
2. **Keep tests fast** - Unit tests should complete in milliseconds
3. **Make tests deterministic** - No random behavior
4. **Use meaningful assertions** - Clear error messages
5. **Test one thing at a time** - Single responsibility
6. **Clean up resources** - Use defer for cleanup
7. **Document complex tests** - Add comments for clarity

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Comments](https://go.dev/blog/subtests)
