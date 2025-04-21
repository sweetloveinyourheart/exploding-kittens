package game

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

const (
	GAME_PHASE_INITIALIZING = "INITIALZING"
	GAME_PHASE_CARD_PLAYING = "CARD_PLAYING"
	GAME_PHASE_CARD_DRAWING = "CARD_DRAWING"
)

type Game struct {
	GameID      uuid.UUID               `json:"game_id"`
	GamePhase   string                  `json:"game_phase"`
	Desk        uuid.UUID               `json:"desk"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
	PlayerTurn  uuid.UUID               `json:"player_turn"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

var _ = common.Entity(&Game{})

func (d *Game) EntityID() string {
	return d.GameID.String()
}

func (d *Game) GetGameID() uuid.UUID {
	return d.GameID
}

func (d *Game) GetDesk() uuid.UUID {
	return d.Desk
}

func (d *Game) GetPlayers() map[uuid.UUID]uuid.UUID {
	return d.PlayerHands
}

func (d *Game) GetGamePhase() string {
	return d.GamePhase
}

func (d *Game) GetPlayerTurn() uuid.UUID {
	return d.PlayerTurn
}

func (t *Game) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Game) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
