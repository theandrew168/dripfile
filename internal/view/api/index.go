package api

import (
	"io"
)

type IndexParams struct{}

func (t *Template) Index(w io.Writer, p IndexParams) error {
	patterns := []string{
		"index.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
