package actions

import (
	"context"

	"github.com/samber/do"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/repos"
)

type actions struct {
	context     context.Context
	defaultAuth func(context.Context, string) (context.Context, error)

	userRepo           repos.IUserRepository
	userCredentialRepo repos.IUserCredentialRepository
	userSessionRepo    repos.IUserSessionRepository

	tracer trace.Tracer
}

// AuthFuncOverride is a callback function that overrides the default authorization middleware in the GRPC layer. This is
// used to allow unauthenticated endpoints (such as login) to be called without a token.
func (a *actions) AuthFuncOverride(ctx context.Context, token string, fullMethodName string) (context.Context, error) {
	if fullMethodName == grpcconnect.UserServerCreateNewUserProcedure {
		return ctx, nil
	}

	if fullMethodName == grpcconnect.UserServerSignInProcedure {
		return ctx, nil
	}

	return a.defaultAuth(ctx, token)
}

func NewActions(ctx context.Context, signingToken string) *actions {
	userRepo := do.MustInvoke[repos.IUserRepository](nil)
	userCredentialRepo := do.MustInvoke[repos.IUserCredentialRepository](nil)
	userSessionRepo := do.MustInvoke[repos.IUserSessionRepository](nil)

	tracer := otel.Tracer("com.sweetloveinyourheart.kittens.userserver.actions")

	return &actions{
		context:            ctx,
		defaultAuth:        interceptors.ConnectServerAuthHandler(signingToken),
		userRepo:           userRepo,
		userCredentialRepo: userCredentialRepo,
		userSessionRepo:    userSessionRepo,
		tracer:             tracer,
	}
}
