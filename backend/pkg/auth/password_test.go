package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "password with special characters",
			password: "P@ssw0rd!#$%",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plaintext password")
				}
				// Verify the hash starts with bcrypt prefix
				if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
					t.Error("HashPassword() did not return valid bcrypt hash")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name          string
		hashedPassword string
		password      string
		wantErr       bool
	}{
		{
			name:          "correct password",
			hashedPassword: hash,
			password:      password,
			wantErr:       false,
		},
		{
			name:          "incorrect password",
			hashedPassword: hash,
			password:      "wrongpassword",
			wantErr:       true,
		},
		{
			name:          "empty password",
			hashedPassword: hash,
			password:      "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password - 8 characters",
			password: "password",
			wantErr:  false,
		},
		{
			name:     "valid password - with numbers and special chars",
			password: "Pass123!@#",
			wantErr:  false,
		},
		{
			name:     "too short - 7 characters",
			password: "pass123",
			wantErr:  true,
		},
		{
			name:     "too long - 73 characters",
			password: strings.Repeat("a", 73),
			wantErr:  true,
		},
		{
			name:     "exactly 72 characters",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Each hash should be different due to unique salt
	if hash1 == hash2 {
		t.Error("HashPassword() should generate unique hashes with different salts")
	}

	// But both should verify correctly
	if err := VerifyPassword(hash1, password); err != nil {
		t.Error("First hash should verify correctly")
	}
	if err := VerifyPassword(hash2, password); err != nil {
		t.Error("Second hash should verify correctly")
	}
}
