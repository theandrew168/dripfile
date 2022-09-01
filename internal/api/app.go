package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/template"
)

type Application struct {
	logger *jsonlog.Logger
	tmpl   *template.Map
}

func NewApplication(logger *jsonlog.Logger, tmpl *template.Map) *Application {
	app := Application{
		logger: logger,
		tmpl:   tmpl,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
