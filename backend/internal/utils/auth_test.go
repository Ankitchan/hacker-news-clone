package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ankitchan/hackernews-clone/pkg/auth"
)

func TestGetUserFromContext(t *testing.T) {
	tests := []struct {
		name      string
		claims    *auth.Claims
		wantFound bool
	}{
		{
			name: "valid claims in context",
			claims: &auth.Claims{
				UserID:   1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantFound: true,
		},
		{
			name:      "no claims in context",
			claims:    nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.claims != nil {
				ctx = context.WithValue(ctx, UserContextKey, tt.claims)
			}

			claims, found := GetUserFromContext(ctx)
			if found != tt.wantFound {
				t.Errorf("GetUserFromContext() found = %v, want %v", found, tt.wantFound)
			}

			if tt.wantFound && claims == nil {
				t.Error("GetUserFromContext() returned nil claims when expected")
			}

			if tt.wantFound && claims != nil {
				if claims.UserID != tt.claims.UserID {
					t.Errorf("GetUserFromContext() UserID = %v, want %v", claims.UserID, tt.claims.UserID)
				}
			}
		})
	}
}

func TestGetUserFromRequest(t *testing.T) {
	tests := []struct {
		name      string
		claims    *auth.Claims
		wantFound bool
	}{
		{
			name: "valid claims in request context",
			claims: &auth.Claims{
				UserID:   42,
				Username: "alice",
				Email:    "alice@example.com",
			},
			wantFound: true,
		},
		{
			name:      "no claims in request context",
			claims:    nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.claims != nil {
				ctx := context.WithValue(req.Context(), UserContextKey, tt.claims)
				req = req.WithContext(ctx)
			}

			claims, found := GetUserFromRequest(req)
			if found != tt.wantFound {
				t.Errorf("GetUserFromRequest() found = %v, want %v", found, tt.wantFound)
			}

			if tt.wantFound && claims == nil {
				t.Error("GetUserFromRequest() returned nil claims when expected")
			}

			if tt.wantFound && claims != nil {
				if claims.UserID != tt.claims.UserID {
					t.Errorf("GetUserFromRequest() UserID = %v, want %v", claims.UserID, tt.claims.UserID)
				}
				if claims.Username != tt.claims.Username {
					t.Errorf("GetUserFromRequest() Username = %v, want %v", claims.Username, tt.claims.Username)
				}
				if claims.Email != tt.claims.Email {
					t.Errorf("GetUserFromRequest() Email = %v, want %v", claims.Email, tt.claims.Email)
				}
			}
		})
	}
}

func TestUserContextKey(t *testing.T) {
	// Test that the context key is of the correct type
	var key ContextKey = UserContextKey
	if key != "user" {
		t.Errorf("UserContextKey = %v, want 'user'", key)
	}
}
