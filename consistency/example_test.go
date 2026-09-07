package consistency_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/andygeiss/cloud-native-utils/consistency"
)

func ExampleNewJsonFileLogger() {
	dir, err := os.MkdirTemp("", "consistency")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	path := filepath.Join(dir, "events.log")

	// Every change is appended to the log, in order.
	logger := consistency.NewJsonFileLogger[string, string](path)
	logger.WritePut("user-1", "Alice")
	logger.WritePut("user-2", "Bob")
	logger.WriteDelete("user-1")

	// Close flushes what is still pending.
	if err := logger.Close(); err != nil {
		log.Fatal(err)
	}

	// Replaying the log rebuilds the state it describes.
	replay := consistency.NewJsonFileLogger[string, string](path)
	defer func() { _ = replay.Close() }()

	events, errs := replay.ReadEvents()
	for event := range events {
		fmt.Println(event.Sequence, event.Key, event.Value)
	}
	if err := <-errs; err != nil {
		log.Fatal(err)
	}
	// Output:
	// 1 user-1 Alice
	// 2 user-2 Bob
	// 3 user-1
}
