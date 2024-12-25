package requests

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
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
