package tracing

// MapMessageCarrier injects and extracts traces from a map.
type MapMessageCarrier struct {
	m map[string]any
}

// NewMapMessageCarrier creates a new MapMessageCarrier.
func NewMapMessageCarrier(m map[string]any) MapMessageCarrier {
	return MapMessageCarrier{m: m}
}

// Get retrieves a single value for a given key.
func (c MapMessageCarrier) Get(key string) string {
	value, ok := c.m[key]
	if ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// Set sets a header.
func (c MapMessageCarrier) Set(key, val string) {
	c.m[key] = val
}

// Keys returns a slice of all key identifiers in the carrier.
func (c MapMessageCarrier) Keys() []string {
	out := make([]string, 0, len(c.m))
	for key := range c.m {
		out = append(out, key)
	}
	return out
}
