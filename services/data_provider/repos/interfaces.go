package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/models"
)

type CardDetail struct {
	models.Card
	Effects      []byte `json:"effects"`
	ComboEffects []byte `json:"combo_effects"`
}

type ICardRepository interface {
	GetCardsInformation(ctx context.Context) ([]models.Card, error)
	GetCards(ctx context.Context) ([]CardDetail, error)
}
