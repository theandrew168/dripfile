package site

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"

	"github.com/theandrew168/dripfile/internal/validator"
)

//go:embed template
var templateFS embed.FS

type Template struct {
	debug bool
	files fs.FS
	cache map[string]*template.Template
}

func New() *Template {
	debug := false
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	
	var files fs.FS
	if debug {
		// reload templates from filesystem if env var DEBUG is defined
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		files = os.DirFS("./internal/html/site/template/")
	} else {
		// else use the embedded template dir
		files, _ = fs.Sub(templateFS, "template")
	}

	cache := make(map[string]*template.Template)

	t := Template{
		debug: debug,
		files: files,
		cache: cache,
	}
	return &t
}

type IndexParams struct{}

func (t *Template) Index(w io.Writer, p IndexParams) error {
	page := "index.html"
}

type AuthLoginForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
}

type AuthLoginParams struct {
	Form AuthLoginForm
}

func (t *Template) AuthLogin(w io.Writer, p AuthLoginParams) error {
	page := "auth/login.html"

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
	page := "auth/register.html"

}

func (t *Template) parse(page string) *template.Template {
	base := "layout.html"
	tmpl, err := template.New(base).ParseFS(t.files, base, page)
	return template.Must(tmpl, err)
}
