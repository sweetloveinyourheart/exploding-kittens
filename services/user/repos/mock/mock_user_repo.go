package userserver_mock

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/planning-pocker/services/user/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, bool, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.User), args.Bool(1), args.Error(2)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (models.User, bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Bool(1), args.Error(2)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserData(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
