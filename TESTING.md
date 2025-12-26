# Testing Documentation

## Test Coverage

The project follows Test-Driven Development (TDD) practices with comprehensive unit tests for all core packages.

### Current Coverage

| Package | Coverage | Test Files |
|---------|----------|------------|
| `pkg/auth` | 84.2% | `password_test.go`, `jwt_test.go` |
| `internal/utils` | 100.0% | `auth_test.go`, `request_test.go`, `response_test.go` |

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

## Future Test Plans

### Repository Layer Tests
- Database mocking
- CRUD operations
- Transaction handling
- Error scenarios

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
- ⏳ Integration tests (planned)
- ⏳ Handler tests (planned)
- ⏳ Repository tests (planned)

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
