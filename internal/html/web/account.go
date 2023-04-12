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

func (html *HTML) AccountRead(w io.Writer, p AccountReadParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"account/read.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}
