package actions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"

	"github.com/sweetloveinyourheart/planning-pocker/services/user/repos"
)

type actions struct {
	context context.Context

	userRepo           repos.IUserRepository
	userCredentialRepo repos.IUserCredentialRepository
	userSessionRepo    repos.IUserSessionRepository
}

func NewActions(ctx context.Context, signingToken string) *actions {
	dbConn := do.MustInvoke[*pgxpool.Pool](nil)

	userRepo := repos.NewUserRepository(dbConn)
	userCredentialRepo := repos.NewUserCredentialRepository(dbConn)
	userSessionRepo := repos.NewUserSessionRepository(dbConn)

	return &actions{
		context:            ctx,
		userRepo:           userRepo,
		userCredentialRepo: userCredentialRepo,
		userSessionRepo:    userSessionRepo,
	}
}
