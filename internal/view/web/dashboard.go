package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/validator"
)

type DashboardForm struct {
	validator.Validator `form:"-"`

	Search string `form:"Search"`
}

type DashboardParams struct {
	Form DashboardForm
}

func (t *Template) Dashboard(w io.Writer, p DashboardParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"dashboard.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
