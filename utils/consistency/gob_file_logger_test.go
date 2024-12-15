package consistency_test

import (
	"cloud-native/utils/consistency"
	"encoding/gob"
	"os"
	"testing"
	"time"
)

func decodeGob[K, V any](logFile string) (events []consistency.Event[K, V], err error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	for {
		var event consistency.Event[K, V]
		if err := decoder.Decode(&event); err != nil {
			break
		}
		events = append(events, event)
	}
	return events, nil
}

func TestGobFileLogger_Succeeds(t *testing.T) {
	logFile := "gob_file_logger_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewGobFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	time.Sleep(200 * time.Millisecond)
	events, err := decodeGob[string, string](logFile)
	if err != nil {
		t.Fatalf("err must be nil, but got %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 events, but got %d", len(events))
	}
	for i := range 4 {
		if events[i].Sequence != uint64(i+1) {
			t.Fatalf("sequence must be %d, but got %d", events[0].Sequence, uint64(i+1))
		}
	}
}

func TestGobFileLogger_Error_Handling(t *testing.T) {
	logFile := "/non-existent/gob_file_logger_error_handling.log"
	logger := consistency.NewGobFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	time.Sleep(200 * time.Millisecond)
	_, err := decodeGob[string, string](logFile)
	if err == nil {
		t.Fatal("err must be not nil")
	}
}

func TestGobFileLogger_Graceful_Shutdown(t *testing.T) {
	logFile := "gob_file_logger_graceful_shutdown.log"
	defer os.Remove(logFile)
	logger := consistency.NewGobFileLogger[string, string](logFile)
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	// Close the logger gracefully
	if err := logger.Close(); err != nil {
		t.Fatalf("failed to close logger: %v", err)
	}
	// Verify all events are written before shutdown
	events, err := decodeGob[string, string](logFile)
	if err != nil {
		t.Fatalf("err must be nil, but got %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 events, but got %d", len(events))
	}
}

func TestGobFileLogger_ReadEvents_Error(t *testing.T) {
	logFile := "/non-existent/gob_file_logger_read_events_error.log"
	logger := consistency.NewGobFileLogger[string, string](logFile)
	defer logger.Close()
	// Call ReadEvents to attempt to read events from the file.
	_, errorCh := logger.ReadEvents()
	// Use a select statement to capture the first error from the error channel.
	select {
	case err := <-errorCh:
		// Verify that an error is received, as the file does not exist.
		if err == nil {
			t.Fatal("err must not be nil")
		}
	}
}

func TestGobFileLogger_ReadEvents_Succeeds(t *testing.T) {
	logFile := "gob_file_logger_read_events_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewGobFileLogger[string, string](logFile)
	logger.WritePut("1", "value")
	logger.Close()
	// Call ReadEvents to read back the events from the file.
	eventCh, errorCh := logger.ReadEvents()
	select {
	case event := <-eventCh:
		// Verify that the event read back matches the one that was written.
		if event.Key != "1" {
			t.Fatalf("key must be correct, but got %v", event.Key)
		}
	case err := <-errorCh:
		// Verify that no error occurred during reading.
		if err != nil {
			t.Fatal("err must be nil") // Fail the test if an error is returned.
		}
	}
}
