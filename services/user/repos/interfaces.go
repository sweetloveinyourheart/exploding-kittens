package repos

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/planning-pocker/services/user/models"
)

type IUserRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, bool, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUserData(ctx context.Context, user *models.User) error
}

type IUserCredentialRepository interface {
	CreateCredential(ctx context.Context, userCredential *models.UserCredential) error
	GetUserCredentials(ctx context.Context, userID uuid.UUID) ([]models.UserCredential, error)
}

type IUserSessionRepository interface {
	CreateSession(ctx context.Context, userSession *models.UserSession) error
	GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error)
}
