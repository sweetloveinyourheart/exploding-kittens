package common

// ProjectorType is the type of a projector, used as its unique identifier.
type ProjectorType string

// String returns the string representation of a projector type.
func (t ProjectorType) String() string {
	return string(t)
}

// Entity is an item which is identified by an ID.
//
// From http://cqrs.nu/Faq:
// "Entities or reference types are characterized by having an identity that's
// not tied to their attribute values. All attributes in an entity can change
// and it's still "the same" entity. Conversely, two entities might be
// equivalent in all their attributes, but will still be distinct".
type Entity interface {
	// EntityID returns the ID of the entity.
	EntityID() string
}
