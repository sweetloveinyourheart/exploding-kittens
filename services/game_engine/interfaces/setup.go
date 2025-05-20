package interfaces

import "github.com/gofrs/uuid"

type CardSetup struct {
	StandardCards        []uuid.UUID
	ExplodingKittenCards []uuid.UUID
	DefuseCards          []uuid.UUID
}
