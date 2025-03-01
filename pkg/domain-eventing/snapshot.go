package eventing

import (
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"

	"github.com/gofrs/uuid"
)

// Snapshotable is an interface for creating and applying a Snapshot record.
type Snapshotable interface {
	CreateSnapshot() *Snapshot
	ApplySnapshot(snapshot *Snapshot)
}

// Snapshot is a recording of the state of an aggregate at a point in time
type Snapshot struct {
	Version       uint64
	AggregateType common.AggregateType
	Timestamp     time.Time
	State         any
}

var snapshotDataFactories = make(map[common.AggregateType]func(uuid.UUID) SnapshotData)

type SnapshotData any

var snapshotDataFactoriesMu sync.RWMutex

var ErrSnapshotDataNotRegistered = errors.New("snapshot data not registered")

// RegisterSnapshotData registers a snapshot factory for a type. The factory is
// used to create concrete snapshot state type when unmarshalling.
//
// An example would be:
//
//	RegisterSnapshotData("aggregateType1", func() SnapshotData { return &MySnapshotData{} })
func RegisterSnapshotData[T any](aggregateType common.AggregateType) {
	if aggregateType == common.AggregateType("") {
		panic("eventing: attempt to register empty aggregate type")
	}

	snapshotDataFactoriesMu.Lock()
	defer snapshotDataFactoriesMu.Unlock()

	if _, ok := snapshotDataFactories[aggregateType]; ok {
		panic(fmt.Sprintf("eventing: registering duplicate types for %q", aggregateType))
	}

	snapshotDataFactories[aggregateType] = func(u uuid.UUID) SnapshotData {
		return new(T)
	}
}

// CreateSnapshotData create a concrete instance using the registered snapshot factories.
func CreateSnapshotData(AggregateID uuid.UUID, aggregateType common.AggregateType) (SnapshotData, error) {
	snapshotDataFactoriesMu.RLock()
	defer snapshotDataFactoriesMu.RUnlock()

	if factory, ok := snapshotDataFactories[aggregateType]; ok {
		return factory(AggregateID), nil
	}

	return nil, ErrSnapshotDataNotRegistered
}
