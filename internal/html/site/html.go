package site

import (
	"embed"
	"html/template"
	"io"

	"github.com/theandrew168/dripfile/internal/validator"
)

//go:embed *
var files embed.FS

var (
	index        = parse("index.html")
	authLogin    = parse("auth/login.html")
	authRegister = parse("auth/register.html")
)

type IndexParams struct{}

func Index(w io.Writer, p IndexParams) error {
	return index.Execute(w, p)
}

type AuthLoginForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
}

type AuthLoginParams struct {
	Form AuthLoginForm
}

func AuthLogin(w io.Writer, p AuthLoginParams) error {
	return authLogin.Execute(w, p)
}

func parse(file string) *template.Template {
	t, err := template.New("layout.html").ParseFS(files, "layout.html", file)
	return template.Must(t, err)
}
