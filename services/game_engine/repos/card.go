package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
)

type CardRepository struct {
	Tx db.DbOrTx
}

func NewCardRepository(tx db.DbOrTx) ICardRepository {
	return &CardRepository{
		Tx: tx,
	}
}

func (r *CardRepository) GetCards(ctx context.Context) ([]CardDetail, error) {
	var cards []CardDetail
	query := `
		SELECT 
			cards.id, 
			cards.name, 
			cards.description, 
			cards.effect, 
			cards.quantity,
			cart_types.name,
			cart_types.description
		FROM cards 
		INNER JOIN card_types ON cards.type_id = card_types.id;
	`
	rows, err := r.Tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card CardDetail
		if err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Description,
			&card.Effect,
			&card.Quantity,
			&card.Type.Name,
			&card.Type.Description,
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
