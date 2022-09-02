package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/html"
	"github.com/theandrew168/dripfile/internal/jsonlog"
)

type Application struct {
	logger *jsonlog.Logger
	html   *html.Template
}

func NewApplication(logger *jsonlog.Logger, html *html.Template) *Application {
	app := Application{
		logger: logger,
		html:   html,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
