package consistency

import (
	"encoding/json"
	"os"
	"sync"
)

// FileLogger is a file-based implementation of the Logger interface.
// It writes events to a JSON-formatted file for persistence.
type FileLogger[K, V any] struct {
	errorCh      chan error       // Channel for propagating errors to the caller.
	eventCh      chan Event[K, V] // Channel for queuing events to be written.
	file         string           // Path to the log file.
	lastSequence uint64           // Sequence number of the last event.
	mutex        sync.Mutex       // Mutex to protect shared resources.
	wg           sync.WaitGroup   // WaitGroup for ensuring graceful shutdown.
	closeOnce    sync.Once        // Ensures the Close method is called only once.
}

// NewFileLogger initializes a new FileLogger for the given file path.
func NewFileLogger[K, V any](file string) *FileLogger[K, V] {
	errorCh := make(chan error, 1)
	eventCh := make(chan Event[K, V], 100) // Buffered channel for queuing events.
	logger := &FileLogger[K, V]{
		errorCh: errorCh,
		eventCh: eventCh,
		file:    file,
	}
	// Start the event processing goroutine.
	logger.wg.Add(1)
	go logger.run()
	return logger
}

// run processes events from the event channel and writes them to the file.
func (a *FileLogger[K, V]) run() {
	// Mark the goroutine as done when this method exits.
	defer a.wg.Done()
	// Open the log file for appending or create it if it doesn't exist.
	file, err := os.OpenFile(a.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		a.errorCh <- err // Report the error if the file can't be opened.
		return
	}
	defer file.Close()
	// JSON encoder for writing events to the file.
	encoder := json.NewEncoder(file)
	// Process events from the event channel.
	for event := range a.eventCh {
		if err := encoder.Encode(event); err != nil {
			a.errorCh <- err
			return
		}
	}
}

// Close shuts down the logger, ensuring all pending events are written.
func (a *FileLogger[K, V]) Close() error {
	var closeErr error
	// Ensure Close is executed only once.
	a.closeOnce.Do(func() {
		close(a.eventCh) // Signal the event processing loop to stop.
		a.wg.Wait()      // Wait for the processing goroutine to finish.
		// Close the error channel and capture any errors that occurred.
		close(a.errorCh)
		for err := range a.errorCh {
			closeErr = err
		}
	})
	return closeErr
}

// Error returns a read-only channel for retrieving errors.
func (a *FileLogger[K, V]) Error() <-chan error {
	return a.errorCh
}

// ReadEvents reads events from the log file and returns two channels.
// The method uses a goroutine to read events asynchronously, allowing the caller
// to process events and handle errors as they are received.
func (a *FileLogger[K, V]) ReadEvents() (<-chan Event[K, V], <-chan error) {
	errorCh := make(chan error, 1)
	eventCh := make(chan Event[K, V], 100)
	// Launch a goroutine to handle the file reading process asynchronously.
	go func() {
		defer close(errorCh)
		defer close(eventCh)
		// Open the log file for reading.
		file, err := os.Open(a.file)
		if err != nil {
			errorCh <- err
			return
		}
		defer file.Close()
		// Create a JSON decoder to read events from the file.
		decoder := json.NewDecoder(file)
		// Read events in a loop until EOF or an error occurs.
		for {
			var event Event[K, V]
			// Decode the next event from the file.
			if err := decoder.Decode(&event); err != nil {
				if err.Error() == "EOF" {
					// Exit gracefully if all events have been read.
					return
				}
				// Report any other decoding errors and terminate the loop.
				errorCh <- err
				return
			}
			// Send the successfully decoded event to the event channel.
			eventCh <- event
		}
	}()
	return eventCh, errorCh
}

// WriteDelete writes a delete event to the log.
func (a *FileLogger[K, V]) WriteDelete(key K) {
	a.mutex.Lock()         // Lock the logger to ensure thread-safe access.
	defer a.mutex.Unlock() // Unlock the logger when the method exits.
	a.lastSequence++       // Increment the sequence number for this event.
	a.eventCh <- Event[K, V]{
		Sequence:  a.lastSequence,
		EventType: EventTypeDelete,
		Key:       key,
	}
}

// WritePut writes a put event to the log.
func (a *FileLogger[K, V]) WritePut(key K, value V) {
	a.mutex.Lock()         // Lock the logger to ensure thread-safe access.
	defer a.mutex.Unlock() // Unlock the logger when the method exits.
	a.lastSequence++       // Increment the sequence number for this event.
	a.eventCh <- Event[K, V]{
		Sequence:  a.lastSequence,
		EventType: EventTypePut,
		Key:       key,
		Value:     value,
	}
}
