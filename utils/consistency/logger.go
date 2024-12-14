package consistency

// Logger is an interface that defines the operations for a transactional log.
type Logger[K, V any] interface {
	// Close closes the logger and ensures all pending events are processed.
	Close() error
	// WriteDelete writes a delete event to the log.
	WriteDelete(key K)
	// WritePut writes a put event to the log.
	WritePut(key K, value V)
}
