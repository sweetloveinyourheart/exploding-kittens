package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Card represents an individual card with effects stored in JSONB
type Card struct {
	CardID      uuid.UUID `json:"card_id"`
	TypeID      int       `json:"type_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Effect      []byte    `json:"effect"` // JSONB field
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
