package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/validator"
)

type AccountDeleteForm struct {
	validator.Validator `form:"-"`

	AccountID string `form:"AccountID"`
}

type AccountReadParams struct {
	Account model.Account
}

func (t *Template) AccountRead(w io.Writer, p AccountReadParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"account/read.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
