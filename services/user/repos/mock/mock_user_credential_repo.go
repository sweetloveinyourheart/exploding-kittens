package userserver_mock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/planning-pocker/services/user/models"
)

type MockUserCredentialRepository struct {
	mock.Mock
}

func (m *MockUserCredentialRepository) CreateCredential(ctx context.Context, userCredential *models.UserCredential) error {
	args := m.Called(ctx, userCredential)
	return args.Error(0)
}

func (m *MockUserCredentialRepository) GetUserCredentials(ctx context.Context, userID uuid.UUID) ([]models.UserCredential, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserCredential), args.Error(1)
}
