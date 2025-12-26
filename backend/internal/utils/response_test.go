package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{
			name:    "bad request error",
			code:    http.StatusBadRequest,
			message: "Invalid input",
		},
		{
			name:    "unauthorized error",
			code:    http.StatusUnauthorized,
			message: "Not authorized",
		},
		{
			name:    "internal server error",
			code:    http.StatusInternalServerError,
			message: "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondWithError(w, tt.code, tt.message)

			if w.Code != tt.code {
				t.Errorf("RespondWithError() status = %v, want %v", w.Code, tt.code)
			}

			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Error != http.StatusText(tt.code) {
				t.Errorf("RespondWithError() error = %v, want %v", response.Error, http.StatusText(tt.code))
			}

			if response.Message != tt.message {
				t.Errorf("RespondWithError() message = %v, want %v", response.Message, tt.message)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondWithError() Content-Type = %v, want application/json", contentType)
			}
		})
	}
}

func TestRespondWithJSON(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name    string
		code    int
		payload interface{}
	}{
		{
			name: "simple struct",
			code: http.StatusOK,
			payload: TestData{
				Name:  "test",
				Value: 123,
			},
		},
		{
			name:    "map",
			code:    http.StatusCreated,
			payload: map[string]string{"key": "value"},
		},
		{
			name:    "slice",
			code:    http.StatusOK,
			payload: []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondWithJSON(w, tt.code, tt.payload)

			if w.Code != tt.code {
				t.Errorf("RespondWithJSON() status = %v, want %v", w.Code, tt.code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondWithJSON() Content-Type = %v, want application/json", contentType)
			}

			// Verify it's valid JSON
			var result interface{}
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Errorf("RespondWithJSON() returned invalid JSON: %v", err)
			}
		})
	}
}

func TestRespondWithSuccess(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		data    interface{}
		message string
	}{
		{
			name:    "success with data and message",
			code:    http.StatusOK,
			data:    map[string]string{"result": "success"},
			message: "Operation completed",
		},
		{
			name:    "success with nil data",
			code:    http.StatusOK,
			data:    nil,
			message: "Deleted successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			RespondWithSuccess(w, tt.code, tt.data, tt.message)

			if w.Code != tt.code {
				t.Errorf("RespondWithSuccess() status = %v, want %v", w.Code, tt.code)
			}

			var response SuccessResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if !response.Success {
				t.Error("RespondWithSuccess() success = false, want true")
			}

			if response.Message != tt.message {
				t.Errorf("RespondWithSuccess() message = %v, want %v", response.Message, tt.message)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondWithSuccess() Content-Type = %v, want application/json", contentType)
			}
		})
	}
}
