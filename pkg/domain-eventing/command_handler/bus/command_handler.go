package bus

import (
	"sync"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// CommandHandler in the domain-eventing package is an eventing.CommandHandler that handles commands
// by routing them to the other eventing.CommandHandlers that are registered in it.
type CommandHandler struct {
	handlers   map[common.CommandType]eventing.CommandHandler
	handlersMu sync.RWMutex
}
