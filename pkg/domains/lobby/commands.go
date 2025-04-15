package lobby

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
)

func init() {
	eventing.RegisterCommand[CreateLobby, *CreateLobby]()
	eventing.RegisterCommand[JoinLobby, *JoinLobby]()
	eventing.RegisterCommand[LeaveLobby, *LeaveLobby]()
	eventing.RegisterCommand[StartGame, *StartGame]()
}

const (
	CreateLobbyCommand = common.CommandType("lobby:create")
	JoinLobbyCommand   = common.CommandType("lobby:join")
	LeaveLobbyCommand  = common.CommandType("lobby:leave")
	StartGameCommand   = common.CommandType("lobby:game:start")
)

var AllCommands = []common.CommandType{
	CreateLobbyCommand,
	JoinLobbyCommand,
	LeaveLobbyCommand,
	StartGameCommand,
}

// Static type check that the eventing.Command interface is implemented.
var _ = eventing.Command(&CreateLobby{})
var _ = eventing.Command(&JoinLobby{})
var _ = eventing.Command(&LeaveLobby{})
var _ = eventing.Command(&StartGame{})

type CreateLobby struct {
	LobbyID    uuid.UUID `json:"lobby_id"`
	LobbyCode  string    `json:"lobby_code"`
	LobbyName  string    `json:"lobby_name"`
	HostUserID uuid.UUID `json:"host_user_id"`
}

func (c *CreateLobby) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateLobby) AggregateID() string { return c.LobbyID.String() }

func (c *CreateLobby) CommandType() common.CommandType { return CreateLobbyCommand }

func (c *CreateLobby) Validate() error {
	if c.LobbyID == uuid.Nil {
		return &common.CommandFieldError{Field: "lobby_id", Details: "empty field"}
	}

	if stringsutil.IsBlank(c.LobbyName) {
		return &common.CommandFieldError{Field: "lobby_name", Details: "empty field"}
	}

	if stringsutil.IsBlank(c.LobbyCode) {
		return &common.CommandFieldError{Field: "lobby_code", Details: "empty field"}
	}

	if c.HostUserID == uuid.Nil {
		return &common.CommandFieldError{Field: "host_user_id", Details: "empty field"}
	}

	return nil
}

type JoinLobby struct {
	LobbyID uuid.UUID `json:"lobby_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (c *JoinLobby) AggregateType() common.AggregateType { return AggregateType }

func (c *JoinLobby) AggregateID() string { return c.LobbyID.String() }

func (c *JoinLobby) CommandType() common.CommandType { return JoinLobbyCommand }

func (c *JoinLobby) Validate() error {
	if c.LobbyID == uuid.Nil {
		return &common.CommandFieldError{Field: "lobby_id", Details: "empty field"}
	}

	if c.UserID == uuid.Nil {
		return &common.CommandFieldError{Field: "user_id", Details: "empty field"}
	}

	return nil
}

type LeaveLobby struct {
	LobbyID uuid.UUID `json:"lobby_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (c *LeaveLobby) AggregateType() common.AggregateType { return AggregateType }

func (c *LeaveLobby) AggregateID() string { return c.LobbyID.String() }

func (c *LeaveLobby) CommandType() common.CommandType { return LeaveLobbyCommand }

func (c *LeaveLobby) Validate() error {
	if c.LobbyID == uuid.Nil {
		return &common.CommandFieldError{Field: "lobby_id", Details: "empty field"}
	}

	if c.UserID == uuid.Nil {
		return &common.CommandFieldError{Field: "user_id", Details: "empty field"}
	}

	return nil
}

type StartGame struct {
	LobbyID    uuid.UUID `json:"lobby_id"`
	HostUserID uuid.UUID `json:"host_user_id"`
	GameID     uuid.UUID `json:"game_id"`
}

func (c *StartGame) AggregateType() common.AggregateType { return AggregateType }

func (c *StartGame) AggregateID() string { return c.LobbyID.String() }

func (c *StartGame) CommandType() common.CommandType { return StartGameCommand }

func (c *StartGame) Validate() error {
	if c.LobbyID == uuid.Nil {
		return &common.CommandFieldError{Field: "lobby_id", Details: "empty field"}
	}

	if c.HostUserID == uuid.Nil {
		return &common.CommandFieldError{Field: "host_user_id", Details: "empty field"}
	}

	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	return nil
}
