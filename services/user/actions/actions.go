package actions

import (
	"context"

	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/services/user/repos"
)

type actions struct {
	context context.Context

	userRepo           repos.IUserRepository
	userCredentialRepo repos.IUserCredentialRepository
	userSessionRepo    repos.IUserSessionRepository
}

func NewActions(ctx context.Context, signingToken string) *actions {
	userRepo := do.MustInvoke[repos.IUserRepository](nil)
	userCredentialRepo := do.MustInvoke[repos.IUserCredentialRepository](nil)
	userSessionRepo := do.MustInvoke[repos.IUserSessionRepository](nil)

	return &actions{
		context:            ctx,
		userRepo:           userRepo,
		userCredentialRepo: userCredentialRepo,
		userSessionRepo:    userSessionRepo,
	}
}
