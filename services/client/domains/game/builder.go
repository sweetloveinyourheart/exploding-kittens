package game

import (
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
)

type GameResponseBuilder struct {
	UserID uuid.UUID
}

func NewGameResponseBuilder(UserID uuid.UUID) *GameResponseBuilder {
	return &GameResponseBuilder{
		UserID: UserID,
	}
}

func (a *GameResponseBuilder) Build(gameState *game.Game, deskState *desk.Desk, handState map[string]*hand.Hand) (*proto.Game, error) {
	if gameState == nil {
		return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "game is missing"))
	}

	players := make([]*proto.Game_Player, 0)
	playerHands := make(map[string]*proto.Game_PlayerHand)
	for _, player := range gameState.GetPlayers() {
		playerID := player.PlayerID

		playerData := &proto.Game_Player{
			PlayerId: playerID.String(),
			Active:   player.IsActive(),
		}
		players = append(players, playerData)

		playerHand := handState[playerData.GetPlayerId()]
		playerHandData := &proto.Game_PlayerHand{
			RemainingCards: int32(len(playerHand.GetCardIDs())),
		}

		// Only the authorized player can view their own card
		if playerID == a.UserID {
			playerHandData.Hands = stringsutil.ConvertUUIDsToStrings(playerHand.GetCardIDs())
		}

		playerHands[playerID.String()] = playerHandData
	}

	desk := &proto.Game_Desk{
		DeskId:         deskState.GetDeskID().String(),
		RemainingCards: int32(len(deskState.GetCardIDs())),
		DiscardPile:    stringsutil.ConvertUUIDsToStrings(deskState.GetDiscardPile()),
	}

	result := &proto.Game{
		GameId:      gameState.GetGameID().String(),
		GamePhase:   proto.Game_Phase(gameState.GamePhase),
		PlayerTurn:  gameState.PlayerTurn.String(),
		Players:     players,
		PlayerHands: playerHands,
		Desk:        desk,
	}

	return result, nil
}
