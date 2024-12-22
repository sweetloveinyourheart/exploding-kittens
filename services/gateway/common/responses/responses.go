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

func BadRequestException(w http.ResponseWriter, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Bad request"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, 400, msg, nil)
}

func UnAuthorizedException(w http.ResponseWriter, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Unauthorized"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, 401, msg, nil)
}

func ForbiddenException(w http.ResponseWriter, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Forbidden resource"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, 403, msg, nil)
}

func NotFoundException(w http.ResponseWriter, messages ...string) {
	// Default message if no custom messages are provided
	msg := "Resource not found"
	if len(messages) > 0 {
		msg = strings.Join(messages, "; ")
	}

	JSON(w, 404, msg, nil)
}
