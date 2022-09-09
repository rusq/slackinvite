// Package chtml provides convenience functions to wrap the standard html/template package.
package chtml

import (
	"fmt"
	"io/fs"
	"text/template"
)

// Layout is used to combine a layout and partials with a content template.
//
// Usage:
// l := chtml.Layout{
// 	LayoutPath: "layout.html",
// 	PartialPattern: "partials/*.html",
// 	FuncMap: template.FuncMap{
// 		"foo": func() string {
// 			return "bar"
// 		},
// 	},
// }
// tmpl, err := l.ReadFS(fsys, "content.html")
// if err != nil {
// 	log.Fatal(err) // shouldn't happen :huehuehue: :trollface:
// }
// err = tmpl.Execute(os.Stdout, nil)
// if err != nil {
// 	log.Fatal(err) // shouldn't happen :huehuehue: :trollface:
// }
type Layout struct {
	// LayoutPath is the path to a html template file that defines regions
	// that will be filled in by partial and page content templates.
	LayoutPath string
	// PartialPattern is a glob pattern that matches html partial templates.
	PartialPattern string
	// FuncMap is a map of functions that will be available to the templates.
	FuncMap template.FuncMap
}

// ReadFS, accepts a GlobFS and parses the layout, partial, and template at filename
// into a template.Template.
//
// The name of the tempalte will be an empty string so it is best to use Execute
// rather then ExecuteTemplate.
func (l *Layout) ReadFS(fsys fs.GlobFS, filename string) (*template.Template, error) {
	tmpl := template.New("").Funcs(l.FuncMap)

	tmpl, err := tmpl.ParseFS(fsys, l.LayoutPath)
	if err != nil {
		return nil, fmt.Errorf("Layout.ReadFS failed to parse layout: %w", err)
	}

	tmpl, err = tmpl.ParseFS(fsys, l.PartialPattern)
	if err != nil {
		return nil, fmt.Errorf("Layout.ReadFS failed to parse partials: %w", err)
	}

	tmpl, err = tmpl.ParseFS(fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("Layout.ReadFS failed to parse template: %w", err)
	}

	return tmpl, nil
}
