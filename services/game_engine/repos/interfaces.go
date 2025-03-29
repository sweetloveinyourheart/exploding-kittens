package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
)

type ICardRepository interface {
	GetCards(ctx context.Context) ([]models.Card, error)
}
