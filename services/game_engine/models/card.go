package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// CardType represents the different categories of cards
type CardType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Card represents an individual card with effects stored in JSONB
type Card struct {
	ID          uuid.UUID `json:"id"`
	TypeID      int       `json:"type_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Effect      []byte    `json:"effect"` // JSONB field
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
