package desk

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type Desk struct {
	DeskID    uuid.UUID   `json:"desk_id"`
	Cards     []uuid.UUID `json:"cards"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

var _ = common.Entity(&Desk{})

func (d *Desk) EntityID() string {
	return d.DeskID.String()
}

func (d *Desk) GetDeskID() uuid.UUID {
	return d.DeskID
}

func (d *Desk) GetCards() []uuid.UUID {
	return d.Cards
}

func (t *Desk) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Desk) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
