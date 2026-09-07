package consistency_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andygeiss/cloud-native-utils/consistency"
)

// FuzzReadEvents replays a log file that may have been truncated mid-write or
// corrupted on disk. Reading it must report an error, never panic or hang.
func FuzzReadEvents(f *testing.F) {
	f.Add(`{"key":"1","value":"a","sequence":1,"event_type":1}` + "\n")
	f.Add(`{"key":"1","value":"a","sequence":1,"event_type":1}` + "\n" + `{"key":"1","sequence":2,"event_type":0}` + "\n")
	f.Add(`{"key":"1","value":"a","sequence":1,"event_ty`) // truncated mid-write
	f.Add("not json\n")
	f.Add("{}")
	f.Add("")

	f.Fuzz(func(t *testing.T, content string) {
		path := filepath.Join(t.TempDir(), "events.log")
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("writing the log failed: %v", err)
		}

		logger := consistency.NewJsonFileLogger[string, string](path)
		defer func() { _ = logger.Close() }()

		events, errs := logger.ReadEvents()
		for range events {
			// Draining is the point; the decoded events do not matter here.
		}
		<-errs
	})
}
