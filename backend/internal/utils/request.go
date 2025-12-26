package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ParseJSONBody parses JSON request body
func ParseJSONBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// GetIDParam extracts and parses ID from URL parameter
func GetIDParam(r *http.Request, paramName string) (int, error) {
	vars := mux.Vars(r)
	idStr := vars[paramName]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter")
	}

	if id <= 0 {
		return 0, fmt.Errorf("ID must be positive")
	}

	return id, nil
}

// GetQueryParam gets a query parameter value
func GetQueryParam(r *http.Request, paramName string, defaultValue string) string {
	value := r.URL.Query().Get(paramName)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt gets a query parameter as integer
func GetQueryParamInt(r *http.Request, paramName string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(paramName)
	if value == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", paramName)
	}

	return intValue, nil
}
