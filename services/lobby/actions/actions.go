package actions

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)
}

func NewActions(ctx context.Context, signingToken string) *actions {
	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
	}
}
