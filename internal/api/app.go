package api

import (
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/service"
	"github.com/theandrew168/dripfile/internal/view"
)

type Application struct {
	logger *jsonlog.Logger
	view   *view.View
	srvc   *service.Service
}

func NewApplication(
	logger *jsonlog.Logger,
	view *view.View,
	srvc *service.Service,
) *Application {
	app := Application{
		logger: logger,
		view:   view,
		srvc:   srvc,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
