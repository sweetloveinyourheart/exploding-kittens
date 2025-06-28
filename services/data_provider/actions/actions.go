package actions

import (
	"context"

	"github.com/samber/do"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)

	cardRepo repos.ICardRepository

	tracer trace.Tracer
}

func NewActions(ctx context.Context, signingToken string) *actions {
	cardRepo := do.MustInvoke[repos.ICardRepository](nil)

	tracer := otel.Tracer("com.sweetloveinyourheart.kittens.dataprovider.actions")

	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
		cardRepo:    cardRepo,
		tracer:      tracer,
	}
}
