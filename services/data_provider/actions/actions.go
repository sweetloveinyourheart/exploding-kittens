package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)

	cardRepo repos.ICardRepository
}

func NewActions(ctx context.Context, signingToken string) *actions {
	cardRepo := do.MustInvoke[repos.ICardRepository](nil)

	return &actions{
		context:     ctx,
		defaultAuth: interceptors.ConnectServerAuthHandler(signingToken),
		cardRepo:    cardRepo,
	}
}
