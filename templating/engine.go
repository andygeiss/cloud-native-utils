package templating

import (
	"fmt"
	"io"
	"text/template"
)

// Engine is a simple wrapper around the Go templating engine.
type Engine struct {
	tmpl *template.Template
}

// NewEngine creates a new templating engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Parse parses the templates in the given path.
func (a *Engine) Parse(path string) {
	tmpl, err := template.ParseGlob(path)
	if err != nil {
		panic(fmt.Sprintf("templating: could not parse templates: %v", err))
	}
	a.tmpl = tmpl
}

// Render renders the template with the given name and data.
func (a *Engine) Render(w io.Writer, name string, data any) error {
	return a.tmpl.ExecuteTemplate(w, name, data)
}
