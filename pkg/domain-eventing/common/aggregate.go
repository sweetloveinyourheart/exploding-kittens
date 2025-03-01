package common

// AggregateType gets used to match events with aggregates.
type AggregateType string

// String returns the string representation of an aggregate type.
func (at AggregateType) String() string {
	return string(at)
}
