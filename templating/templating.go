// Package templating provides an `Engine` for managing templates stored in an
// embedded filesystem. Use `Parse` to load multiple templates (via glob
// patterns), and `Render` to execute them with custom data.
//
// Templates are HTML: data is escaped for the context it lands in, so a value
// carrying markup is rendered as text. Pass template.HTML to opt a value out,
// and only for markup you produced yourself.
package templating
