package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	ExpirationHours  int
}

var jwtConfig *JWTConfig

// InitJWT initializes the JWT configuration
func InitJWT(secret string, expirationHours int) {
	jwtConfig = &JWTConfig{
		Secret:          secret,
		ExpirationHours: expirationHours,
	}
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int, username, email string) (string, error) {
	if jwtConfig == nil {
		return "", fmt.Errorf("JWT configuration not initialized")
	}

	// Create claims with user information and standard claims
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(jwtConfig.ExpirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hackernews-clone",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	if jwtConfig == nil {
		return nil, fmt.Errorf("JWT configuration not initialized")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RefreshToken generates a new token for an existing valid token
func RefreshToken(tokenString string) (string, error) {
	// Validate the existing token
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate token for refresh: %w", err)
	}

	// Generate a new token with the same user information
	return GenerateToken(claims.UserID, claims.Username, claims.Email)
}
