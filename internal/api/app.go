package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/view"
	"github.com/theandrew168/dripfile/internal/jsonlog"
)

type Application struct {
	logger *jsonlog.Logger
	view   *view.Template
}

func NewApplication(logger *jsonlog.Logger, view *view.Template) *Application {
	app := Application{
		logger: logger,
		view:   view,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
