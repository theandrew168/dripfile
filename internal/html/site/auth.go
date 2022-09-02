package site

import (
	"io"

	"github.com/theandrew168/dripfile/internal/validator"
)

type AuthLoginForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
}

type AuthLoginParams struct {
	Form AuthLoginForm
}

func (t *Template) AuthLogin(w io.Writer, p AuthLoginParams) error {
	patterns := []string{
		"layout/base.html",
		"partial/*.html",
		"auth/login.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}

type AuthRegisterForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
}

type AuthRegisterParams struct {
	Form AuthRegisterForm
}

func (t *Template) AuthRegister(w io.Writer, p AuthRegisterParams) error {
	patterns := []string{
		"layout/base.html",
		"partial/*.html",
		"auth/register.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
