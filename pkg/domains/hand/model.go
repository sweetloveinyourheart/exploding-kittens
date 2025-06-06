package hand

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type Hand struct {
	HandID     uuid.UUID   `json:"hand_id"`
	CardIDs    []uuid.UUID `json:"card_ids"`
	ShuffledAt time.Time   `json:"shuffled_at"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

var _ = common.Entity(&Hand{})

func (d *Hand) EntityID() string {
	return d.HandID.String()
}

func (d *Hand) GetHandID() uuid.UUID {
	return d.HandID
}

func (d *Hand) GetCardIDs() []uuid.UUID {
	return d.CardIDs
}

func (t *Hand) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Hand) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
