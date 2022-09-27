// Package chtml provides convenience functions to wrap the standard html/template package.
package chtml

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
)

func Execute(wr io.Writer, t *template.Template, data any) error {
	return t.ExecuteTemplate(wr, "base", data)
}

// Layout is used to combine a layout and partials with a content template.
type Layout struct {
	// layoutPath is the path to a html template file that defines regions
	// that will be filled in by partial and page content templates.
	//
	// Layouts should define a region named "base" that contains all the html.
	// This is so the root name of the template will be "base" and then we can
	// use ExecuteTemplate to render the template.
	layoutPath *string
	// partialPattern is a glob pattern that matches html partial templates.
	partialPattern *string
	// funcMap is a map of functions that will be available to the templates.
	funcMap template.FuncMap
}

func NewLayout() *Layout {
	return &Layout{
		funcMap: template.FuncMap{},
	}
}

func (l *Layout) WithLayout(filepath string) *Layout {
	l.layoutPath = &filepath
	return l
}

func (l *Layout) WithPartialPattern(pattern string) *Layout {
	l.partialPattern = &pattern
	return l
}

func (l *Layout) WithFunc(name string, fn any) *Layout {
	l.funcMap[name] = fn
	return l
}

// ParseFS, accepts a FS and parses the layout, partial, and template at filename
// into a template.Template.
//
// The name of the tempalte will be the top level {{define}} name in the template.
// This should be "base" to keep things simple.
func (l *Layout) ParseFS(fsys fs.FS, filename string) (*template.Template, error) {
	tmpl := template.New("").Funcs(l.funcMap)
	var err error

	if l.layoutPath != nil {
		tmpl, err = tmpl.ParseFS(fsys, *l.layoutPath)
		if err != nil {
			return nil, fmt.Errorf("Layout.ParseFS failed to parse layout: %w", err)
		}
	}

	if l.partialPattern != nil {
		tmpl, err = tmpl.ParseFS(fsys, *l.partialPattern)
		if err != nil {
			return nil, fmt.Errorf("Layout.ParseFS failed to parse partials: %w", err)
		}
	}

	tmpl, err = tmpl.ParseFS(fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("Layout.ParseFS failed to parse template: %w", err)
	}

	return tmpl, nil
}
