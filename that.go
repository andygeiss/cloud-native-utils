package assert

import (
	"reflect"
	"testing"
)

// That is a utility function to assert that two values are equal during a test.
func That(t *testing.T, desc string, got, expected any) {
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("%s, but got %v (expected %v)", desc, got, expected)
	}
}
