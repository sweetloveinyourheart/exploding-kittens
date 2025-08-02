package actions

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/repos"
)

const (
	// DefaultCacheTTL is the default time-to-live for cached data.
	DefaultCacheTTL = 5 * time.Minute
	GlobalCacheKey  = "global"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)

	redisClient *redis.Client
	cacheTTL    time.Duration

	cardRepo repos.ICardRepository

	tracer trace.Tracer
}

func NewActions(ctx context.Context, signingToken string, redisClient *redis.Client) *actions {
	cardRepo := do.MustInvoke[repos.ICardRepository](nil)

	tracer := otel.Tracer("com.sweetloveinyourheart.kittens.dataprovider.actions")

	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
		redisClient: redisClient,
		cardRepo:    cardRepo,
		tracer:      tracer,
	}
}
