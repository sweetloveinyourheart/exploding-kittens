package handlers

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"
)

type handler struct {
	ctx              context.Context
	userServerClient grpcconnect.UserServerClient
}

func NewGatewayHandler(context context.Context) *handler {
	return &handler{
		ctx:              context,
		userServerClient: do.MustInvoke[grpcconnect.UserServerClient](nil),
	}
}
