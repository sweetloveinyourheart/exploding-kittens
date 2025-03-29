package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
)

type CardDetail struct {
	models.Card
	Effects      []byte
	ComboEffects []byte
}

type ICardRepository interface {
	GetCardsInformation(ctx context.Context) ([]models.Card, error)
	GetCards(ctx context.Context) ([]CardDetail, error)
}
