package extensibility_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/extensibility"
)

type TestPort interface {
	FindByID(id string) (name string, err error)
}

// testPluginPath is the compiled testdata/adapter.go, or empty when
// testPluginErr says why it is not.
var testPluginPath string

// testPluginErr is the reason the tests needing a plugin have to skip.
var testPluginErr error

// TestMain builds the plugin once. Go keys a plugin by its package path, so a
// second build of the same source fails to open with "plugin already loaded".
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "extensibility-plugin")
	if err != nil {
		testPluginErr = err
	} else {
		testPluginPath, testPluginErr = buildTestPlugin(dir)
	}
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

// buildTestPlugin compiles testdata/adapter.go into dir. A plugin must match
// the loading binary's toolchain and race setting exactly, so a committed .so
// would go stale on every Go release and could never satisfy both "go test"
// and "go test -race" at once.
func buildTestPlugin(dir string) (path string, err error) {
	args := []string{"build", "-buildmode=plugin"}
	if isRaceEnabled() {
		args = append(args, "-race")
	}
	if runtime.GOOS == "darwin" {
		// dyld rejects Go's chained fixups in a plugin, so plugin.Open fails
		// with "seg_count does not match number of segments" without this.
		args = append(args, "-ldflags=-extldflags=-Wl,-no_fixup_chains")
	}
	path = filepath.Join(dir, "adapter.so")
	args = append(args, "-o", path, "testdata/adapter.go")
	if out, err := exec.Command("go", args...).CombinedOutput(); err != nil {
		return "", fmt.Errorf("go %s: %w: %s", strings.Join(args, " "), err, out)
	}
	return path, nil
}

// isRaceEnabled reports whether this test binary was built with -race.
func isRaceEnabled() bool {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return false
	}
	for _, setting := range info.Settings {
		if setting.Key == "-race" {
			return setting.Value == "true"
		}
	}
	return false
}

// requireTestPlugin returns the plugin path, skipping where Go cannot build or
// load one.
func requireTestPlugin(t *testing.T) string {
	t.Helper()
	if runtime.GOOS != "darwin" && runtime.GOOS != "freebsd" && runtime.GOOS != "linux" {
		t.Skipf("plugin buildmode is unsupported on %s", runtime.GOOS)
	}
	if testPluginErr != nil {
		t.Skipf("building the test plugin failed: %v", testPluginErr)
	}
	return testPluginPath
}

func Test_LoadPlugin_With_InvalidPluginPath_Should_ReturnError(t *testing.T) {
	// Arrange & Act
	_, err := extensibility.LoadPlugin[TestPort]("testdata/adapter2.so", "Adapter")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_LoadPlugin_With_InvalidSymbol_Should_ReturnError(t *testing.T) {
	// Arrange
	path := requireTestPlugin(t)

	// Act
	_, err := extensibility.LoadPlugin[TestPort](path, "Adapter2")

	// Assert
	assert.That(t, "err must not be nil", err != nil, true)
}

func Test_LoadPlugin_With_ValidPlugin_Should_ReturnAdapter(t *testing.T) {
	// Arrange
	path := requireTestPlugin(t)

	// Act
	adapter, err := extensibility.LoadPlugin[TestPort](path, "Adapter")
	if err != nil {
		t.Fatalf("loading the plugin failed: %v", err)
	}
	name, err := adapter.FindByID("1")

	// Assert
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "adapter.FindByID must return 'Andy'", name, "Andy")
}
