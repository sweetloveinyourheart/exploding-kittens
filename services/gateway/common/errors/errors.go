package errors

import (
	"net/http"

	"github.com/cockroachdb/errors"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func GlobalErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	if ok := errors.As(err, &appErr); ok {
		http.Error(w, appErr.Message, appErr.Code)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
