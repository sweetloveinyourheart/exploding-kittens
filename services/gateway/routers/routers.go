package routers

import (
	"context"

	"github.com/gorilla/mux"

	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/handlers"
	"github.com/sweetloveinyourheart/planning-pocker/services/gateway/middlewares"
)

func NewGatewayRouter(ctx context.Context) *mux.Router {
	handler := handlers.NewGatewayHandler(ctx)
	mux := mux.NewRouter()

	// Apply the error handler middleware
	mux.Use(middlewares.ErrorMiddleware)

	// Path prefix v1
	router := mux.PathPrefix("/api/v1").Subrouter()

	// User routers
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{id}", handler.GetUserByID).Methods("GET")
	userRouter.HandleFunc("/create-guest-user", handler.CreateNewGuestUser).Methods("POST")

	return router
}
