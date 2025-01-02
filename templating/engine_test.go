package templating_test

import (
	"bytes"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/templating"
)

func TestEngine_Parse_Must_Succeed(t *testing.T) {
	engine := templating.NewEngine()
	engine.Parse("testdata/*.tmpl")
}

func TestEngine_Render_Must_Succeed(t *testing.T) {
	engine := templating.NewEngine()
	engine.Parse("testdata/*.tmpl")
	var result bytes.Buffer
	engine.Render(&result, "index", struct{ Name string }{Name: "World"})
	assert.That(t, "engine.Render must succeed", result.String(), "\nHello World\n")
}

func TestEngine_Render_Must_Fail(t *testing.T) {
	engine := templating.NewEngine()
	engine.Parse("testdata/*.tmpl")
	err := engine.Render(nil, "not-existing.tmpl", nil)
	assert.That(t, "engine.Render must fail", err != nil, true)
}
