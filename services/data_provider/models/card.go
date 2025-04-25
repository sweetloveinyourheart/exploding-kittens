package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Card represents an individual card with effects stored in JSONB
type Card struct {
	CardID      uuid.UUID `json:"card_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CardEffect struct {
	EffectID  uuid.UUID `json:"effect_id"`
	CardID    uuid.UUID `json:"card_id"`
	Effect    []byte    `json:"effect"` // JSONB is represented as a byte slice in Go
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ComboEffect struct {
	ComboID       uuid.UUID `json:"combo_id"`
	RequiredCards int       `json:"required_cards"`
	Effect        []byte    `json:"effect"` // JSONB is represented as a byte slice in Go
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
