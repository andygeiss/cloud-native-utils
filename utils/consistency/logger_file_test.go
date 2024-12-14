package consistency_test

import (
	"cloud-native/utils/consistency"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func decode[K, V any](logFile string) (events []consistency.Event[K, V], err error) {
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

func TestFileLogger_Succeeds(t *testing.T) {
	logFile := "filelogger_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	time.Sleep(200 * time.Millisecond)
	events, err := decode[string, string](logFile)
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

func TestFileLogger_Error_Handling(t *testing.T) {
	logFile := "/non-existent/filelogger_error_handling.log"
	logger := consistency.NewFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")
	time.Sleep(200 * time.Millisecond)
	_, err := decode[string, string](logFile)
	if err == nil {
		t.Fatal("err must be not nil")
	}
}

func TestFileLogger_Graceful_Shutdown(t *testing.T) {
	logFile := "filelogger_graceful_shutdown.log"
	defer os.Remove(logFile)
	logger := consistency.NewFileLogger[string, string](logFile)
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	// Close the logger gracefully
	if err := logger.Close(); err != nil {
		t.Fatalf("failed to close logger: %v", err)
	}
	// Verify all events are written before shutdown
	events, err := decode[string, string](logFile)
	if err != nil {
		t.Fatalf("err must be nil, but got %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 events, but got %d", len(events))
	}
}
