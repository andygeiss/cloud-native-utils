package templating

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"text/template"
)

// Engine is a simple wrapper around the Go templating engine.
type Engine struct {
	efs  fs.FS
	tmpl *template.Template
}

// NewEngine creates a new templating engine.
func NewEngine(efs fs.FS) *Engine {
	return &Engine{
		efs: efs,
	}
}

// Parse parses the templates in the given path.
func (a *Engine) Parse(patterns ...string) {
	// Add custom functions.
	tmpl := template.New("tmpl").Funcs(template.FuncMap{
		"add_int": func(a, b int) int { return a + b },
	})
	parsed, err := tmpl.ParseFS(a.efs, patterns...)
	if err != nil {
		panic(fmt.Sprintf("templating: could not parse templates: %v", err))
	}
	a.tmpl = parsed
}

// Render renders the template with the given name and data.
func (a *Engine) Render(w io.Writer, name string, data any) error {
	return a.tmpl.ExecuteTemplate(w, name, data)
}

// View returns an http.HandlerFunc that renders the template with the given name and data.
func (a *Engine) View(name string, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.Render(w, name, data); err != nil {
			http.Error(w, fmt.Sprintf("templating: render %q: %v", name, err), http.StatusInternalServerError)
			return
		}
	}
}
