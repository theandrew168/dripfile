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

func (html *HTML) Dashboard(w io.Writer, p DashboardParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"dashboard.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}
