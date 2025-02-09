package templating

import (
	"embed"
	"fmt"
	"io"
	"text/template"
)

// Engine is a simple wrapper around the Go templating engine.
type Engine struct {
	efs  embed.FS
	tmpl *template.Template
}

// NewEngine creates a new templating engine.
func NewEngine(efs embed.FS) *Engine {
	return &Engine{
		efs: efs,
	}
}

// Parse parses the templates in the given path.
func (a *Engine) Parse(patterns ...string) {
	tmpl, err := template.ParseFS(a.efs, patterns...)
	if err != nil {
		panic(fmt.Sprintf("templating: could not parse templates: %v", err))
	}
	// Add custom functions.
	tmpl = tmpl.Funcs(template.FuncMap{
		"add_int": func(a, b int) int { return a + b },
	})
	a.tmpl = tmpl
}

// Render renders the template with the given name and data.
func (a *Engine) Render(w io.Writer, name string, data any) error {
	return a.tmpl.ExecuteTemplate(w, name, data)
}
