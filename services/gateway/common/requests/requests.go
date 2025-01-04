package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sweetloveinyourheart/exploding-kittens/services/gateway/common/validations"
)

// GetVar extracts a variable from the request and returns an error if not found.
func GetVar(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	value, exists := vars[key]
	if !exists {
		return "", errors.New("missing required parameter: " + key)
	}
	return value, nil
}

// GetVars extracts multiple variables from the request.
// It returns a map of found variables and an error for any missing ones.
func GetVars(r *http.Request, keys ...string) (map[string]string, error) {
	vars := mux.Vars(r)
	result := make(map[string]string)
	for _, key := range keys {
		value, exists := vars[key]
		if !exists {
			return nil, errors.New("missing required parameter: " + key)
		}
		result[key] = value
	}
	return result, nil
}

// ParseBody reads and parses the body into a struct of any type
func ParseBody[T any](r *http.Request) (*T, error) {
	var schema T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	// Parse JSON into the struct
	if err := json.Unmarshal(body, &schema); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return &schema, nil
}

// ParseBodyWithValidation reads, parses the body into a struct of any type, and validates it
func ParseBodyWithValidation[T any](r *http.Request) (*T, error) {
	var schema T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	// Parse JSON into the struct
	if err := json.Unmarshal(body, &schema); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate the parsed struct
	if err := validations.Validate(schema); err != nil {
		return nil, err
	}

	return &schema, nil
}
