package nats

import (
	"context"
	"os"
	"strings"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"go.uber.org/zap"
)

func HandleError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "Server Shutdown") {
		if strings.EqualFold(os.Getenv("GO_ENV"), "TEST") {
			log.Global().DebugContext(ctx, "eventing: consume error", zap.Error(err))
		} else {
			log.Global().WarnContext(ctx, "eventing: consume error", zap.Error(err))
		}
		return
	}
	if strings.Contains(err.Error(), "consumer deleted") {
		log.Global().InfoContext(ctx, "nats(bus) error", zap.String("type", "nats"), zap.Error(err))
	} else {
		log.Global().ErrorContext(ctx, "nats(bus) error", zap.String("type", "nats"), zap.Error(err))
	}
}
