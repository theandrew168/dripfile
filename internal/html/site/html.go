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
	tmpl := t.parse(page)
	return tmpl.Execute(w, p)
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
	tmpl := t.parse(page)
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
	page := "auth/register.html"
	tmpl := t.parse(page)
	return tmpl.Execute(w, p)
}

func (t *Template) parse(page string) *template.Template {
	// load from cache if not in debug mode
	if !t.debug {
		tmpl, ok := t.cache[page]
		if ok {
			return tmpl
		}
	}

	patterns := []string{
		"layout.html",
		"partial/*.html",
		page,
	}
	tmpl, err := template.New(patterns[0]).ParseFS(t.files, patterns...)
	if err != nil {
		panic(err)
	}

	t.cache[page] = tmpl
	return tmpl
}
