package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/models"
)

type CardRepository struct {
	Tx db.DbOrTx
}

func NewCardRepository(tx db.DbOrTx) ICardRepository {
	return &CardRepository{
		Tx: tx,
	}
}

func (r *CardRepository) GetCards(ctx context.Context) ([]models.Card, error) {
	var cards []models.Card
	query := `
		SELECT 
			cards.card_id, 
			cards.name, 
			cards.description, 
			cards.effect, 
			cards.quantity
		FROM cards 
		INNER JOIN card_types ON cards.type_id = card_types.type_id;
	`
	rows, err := r.Tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card models.Card
		if err := rows.Scan(
			&card.CardID,
			&card.Name,
			&card.Description,
			&card.Effect,
			&card.Quantity,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
