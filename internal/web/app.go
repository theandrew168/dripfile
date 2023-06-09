package web

import (
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"
)

type Application struct {
	logger *slog.Logger
}

func NewApplication(
	logger *slog.Logger,
) *Application {
	app := Application{
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
