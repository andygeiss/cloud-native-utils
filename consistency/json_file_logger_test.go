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

func Test_JsonFileLogger_With_CloseGracefully_Should_WriteAllEvents(t *testing.T) {
	// Arrange
	logFile := "json_file_graceful_shutdown.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")

	// Act
	errClose := logger.Close()

	// Assert
	assert.That(t, "err must be nil", errClose == nil, true)
	events, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "events length must be 4", len(events), 4)
}

func Test_JsonFileLogger_With_NonExistentPath_Should_ReturnError(t *testing.T) {
	// Arrange
	logFile := "/non-existent/json_file_error_handling.log"
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()
	logger.WritePut("key1", "value1")

	// Act
	time.Sleep(200 * time.Millisecond)
	_, err := decodeJson[string, string](logFile)

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_JsonFileLogger_With_ReadEventsError_Should_ReturnError(t *testing.T) {
	// Arrange
	logFile := "/non-existent/json_file_read_events_error.log"
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()

	// Act
	_, errorCh := logger.ReadEvents()

	// Assert
	select {
	case err := <-errorCh:
		assert.That(t, "err must not be nil", err != nil, true)
	}
}

func Test_JsonFileLogger_With_ReadEventsSuccess_Should_ReturnEvents(t *testing.T) {
	// Arrange
	logFile := "json_file_read_events_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	logger.WritePut("1", "value")
	logger.Close()

	// Act
	eventCh, errorCh := logger.ReadEvents()

	// Assert
	select {
	case event := <-eventCh:
		assert.That(t, "key must be correct", event.Key, "1")
		assert.That(t, "value must be correct", event.Value, "value")
	case err := <-errorCh:
		assert.That(t, "err must be nil", err == nil, true)
	}
}

func Test_JsonFileLogger_With_WritePutAndDelete_Should_LogAllEvents(t *testing.T) {
	// Arrange
	logFile := "json_file_logger_succeeds.log"
	defer os.Remove(logFile)
	logger := consistency.NewJsonFileLogger[string, string](logFile)
	defer logger.Close()

	// Act
	logger.WritePut("key1", "value1")
	logger.WritePut("key2", "value2")
	logger.WritePut("key3", "value3")
	logger.WriteDelete("key2")
	time.Sleep(200 * time.Millisecond)

	// Assert
	events, err := decodeJson[string, string](logFile)
	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "events length must be 4", len(events), 4)
	for i := range 4 {
		assert.That(t, "sequence must be correct", events[i].Sequence, uint64(i+1))
	}
}
