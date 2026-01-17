package templating_test

import (
	"bytes"
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/templating"
)

//go:embed testdata/*.tmpl
var efs embed.FS

func Test_Engine_With_AddIntTemplate_Should_CalculateSum(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")
	var result bytes.Buffer

	// Act
	err := engine.Render(&result, "add_int", struct{ A, B int }{A: 1, B: 2})

	// Assert
	assert.That(t, "engine.Render must succeed", err, nil)
	assert.That(t, "engine.Render must succeed", result.String(), "3")
}

func Test_Engine_With_MultipleRenders_Should_RenderCorrectly(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")

	// Act
	var result1 bytes.Buffer
	err1 := engine.Render(&result1, "index", struct{ Name string }{Name: "Alice"})

	var result2 bytes.Buffer
	err2 := engine.Render(&result2, "index", struct{ Name string }{Name: "Bob"})

	// Assert
	assert.That(t, "first render must succeed", err1, nil)
	assert.That(t, "second render must succeed", err2, nil)
	assert.That(t, "first render must succeed", result1.String(), "\nHello Alice\n")
	assert.That(t, "second render must succeed", result2.String(), "\nHello Bob\n")
}

func Test_Engine_With_NilData_Should_ReturnError(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")
	var result bytes.Buffer

	// Act
	err := engine.Render(&result, "add_int", nil)

	// Assert
	assert.That(t, "render with nil data should fail", err != nil, true)
}

func Test_Engine_With_NotExistingTemplate_Should_ReturnError(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")

	// Act
	err := engine.Render(nil, "not-existing.tmpl", nil)

	// Assert
	assert.That(t, "engine.Render must fail", err != nil, true)
}

func Test_Engine_With_ParsePattern_Should_NotPanic(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)

	// Act & Assert (no panic means success)
	engine.Parse("testdata/*.tmpl")
}

func Test_Engine_With_ValidInput_Should_ReturnNonNilEngine(t *testing.T) {
	// Arrange & Act
	engine := templating.NewEngine(efs)

	// Assert
	assert.That(t, "engine should not be nil", engine != nil, true)
}

func Test_Engine_With_ValidTemplate_Should_RenderCorrectly(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")
	var result bytes.Buffer

	// Act
	err := engine.Render(&result, "index", struct{ Name string }{Name: "World"})

	// Assert
	assert.That(t, "engine.Render must succeed", err, nil)
	assert.That(t, "engine.Render must succeed", result.String(), "\nHello World\n")
}

func Test_Engine_View_With_ValidTemplate_Should_ReturnHandler(t *testing.T) {
	// Arrange
	engine := templating.NewEngine(efs)
	engine.Parse("testdata/*.tmpl")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Act
	handler := engine.View("index", struct{ Name string }{Name: "View"})
	handler(w, req)

	// Assert
	assert.That(t, "handler must return 200", w.Code, http.StatusOK)
	assert.That(t, "handler must render content", w.Body.String(), "\nHello View\n")
}
