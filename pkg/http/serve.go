package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
)

func ServeHTTP(ctx context.Context, handler *mux.Router, port uint64, serviceName string) {
	log.GlobalSugared().Infof("HTTP %s listening on port %d\n", serviceName, port)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.GlobalSugared().Panicf("%s failed to serve: %v", serviceName, err)
			} else {
				log.GlobalSugared().Infof("HTTP %s server closed", serviceName)
			}
		}
	}()
	<-ctx.Done()
	log.GlobalSugared().Infof("HTTP %s shutting down", serviceName)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.GlobalSugared().Panicf("%s failed to shutdown: %v", serviceName, err)
	}
}
