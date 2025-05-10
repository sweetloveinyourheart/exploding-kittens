package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
)

type actions struct {
	context      context.Context
	defaultAuth  func(context.Context, string) (context.Context, error)
	dataProvider dataProviderGrpc.DataProviderClient
}

func NewActions(ctx context.Context, signingToken string) *actions {
	return &actions{
		context:      ctx,
		defaultAuth:  interceptors.ConnectServerAuthHandler(signingToken),
		dataProvider: do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
	}
}
