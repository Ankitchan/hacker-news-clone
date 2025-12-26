package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestParseJSONBody(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			body:    `{"name":"John","email":"john@example.com"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `{"name":"John"`,
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    `{}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			var data TestStruct
			err := ParseJSONBody(req, &data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetIDParam(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		paramValue string
		wantID    int
		wantErr   bool
	}{
		{
			name:      "valid ID",
			paramName: "id",
			paramValue: "123",
			wantID:    123,
			wantErr:   false,
		},
		{
			name:      "invalid ID - not a number",
			paramName: "id",
			paramValue: "abc",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "invalid ID - negative",
			paramName: "id",
			paramValue: "-5",
			wantID:    0,
			wantErr:   true,
		},
		{
			name:      "invalid ID - zero",
			paramName: "id",
			paramValue: "0",
			wantID:    0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.paramValue, nil)
			req = mux.SetURLVars(req, map[string]string{tt.paramName: tt.paramValue})

			id, err := GetIDParam(req, tt.paramName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id != tt.wantID {
				t.Errorf("GetIDParam() = %v, want %v", id, tt.wantID)
			}
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		paramName    string
		defaultValue string
		want         string
	}{
		{
			name:         "parameter exists",
			url:          "/test?search=golang",
			paramName:    "search",
			defaultValue: "default",
			want:         "golang",
		},
		{
			name:         "parameter missing",
			url:          "/test",
			paramName:    "search",
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "empty parameter value",
			url:          "/test?search=",
			paramName:    "search",
			defaultValue: "default",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			got := GetQueryParam(req, tt.paramName, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetQueryParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetQueryParamInt(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		paramName    string
		defaultValue int
		want         int
		wantErr      bool
	}{
		{
			name:         "valid integer",
			url:          "/test?page=5",
			paramName:    "page",
			defaultValue: 1,
			want:         5,
			wantErr:      false,
		},
		{
			name:         "parameter missing - use default",
			url:          "/test",
			paramName:    "page",
			defaultValue: 1,
			want:         1,
			wantErr:      false,
		},
		{
			name:         "invalid integer",
			url:          "/test?page=abc",
			paramName:    "page",
			defaultValue: 1,
			want:         0,
			wantErr:      true,
		},
		{
			name:         "negative integer",
			url:          "/test?page=-5",
			paramName:    "page",
			defaultValue: 1,
			want:         -5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			got, err := GetQueryParamInt(req, tt.paramName, tt.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueryParamInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetQueryParamInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
