package app

import (
	"io"
)

type DashboardParams struct{}

func (t *Template) Dashboard(w io.Writer, p DashboardParams) error {
	patterns := []string{
		"layout/base.html",
		"dashboard.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
