package nats

import (
	"context"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/cockroachdb/errors"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
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

func BusErrors(ctx context.Context, eventBus *EventBus) {
	busName := eventBus.busName
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = eventBus.Close()
				return
			case err, ok := <-eventBus.Errors():
				HandleError(ctx, errors.Wrap(err, busName))
				if !ok {
					return
				}
			}
		}
	}()
}
