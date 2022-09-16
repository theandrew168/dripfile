package api

import (
	"io"
)

type IndexParams struct{}

func (v *View) Index(w io.Writer, p IndexParams) error {
	patterns := []string{
		"index.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}
