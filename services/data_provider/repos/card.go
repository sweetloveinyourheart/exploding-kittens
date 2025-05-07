package repos

import (
	"context"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/models"
)

type CardRepository struct {
	Tx db.DbOrTx
}

func NewCardRepository(tx db.DbOrTx) ICardRepository {
	return &CardRepository{
		Tx: tx,
	}
}

func (r *CardRepository) GetCardsInformation(ctx context.Context) ([]models.Card, error) {
	var cards []models.Card
	query := `
		SELECT 
			cards.card_id, 
			cards.name, 
			cards.description,
			cards.quantity
		FROM cards;
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

func (r *CardRepository) GetCards(ctx context.Context) ([]CardDetail, error) {
	var cards []CardDetail
	query := `
		SELECT
			c.card_id,
			c.name AS card_name,
			c.description AS card_description,
			c.quantity,
			COALESCE(ce.effect, '{}') AS card_effect, -- Defaults to empty JSON if no effect
			COALESCE(
				jsonb_agg(DISTINCT jsonb_build_object('type', ce2.effect ->> 'type', 'required_cards', ce2.required_cards)) FILTER (
					WHERE
						ce2.effect IS NOT NULL
				),
				'[]'
			) AS combo_effects
		FROM
			cards c
			LEFT JOIN card_effects ce ON c.card_id = ce.card_id
			LEFT JOIN card_combo cc ON c.card_id = cc.card_id
			LEFT JOIN combo_effects ce2 ON cc.combo_id = ce2.combo_id
		GROUP BY
			c.card_id,
			c.name,
			c.description,
			c.quantity,
			ce.effect;
	`
	rows, err := r.Tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card CardDetail
		if err := rows.Scan(
			&card.CardID,
			&card.Name,
			&card.Code,
			&card.Description,
			&card.Quantity,
			&card.Effects,
			&card.ComboEffects,
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
