package templating_test

import (
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/templating"
)

func TestExecute_Must_Succeed(t *testing.T) {
	// Given a template and some data.
	tmpl := templating.Parse("Hello {{.Name}}")
	data := struct {
		Name string
	}{
		Name: "World",
	}

	// When the template is executed.
	got, err := templating.Execute(tmpl, data)

	// Then the result must be 'Hello World'.
	assert.That(t, "error must be nil", err, nil)
	assert.That(t, "result must be 'Hello World'", got, "Hello World")
}

func TestExecute_Must_Fail(t *testing.T) {
	// Given a template and some data.
	tmpl := templating.Parse("Hello {{.Name}}")
	data := struct {
		NotFound string
	}{
		NotFound: "World",
	}

	// When the template is executed.
	_, err := templating.Execute(tmpl, data)

	// Then the error must be not nil.
	assert.That(t, "error must be not nil", err != nil, true)
}

func BenchmarkExecute(b *testing.B) {
	tmpl := templating.Parse("Hello {{.Name}}")
	data := struct {
		Name string
	}{
		Name: "World",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = templating.Execute(tmpl, data)
	}
}
