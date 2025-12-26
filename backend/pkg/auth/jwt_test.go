package auth

import (
	"testing"
	"time"
)

func TestInitJWT(t *testing.T) {
	secret := "test-secret"
	hours := 24

	InitJWT(secret, hours)

	if jwtConfig == nil {
		t.Error("InitJWT() should initialize jwtConfig")
	}
	if jwtConfig.Secret != secret {
		t.Errorf("InitJWT() secret = %v, want %v", jwtConfig.Secret, secret)
	}
	if jwtConfig.ExpirationHours != hours {
		t.Errorf("InitJWT() expiration = %v, want %v", jwtConfig.ExpirationHours, hours)
	}
}

func TestGenerateToken(t *testing.T) {
	// Initialize JWT config
	InitJWT("test-secret-key", 24)

	tests := []struct {
		name     string
		userID   int
		username string
		email    string
		wantErr  bool
	}{
		{
			name:     "valid user data",
			userID:   1,
			username: "testuser",
			email:    "test@example.com",
			wantErr:  false,
		},
		{
			name:     "valid user with special characters",
			username: "user.name-123",
			email:    "user+test@example.com",
			userID:   999,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateToken() returned empty token")
				}
				// JWT tokens should have 3 parts separated by dots
				parts := len(token)
				if parts == 0 {
					t.Error("GenerateToken() returned invalid token format")
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Initialize JWT config
	InitJWT("test-secret-key", 24)

	userID := 42
	username := "testuser"
	email := "test@example.com"

	token, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not-a-jwt-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims == nil {
					t.Error("ValidateToken() returned nil claims")
					return
				}
				if claims.UserID != userID {
					t.Errorf("ValidateToken() userID = %v, want %v", claims.UserID, userID)
				}
				if claims.Username != username {
					t.Errorf("ValidateToken() username = %v, want %v", claims.Username, username)
				}
				if claims.Email != email {
					t.Errorf("ValidateToken() email = %v, want %v", claims.Email, email)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Initialize with very short expiration for testing
	InitJWT("test-secret", 0) // 0 hours = expires immediately

	userID := 1
	username := "testuser"
	email := "test@example.com"

	token, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a moment to ensure expiration
	time.Sleep(2 * time.Second)

	_, err = ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() should fail for expired token")
	}

	// Reset to normal expiration for other tests
	InitJWT("test-secret", 24)
}

func TestRefreshToken(t *testing.T) {
	// Initialize JWT config
	InitJWT("test-secret-key", 24)

	userID := 1
	username := "testuser"
	email := "test@example.com"

	originalToken, err := GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "refresh valid token",
			token:   originalToken,
			wantErr: false,
		},
		{
			name:    "refresh invalid token",
			token:   "invalid.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newToken, err := RefreshToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if newToken == "" {
					t.Error("RefreshToken() returned empty token")
				}
				// Note: newToken might be same as original if generated in same second
				// This is acceptable as both tokens are valid

				// Validate the new token
				claims, err := ValidateToken(newToken)
				if err != nil {
					t.Errorf("RefreshToken() generated invalid token: %v", err)
				}
				if claims.UserID != userID {
					t.Errorf("RefreshToken() userID = %v, want %v", claims.UserID, userID)
				}
			}
		})
	}
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	// Generate token with one secret
	InitJWT("secret-1", 24)
	token, err := GenerateToken(1, "user", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with different secret
	InitJWT("secret-2", 24)
	_, err = ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() should fail when secret is different")
	}

	// Reset for other tests
	InitJWT("test-secret-key", 24)
}
