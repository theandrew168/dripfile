package web

import (
	"io"
)

type IndexParams struct{}

func (t *Template) Index(w io.Writer, p IndexParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/site.html",
		"partial/*.html",
		"index.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
