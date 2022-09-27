// asset is a tool for rendering static html templates during development
// so I don't need to run the server to work on page styles.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/rusq/dlog"
	"github.com/rusq/slackinvite/internal/chtml"
)

type params struct {
	layoutPath     string
	partialPattern string
	contentPath    string
}

var cmdline params

func init() {
	flag.StringVar(&cmdline.layoutPath, "layout", "templates/layout.html", "layout `template`")
	flag.StringVar(&cmdline.partialPattern, "partials", "", "partial `templates`")
	flag.StringVar(&cmdline.contentPath, "content", "templates/index.html", "content `template`")
}

func main() {
	flag.Parse()

	if cmdline.layoutPath == "" {
		flag.Usage()
		dlog.Fatal("layout filepath not present")
	}

	if cmdline.contentPath == "" {
		flag.Usage()
		dlog.Fatal("content filepath not present")
	}

	layout := chtml.NewLayout()

	if cmdline.layoutPath != "" {
		layout = layout.WithLayout(cmdline.layoutPath)
	}

	if cmdline.partialPattern != "" {
		layout = layout.WithPartialPattern(cmdline.partialPattern)
	}

	tmpl, err := layout.ParseFS(os.DirFS("."), cmdline.contentPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = chtml.Execute(os.Stdout, tmpl, nil); err != nil {
		log.Fatal(err)
	}
}
