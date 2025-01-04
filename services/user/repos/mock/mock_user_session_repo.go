package userserver_mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/exploding-kittens/services/user/models"
)

type MockUserSessionRepository struct {
	mock.Mock
}

func (m *MockUserSessionRepository) CreateSession(ctx context.Context, userSession *models.UserSession) error {
	args := m.Called(ctx, userSession)
	return args.Error(0)
}

func (m *MockUserSessionRepository) GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(models.UserSession), args.Error(1)
}
