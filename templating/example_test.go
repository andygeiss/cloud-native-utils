package templating_test

import (
	"bytes"
	"fmt"
	"testing/fstest"

	"github.com/andygeiss/cloud-native-utils/templating"
)

func ExampleEngine() {
	fsys := fstest.MapFS{
		"page.tmpl": &fstest.MapFile{Data: []byte(`{{ define "page" }}Hello {{ .Name }}{{ end }}`)},
	}

	engine := templating.NewEngine(fsys)
	engine.Parse("*.tmpl")

	// Data is escaped for the context it lands in.
	var out bytes.Buffer
	_ = engine.Render(&out, "page", struct{ Name string }{Name: "<b>Alice</b>"})

	fmt.Println(out.String())
	// Output: Hello &lt;b&gt;Alice&lt;/b&gt;
}
