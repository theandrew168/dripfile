package web

import (
	"io"
)

type IndexParams struct{}

func (v *View) Index(w io.Writer, p IndexParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/site.html",
		"partial/*.html",
		"index.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}
