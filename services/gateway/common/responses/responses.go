package responses

import (
	"encoding/json"
	"net/http"
	"strings"
)

type AppResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data any) {
	response := AppResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ###
// ### Success responses
// ###

func Ok(w http.ResponseWriter, data any, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Success"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, http.StatusOK, msg, data)
}

func Created(w http.ResponseWriter, data any, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Created"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, http.StatusCreated, msg, data)
}

// ###
// ### Client error responses
// ###

// BadRequestException sends a 400 JSON response with optional error messages.
func BadRequestException(w http.ResponseWriter, errors ...error) {
	// Default message if no custom messages are provided
	msg := "Bad request"
	if len(errors) > 0 {
		// Combine error messages into a single string
		msg = strings.Join(getErrorMessages(errors), "; ")
	}

	JSON(w, http.StatusBadRequest, msg, nil)
}

// UnauthorizedException sends a 401 JSON response with optional error messages.
func UnauthorizedException(w http.ResponseWriter, errors ...error) {
	// Default message if no custom messages are provided
	msg := "Unauthorized"
	if len(errors) > 0 {
		// Combine error messages into a single string
		msg = strings.Join(getErrorMessages(errors), "; ")
	}

	JSON(w, http.StatusUnauthorized, msg, nil)
}

// ForbiddenException sends a 403 JSON response with optional error messages.
func ForbiddenException(w http.ResponseWriter, errors ...error) {
	// Default message if no custom messages are provided
	msg := "Forbidden resource"
	if len(errors) > 0 {
		// Combine error messages into a single string
		msg = strings.Join(getErrorMessages(errors), "; ")
	}

	JSON(w, http.StatusForbidden, msg, nil)
}

// NotFoundException sends a 404 JSON response with optional error messages.
func NotFoundException(w http.ResponseWriter, errors ...error) {
	// Default message if no custom messages are provided
	msg := "Resource not found"
	if len(errors) > 0 {
		// Combine error messages into a single string
		msg = strings.Join(getErrorMessages(errors), "; ")
	}

	JSON(w, http.StatusNotFound, msg, nil)
}

// getErrorMessages extracts the error messages from a slice of errors.
func getErrorMessages(errors []error) []string {
	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}
	return messages
}
