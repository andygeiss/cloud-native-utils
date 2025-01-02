package templating_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/templating"
)

func TestParse_Must_Succeed(t *testing.T) {
	tmpl := templating.Parse("Hello {{.Name}}")
	assert.That(t, "template must not be nil", tmpl != nil, true)
}

func TestParse_Must_Fail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	templating.Parse("Hello {{.Name}")
}
