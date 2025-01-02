package templating

import (
	"bytes"
	"text/template"
)

// Execute executes the given template with the given data and returns the result.
func Execute(tmpl *template.Template, data any) (string, error) {
	// Execute the template and return the result.
	var result bytes.Buffer
	err := tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
