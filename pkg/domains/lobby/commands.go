package lobby

import (
	"github.com/gofrs/uuid"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
)

func init() {
	eventing.RegisterCommand[CreateLobby, *CreateLobby]()
}

const (
	CreateLobbyCommand = common.CommandType("lobby:create")
)

var AllCommands = []common.CommandType{
	CreateLobbyCommand,
}

// Static type check that the eventing.Command interface is implemented.
var _ = eventing.Command(&CreateLobby{})

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

	if stringsutil.IsBlank(c.LobbyCode) {
		return &common.CommandFieldError{Field: "lobby_code", Details: "empty field"}
	}

	if c.HostUserID == uuid.Nil {
		return &common.CommandFieldError{Field: "host_user_id", Details: "empty field"}
	}

	return nil
}
