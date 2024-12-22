package middlewares

import (
	"net/http"

	gateway_errors "github.com/sweetloveinyourheart/planning-pocker/services/gateway/common/errors"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				gateway_errors.GlobalErrorHandler(w, r, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
