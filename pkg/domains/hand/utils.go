package hand

import (
	"bytes"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
)

func NewPlayerHandID(gameID uuid.UUID, playerID uuid.UUID) uuid.UUID {
	buf := bytesBufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bytesBufferPool.Put(buf)
	}()

	buf.Write(gameID.Bytes())
	buf.Write(playerID.Bytes())

	return uuid.NewV5(constants.NameSpaceGames, buf.String())
}

var bytesBufferPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}
