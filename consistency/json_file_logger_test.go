package consistency_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/consistency"
)

func decodeJson[K, V any](logFile string) (events []consistency.Event[K, V], err error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		var event consistency.Event[K, V]
		if err := decoder.Decode(&event); err != nil {
			break
		}
		events = append(events, event)
	}
	return events, nil
}

func TestJsonFileLogger_Succeeds(t *testing.T) {
	logFile := "json_file_logger_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	time.Sleep(200 * time.Millisecond)
	events, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "events length must be 4", len(events), 4)
	for i := range 4 {
		assert.That(t, "sequence must be correct", events[i].Sequence, uint64(i+1))
	}
}

func TestJsonFileLogger_Error_Handling(t *testing.T) {
	logFile := "/non-existent/json_file_error_handling.log"
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	time.Sleep(200 * time.Millisecond)
	_, err := decodeJson[string, string](logFile)
	assert.That(t, "err must not be nil", err != nil, true)
}

func TestJsonFileLogger_Graceful_Shutdown(t *testing.T) {
	logFile := "json_file_graceful_shutdown.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	// Close the logger gracefully
	errClose := logger.Close()
	assert.That(t, "err must be nil", errClose == nil, true)
	// Verify all events are written before shutdown
	events, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "events length must be 4", len(events), 4)
}

func TestJsonFileLogger_ReadEvents_Error(t *testing.T) {
	logFile := "/non-existent/json_file_read_events_error.log"
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()
	// Call ReadEvents to attempt to read events from the file.
	_, errorCh := logger.ReadEvents()
	// Use a select statement to capture the first error from the error channel.
	select {
	case err := <-errorCh:
		// Verify that an error is received, as the file does not exist.
		assert.That(t, "err must not be nil", err != nil, true)
	}
}

func TestJsonFileLogger_ReadEvents_Succeeds(t *testing.T) {
	logFile := "json_file_read_events_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	logger.WritePut("1", "value")
	logger.Close()
	// Call ReadEvents to read back the events from the file.
	eventCh, errorCh := logger.ReadEvents()
	select {
	case event := <-eventCh:
		// Verify that the event read back matches the one that was written.
		assert.That(t, "key must be correct", event.Key, "1")
	case err := <-errorCh:
		// Verify that no error occurred during reading.
		assert.That(t, "err must be nil", err == nil, true)
	}
}

func TestJsonFileLogger_Continue_Sequence(t *testing.T) {
	logFile := "json_file_logger_continue_sequece.log"
	defer os.Remove(logFile)

	// Create a new logger and write some events.
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	logger.Close()

	// Decode the events from the log file.
	events, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "length of events must be 4", len(events), 4)
	for i := 0; i < 4; i++ {
		assert.That(t, "sequence must be correct", events[i].Sequence, uint64(i+1))
	}

	// Restart the logger to test sequence continuity.
	logger2 := consistency.NewJsonFileLogger[string, string](logFile)
	logger2.WritePut("key4", "value4")
	logger2.Close()

	// Decode the updated events from the log file.
	updatedEvents, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil after restart", err == nil, true)
	assert.That(t, "length of events must be 5", len(updatedEvents), 5)
	assert.That(t, "last sequence must be 5", updatedEvents[4].Sequence, uint64(5))
}
