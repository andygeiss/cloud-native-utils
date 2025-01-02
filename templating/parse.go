package templating

import (
	"fmt"
	"text/template"
)

// Parse creates a new template and parses it.
// If the template is invalid, the function will panic.
func Parse(tmpl string) *template.Template {
	// Create a new template and parse it.
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		panic(fmt.Sprintf("error parsing template: %v", err))
	}
	return t
}
