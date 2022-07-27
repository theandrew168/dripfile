package api

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/jsonlog"
)

//go:embed template
var templateFS embed.FS

type Application struct {
	index *template.Template

	logger *jsonlog.Logger
}

func NewApplication(logger *jsonlog.Logger) *Application {
	dir, err := fs.Sub(templateFS, "template")
	if err != nil {
		panic(err)
	}

	index, err := template.ParseFS(dir, "index.html")
	if err != nil {
		panic(err)
	}

	app := Application{
		index: index,

		logger: logger,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
