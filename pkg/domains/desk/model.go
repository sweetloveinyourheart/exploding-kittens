package desk

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type Desk struct {
	DeskID      uuid.UUID   `json:"desk_id"`
	CardIDs     []uuid.UUID `json:"card_ids"`
	DiscardPile []uuid.UUID `json:"discard_pile"`
	ShuffledAt  time.Time   `json:"shuffled_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

var _ = common.Entity(&Desk{})

func (d *Desk) EntityID() string {
	return d.DeskID.String()
}

func (d *Desk) GetDeskID() uuid.UUID {
	return d.DeskID
}

func (d *Desk) GetCardIDs() []uuid.UUID {
	return d.CardIDs
}

func (d *Desk) GetDiscardPile() []uuid.UUID {
	return d.DiscardPile
}

func (d *Desk) GetShuffledAt() time.Time {
	return d.ShuffledAt
}

func (d *Desk) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *Desk) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}
