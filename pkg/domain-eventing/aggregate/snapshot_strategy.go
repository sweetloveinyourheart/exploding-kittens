package aggregate

import (
	"time"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// NoSnapshotStrategy no snapshot should be taken.
type NoSnapshotStrategy struct {
}

func (s *NoSnapshotStrategy) ShouldTakeSnapshot(_ uint64, _ time.Time, _ common.Event) bool {
	return false
}

// EveryNumberEventSnapshotStrategy use to take a snapshot every n number of events.
type EveryNumberEventSnapshotStrategy struct {
	snapshotThreshold uint64
}

func NewEveryNumberEventSnapshotStrategy(threshold uint64) *EveryNumberEventSnapshotStrategy {
	return &EveryNumberEventSnapshotStrategy{
		snapshotThreshold: threshold,
	}
}

func (s *EveryNumberEventSnapshotStrategy) ShouldTakeSnapshot(lastSnapshotVersion uint64, _ time.Time, event common.Event) bool {
	return event.Version()-lastSnapshotVersion >= s.snapshotThreshold
}

// PeriodSnapshotStrategy use to take a snapshot every time a period has elapsed, for example every hour.
type PeriodSnapshotStrategy struct {
	snapshotThreshold time.Duration
}

func NewPeriodSnapshotStrategy(threshold time.Duration) *PeriodSnapshotStrategy {
	return &PeriodSnapshotStrategy{
		snapshotThreshold: threshold,
	}
}

func (s *PeriodSnapshotStrategy) ShouldTakeSnapshot(_ uint64,
	lastSnapshotTimestamp time.Time,
	event common.Event) bool {
	return event.Timestamp().Sub(lastSnapshotTimestamp) >= s.snapshotThreshold
}
