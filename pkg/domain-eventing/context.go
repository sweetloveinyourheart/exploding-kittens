package eventing

type contextKey int

const (
	aggregateIDKey contextKey = iota
	aggregateTypeKey
	commandTypeKey
	namespaceDetchedTypeKey
	eventSubjectKey
)
