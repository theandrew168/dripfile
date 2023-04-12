package web

import (
	"io"
)

type IndexParams struct{}

func (html *HTML) Index(w io.Writer, p IndexParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/site.html",
		"partial/*.html",
		"index.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}
