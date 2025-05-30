package game

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

const (
	GAME_PHASE_INITIALIZING      = iota // Setting up players, shuffling and dealing cards, inserting Exploding Kittens and Defuse cards into the deck
	GAME_PHASE_TURN_START        = 1    // Active player begins their turn, player can play as many action cards as they want
	GAME_PHASE_ACTION_PHASE      = 2    // Execute the action of the played card
	GAME_PHASE_CARD_DRAWING      = 3    // Player draws one card from the deck (mandatory if they didn't Skip/Attack)
	GAME_PHASE_TURN_FINISH       = 4    // Finalize the turn, next player becomes active
	GAME_PHASE_GAME_FINISH       = 5    // When only one player remains
	GAME_PHASE_EXPLODING_DRAWN   = 6    // When a player draws an Exploding Kitten card
	GAME_PHASE_EXPLODING_DEFUSED = 7    // When a player defuses an Exploding Kitten card
	GAME_PHASE_PLAYER_ELIMINATED = 8    // When a player is eliminated from the game
)

type Game struct {
	GameID      uuid.UUID               `json:"game_id"`
	GamePhase   int                     `json:"game_phase"`
	DeskID      uuid.UUID               `json:"desk_id"`
	Players     []Player                `json:"players"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
	WinnerID    uuid.UUID               `json:"winner_id"`

	PlayerTurn      uuid.UUID `json:"player_turn"`      // The player whose turn it is
	ExecutingAction string    `json:"executing_action"` // The action that is currently being executed
	AffectedPlayer  uuid.UUID `json:"affected_player"`  // The player who is affected by the action

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var _ = common.Entity(&Game{})

func (d *Game) EntityID() string {
	return d.GameID.String()
}

func (d *Game) GetGameID() uuid.UUID {
	return d.GameID
}

func (d *Game) GetDeskID() uuid.UUID {
	return d.DeskID
}

func (d *Game) GetPlayers() []Player {
	return d.Players
}

func (d *Game) GetPlayerHands() map[uuid.UUID]uuid.UUID {
	return d.PlayerHands
}

func (d *Game) GetGamePhase() int {
	return d.GamePhase
}

func (d *Game) GetPlayerTurn() uuid.UUID {
	return d.PlayerTurn
}

func (d *Game) GetExecutingAction() string {
	return d.ExecutingAction
}

func (d *Game) GetAffectedPlayer() uuid.UUID {
	return d.AffectedPlayer
}

func (d *Game) GetWinnerID() uuid.UUID {
	return d.WinnerID
}

func (d *Game) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *Game) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}

type Player struct {
	PlayerID uuid.UUID `json:"player_id"`
	Active   bool      `json:"active"`
}

func (p *Player) GetPlayerID() uuid.UUID {
	return p.PlayerID
}

func (p *Player) IsActive() bool {
	return p.Active
}
