package extensibility_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/extensibility"
)

type TestPort interface {
	FindByID(id string) (name string, err error)
}

func TestLoadPlugin(t *testing.T) {
	adapter, err := extensibility.LoadPlugin[TestPort]("testdata/adapter.so", "Adapter")
	name, err := adapter.FindByID("1")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "adapter.FindByID must return 'Andy'", name, "Andy")
}
