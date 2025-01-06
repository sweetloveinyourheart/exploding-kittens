package actions

import (
	"context"

	"github.com/samber/do"

	userServerConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"
)

type actions struct {
	context context.Context

	userServerClient userServerConnect.UserServerClient
}

func NewActions(ctx context.Context, signingToken string) *actions {
	return &actions{
		context:          ctx,
		userServerClient: do.MustInvoke[userServerConnect.UserServerClient](nil),
	}
}
