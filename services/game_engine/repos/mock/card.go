package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/repos"
)

type MockCardRepository struct {
	mock.Mock
}

func (m *MockCardRepository) GetCardsInformation(ctx context.Context) ([]models.Card, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Card), args.Error(1)
}

func (m *MockCardRepository) GetCards(ctx context.Context) ([]repos.CardDetail, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repos.CardDetail), args.Error(1)
}
