package consistency

// EventType represents the type of an event in the transactional log.
type EventType byte

const (
	// EventTypeDelete indicates a delete operation.
	EventTypeDelete EventType = iota
	// EventTypePut indicates a put (write) operation.
	EventTypePut
)

// Event represents an entry in the transactional log.
type Event[K, V any] struct {
	Sequence  uint64    `json:"sequence"`   // The sequence number of the event, ensuring order.
	EventType EventType `json:"event_type"` // The type of event (e.g., Put or Delete).
	Key       K         `json:"key"`        // The key associated with the event.
	Value     V         `json:"value"`      // The value associated with the event (only for Put).
}
