package middlewares

import (
	"net/http"

	"github.com/sweetloveinyourheart/exploding-kittens/services/gateway/common/errors"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errors.GlobalErrorHandler(w, r, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
