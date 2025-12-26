package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password with salt
// bcrypt automatically handles salt generation and storage within the hash
func HashPassword(password string) (string, error) {
	// Cost factor of 12 provides a good balance between security and performance
	// bcrypt automatically generates a unique salt for each password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// VerifyPassword compares a plain text password with a hashed password
// bcrypt automatically extracts the salt from the hash and compares
func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

// ValidatePasswordStrength checks if password meets minimum requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 72 {
		return fmt.Errorf("password must not exceed 72 characters")
	}

	// You can add more validation rules here:
	// - Check for uppercase letters
	// - Check for lowercase letters
	// - Check for numbers
	// - Check for special characters

	return nil
}
