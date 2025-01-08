package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go/grpcconnect"
	userServerConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)

	userServerClient userServerConnect.UserServerClient
}

// AuthFuncOverride is a callback function that overrides the default authorization middleware in the GRPC layer. This is
// used to allow unauthenticated endpoints (such as login) to be called without a token.
func (a *actions) AuthFuncOverride(ctx context.Context, token string, fullMethodName string) (context.Context, error) {
	if fullMethodName == grpcconnect.ClientServerCreateNewGuestUserProcedure {
		return ctx, nil
	}

	if fullMethodName == grpcconnect.ClientServerGuestLoginProcedure {
		return ctx, nil
	}

	return a.defaultAuth(ctx, token)
}

func NewActions(ctx context.Context, signingToken string) *actions {
	return &actions{
		context:          ctx,
		defaultAuth:      interceptors.ConnectAuthHandler(signingToken),
		userServerClient: do.MustInvoke[userServerConnect.UserServerClient](nil),
	}
}
