package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
)

type CardDetail struct {
	models.Card
	Type models.CardType `json:"type"`
}

type ICardRepository interface {
	GetCards(ctx context.Context) ([]CardDetail, error)
}
